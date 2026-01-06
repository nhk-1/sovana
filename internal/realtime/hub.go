package realtime

import (
	"context"
	"encoding/json"
	"log"

	"sovana/internal/metrics"
	"sovana/internal/ws"
)

// Hub manages concurrent websocket clients and broadcasts metrics.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan metrics.Metric
}

// NewHub initializes the hub ready to run.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan metrics.Metric, 16),
	}
}

// Run processes client registrations and broadcast requests.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for c := range h.clients {
				c.Close()
			}
			return
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
		case metric := <-h.broadcast:
			payload, err := json.Marshal(metric)
			if err != nil {
				log.Printf("serialize metric: %v", err)
				continue
			}
			for client := range h.clients {
				select {
				case client.send <- payload:
				default:
					// Slow client; drop connection to keep hub responsive.
					delete(h.clients, client)
					client.Close()
				}
			}
		}
	}
}

// Broadcast enqueues a metric for connected websocket clients.
func (h *Hub) Broadcast(m metrics.Metric) {
	h.broadcast <- m
}

// Register adds a client to the hub.
func (h *Hub) Register(c *Client) { h.register <- c }

// Unregister removes a client.
func (h *Hub) Unregister(c *Client) { h.unregister <- c }

// Client wraps a websocket connection.
type Client struct {
	hub  *Hub
	conn *ws.Conn
	send chan []byte
}

// NewClient prepares a websocket client.
func NewClient(hub *Hub, conn *ws.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 8),
	}
}

// WritePump writes messages to the websocket connection.
// It runs in its own goroutine to avoid blocking the hub.
func (c *Client) WritePump(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.Close()
			return
		case msg, ok := <-c.send:
			if !ok {
				c.Close()
				return
			}
			if err := c.conn.WriteText(msg); err != nil {
				c.Close()
				return
			}
		}
	}
}

// Close closes the underlying connection and channel.
func (c *Client) Close() {
	c.conn.Close()
	close(c.send)
}
