package ws

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

const websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// Conn is a minimal websocket connection supporting text frames.
type Conn struct {
	conn net.Conn
	mu   sync.Mutex
}

// Upgrade performs a WebSocket handshake and returns a connection.
func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("websocket upgrade requires GET")
	}
	if !headerContainsToken(r.Header, "Connection", "Upgrade") || !headerContainsToken(r.Header, "Upgrade", "websocket") {
		return nil, errors.New("missing required upgrade headers")
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return nil, errors.New("missing Sec-WebSocket-Key")
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("response does not support hijacking")
	}
	nc, buf, err := hj.Hijack()
	if err != nil {
		return nil, err
	}

	accept := computeAccept(key)
	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", accept)
	if _, err := buf.WriteString(response); err != nil {
		nc.Close()
		return nil, err
	}
	if err := buf.Flush(); err != nil {
		nc.Close()
		return nil, err
	}

	return &Conn{conn: nc}, nil
}

// WriteText writes a single text frame to the connection.
func (c *Conn) WriteText(payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	header := []byte{0x81}
	length := len(payload)
	switch {
	case length < 126:
		header = append(header, byte(length))
	case length <= 65535:
		header = append(header, 126, byte(length>>8), byte(length))
	default:
		header = append(header, 127,
			byte(length>>56), byte(length>>48), byte(length>>40), byte(length>>32),
			byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
	}

	if _, err := c.conn.Write(header); err != nil {
		return err
	}
	_, err := c.conn.Write(payload)
	return err
}

// Close sends a close frame and closes the underlying connection.
func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// best-effort close frame
	c.conn.Write([]byte{0x88, 0x00})
	return c.conn.Close()
}

// Drain reads frames until the connection closes, useful to detect client disconnects.
func (c *Conn) Drain() error {
	reader := bufio.NewReader(c.conn)
	for {
		opcode, payloadLen, err := readFrameHeader(reader)
		if err != nil {
			return err
		}
		// ignore payload content
		if _, err := io.CopyN(io.Discard, reader, int64(payloadLen)); err != nil {
			return err
		}
		if opcode == 0x8 { // close frame
			return io.EOF
		}
	}
}

func readFrameHeader(r *bufio.Reader) (opcode byte, payloadLen int, err error) {
	first, err := r.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	second, err := r.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	opcode = first & 0x0F
	masked := second&0x80 != 0
	length := int(second & 0x7F)
	switch length {
	case 126:
		b1, err1 := r.ReadByte()
		b2, err2 := r.ReadByte()
		if err := firstErr(err1, err2); err != nil {
			return 0, 0, err
		}
		length = int(b1)<<8 | int(b2)
	case 127:
		var bytes [8]byte
		if _, err := io.ReadFull(r, bytes[:]); err != nil {
			return 0, 0, err
		}
		for _, b := range bytes {
			length = length<<8 + int(b)
		}
	}
	if masked {
		mask := make([]byte, 4)
		if _, err := io.ReadFull(r, mask); err != nil {
			return 0, 0, err
		}
	}
	return opcode, length, nil
}

func computeAccept(key string) string {
	h := sha1.Sum([]byte(key + websocketGUID))
	return base64.StdEncoding.EncodeToString(h[:])
}

func headerContainsToken(h http.Header, name, token string) bool {
	for _, v := range h.Values(name) {
		for _, part := range splitHeader(v) {
			if strings.EqualFold(part, token) {
				return true
			}
		}
	}
	return false
}

func splitHeader(v string) []string {
	parts := strings.Split(v, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
