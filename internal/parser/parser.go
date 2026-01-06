package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"sovana/internal/model"
)

var dateLayouts = []string{
	"2006-01-02",
	"02/01/2006",
	"02-01-2006",
}

var dateRegexes = []*regexp.Regexp{
	regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`),
	regexp.MustCompile(`\b\d{2}/\d{2}/\d{4}\b`),
	regexp.MustCompile(`\b\d{2}-\d{2}-\d{4}\b`),
}

var amountRegex = regexp.MustCompile(`[-+]?\d{1,3}(?:[\s.,]\d{3})*(?:[.,]\d+)?`)

type columnIndexes struct {
	date   int
	label  int
	amount int
	debit  int
	credit int
}

// ParsePDF reads a PDF bank statement and converts it to a slice of transactions.
// The parser extracts plain text, automatically detects separators (; or ,), identifies
// date/label/amount or debit/credit columns, tolerates invalid lines, and logs ignored
// rows at debug level.
func ParsePDF(r io.Reader) ([]model.Transaction, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read pdf: %w", err)
	}
	if len(data) == 0 {
		return nil, errors.New("empty pdf")
	}

	text, err := extractPDFText(data)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("pdf contains no text to parse")
	}

	rows, lines := buildRowsFromText(text)

	tabular := parseTabular(rows)
	var transactions []model.Transaction
	if len(tabular) > 0 {
		transactions = mergeTransactions(tabular, parseLooseLines(filterLooseCandidates(lines)))
	} else {
		transactions = parseLooseLines(lines)
	}

	if len(transactions) == 0 {
		return nil, errors.New("no valid transactions found")
	}

	return transactions, nil
}

func extractPDFText(data []byte) (string, error) {
	content := string(data)
	if !strings.Contains(content, "%PDF") {
		return content, nil
	}

	re := regexp.MustCompile(`\(([^()]*)\)`)
	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		lines := make([]string, 0, len(matches))
		for _, m := range matches {
			if len(m) > 1 {
				lines = append(lines, m[1])
			}
		}
		return strings.Join(lines, "\n"), nil
	}

	lower := strings.ToLower(content)
	var buf strings.Builder
	start := 0
	for {
		streamIdx := strings.Index(lower[start:], "stream")
		if streamIdx == -1 {
			break
		}
		streamIdx += start
		endIdx := strings.Index(lower[streamIdx:], "endstream")
		if endIdx == -1 {
			break
		}
		endIdx += streamIdx
		fragment := content[streamIdx+len("stream") : endIdx]
		buf.WriteString(fragment)
		buf.WriteString("\n")
		start = endIdx + len("endstream")
	}

	if buf.Len() > 0 {
		return buf.String(), nil
	}

	return content, nil
}

func buildRowsFromText(text string) ([][]string, []string) {
	sep := detectSeparator([]byte(text))
	lines := strings.Split(text, "\n")
	rows := make([][]string, 0, len(lines))
	cleaned := make([]string, 0, len(lines))
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		cleaned = append(cleaned, line)
		line = strings.ReplaceAll(line, "\t", string(sep))

		var fields []string
		if strings.Contains(line, string(sep)) {
			fields = strings.Split(line, string(sep))
			if len(fields) < 3 && sep == ',' {
				fields = strings.Fields(line)
			}
		} else {
			fields = strings.Fields(line)
		}
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		if len(fields) == 0 {
			continue
		}
		rows = append(rows, fields)
	}
	return rows, cleaned
}

func detectSeparator(data []byte) rune {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		semicolons := strings.Count(line, ";")
		commas := strings.Count(line, ",")
		switch {
		case semicolons > commas:
			return ';'
		case commas > semicolons:
			return ','
		case semicolons > 0:
			return ';'
		}
	}
	return ','
}

func detectColumns(rows [][]string) (columnIndexes, int) {
	if len(rows) == 0 {
		return columnIndexes{date: 0, label: 1, amount: 2, debit: -1, credit: -1}, 0
	}
	header := rows[0]
	normalized := make([]string, len(header))
	for i, col := range header {
		normalized[i] = normalizeHeader(col)
	}

	idx := columnIndexes{date: -1, label: -1, amount: -1, debit: -1, credit: -1}
	for i, col := range normalized {
		switch {
		case strings.Contains(col, "date"):
			idx.date = i
		case strings.Contains(col, "libelle") || strings.Contains(col, "label") || strings.Contains(col, "description") || strings.Contains(col, "detail") || strings.Contains(col, "operation"):
			if idx.label == -1 {
				idx.label = i
			}
		case strings.Contains(col, "montant") || strings.Contains(col, "amount") || strings.Contains(col, "prix"):
			idx.amount = i
		case strings.Contains(col, "debit") || strings.Contains(col, "retrait"):
			idx.debit = i
		case strings.Contains(col, "credit") || strings.Contains(col, "versement") || strings.Contains(col, "ajout"):
			idx.credit = i
		}
	}

	hasHeader := idx.date >= 0 && idx.label >= 0 && (idx.amount >= 0 || idx.debit >= 0 || idx.credit >= 0)
	if hasHeader {
		return idx, 1
	}

	inferred := inferColumns(rows)
	return inferred, 0
}

func parseDate(value string) (time.Time, error) {
	trimmed := strings.TrimSpace(value)
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", value)
}

func parseAmount(value string) (float64, error) {
	cleaned := strings.ReplaceAll(strings.TrimSpace(value), " ", "")
	cleaned = strings.ReplaceAll(cleaned, "€", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, "\u00a0", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	if cleaned == "" {
		return 0, fmt.Errorf("empty amount")
	}
	amount, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %s", value)
	}
	return amount, nil
}

func extractAmount(row []string, idx columnIndexes) (float64, bool) {
	if idx.amount >= 0 && idx.amount < len(row) {
		val := fieldAt(row, idx.amount)
		amt, err := parseAmount(val)
		if err == nil {
			return amt, true
		}
	}

	debitVal := 0.0
	creditVal := 0.0
	hasDebit := false
	hasCredit := false

	if idx.debit >= 0 && idx.debit < len(row) {
		if amt, err := parseAmount(fieldAt(row, idx.debit)); err == nil {
			debitVal = math.Abs(amt)
			hasDebit = true
		}
	}
	if idx.credit >= 0 && idx.credit < len(row) {
		if amt, err := parseAmount(fieldAt(row, idx.credit)); err == nil {
			creditVal = math.Abs(amt)
			hasCredit = true
		}
	}

	switch {
	case hasDebit && hasCredit:
		if debitVal != 0 {
			return -debitVal, true
		}
		if creditVal != 0 {
			return creditVal, true
		}
	case hasDebit:
		if debitVal != 0 {
			return -debitVal, true
		}
	case hasCredit:
		if creditVal != 0 {
			return creditVal, true
		}
	}
	return 0, false
}

func isNonTransactionLabel(label string) bool {
	lower := strings.ToLower(label)
	blocked := []string{"ajout de fonds", "solde", "virement interne", "balance"}
	for _, b := range blocked {
		if strings.Contains(lower, b) {
			return true
		}
	}
	return false
}

func normalizeHeader(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	replacer := strings.NewReplacer("é", "e", "è", "e", "ê", "e", "à", "a", "ù", "u", "ç", "c")
	lower = replacer.Replace(lower)
	lower = strings.ReplaceAll(lower, "_", "")
	lower = strings.ReplaceAll(lower, " ", "")
	lower = strings.ReplaceAll(lower, "-", "")
	return lower
}

func inferColumns(rows [][]string) columnIndexes {
	stats := make([]struct {
		dateCount   int
		amountCount int
		negCount    int
		posCount    int
		totalLen    int
		samples     int
	}, maxColumns(rows))

	for _, row := range rows {
		for i, field := range row {
			if _, err := parseDate(field); err == nil {
				stats[i].dateCount++
			}
			if amt, err := parseAmount(field); err == nil {
				stats[i].amountCount++
				if amt < 0 {
					stats[i].negCount++
				} else {
					stats[i].posCount++
				}
			}
			stats[i].totalLen += len(field)
			stats[i].samples++
		}
	}

	idx := columnIndexes{date: -1, label: -1, amount: -1, debit: -1, credit: -1}
	idx.date = bestIndex(stats, func(s *struct{ dateCount, amountCount, negCount, posCount, totalLen, samples int }) int {
		return s.dateCount
	})

	amountCandidates := indexesAbove(stats, func(s *struct{ dateCount, amountCount, negCount, posCount, totalLen, samples int }) int {
		return s.amountCount
	}, 0)
	switch len(amountCandidates) {
	case 0:
		idx.amount = -1
	case 1:
		idx.amount = amountCandidates[0]
	default:
		debitIdx := -1
		creditIdx := -1
		for _, c := range amountCandidates {
			if stats[c].negCount > stats[c].posCount && debitIdx == -1 {
				debitIdx = c
				continue
			}
			if stats[c].posCount >= stats[c].negCount && creditIdx == -1 {
				creditIdx = c
			}
		}
		if debitIdx != -1 || creditIdx != -1 {
			idx.debit = debitIdx
			idx.credit = creditIdx
		}
		if idx.debit == -1 && idx.credit == -1 {
			idx.amount = bestIndex(stats, func(s *struct{ dateCount, amountCount, negCount, posCount, totalLen, samples int }) int {
				return s.amountCount
			})
		}
	}

	idx.label = bestIndex(stats, func(s *struct{ dateCount, amountCount, negCount, posCount, totalLen, samples int }) int {
		return s.totalLen
	}, excludeIndexes(idx)...)

	if idx.label == -1 {
		idx.label = 1
	}
	if idx.date == -1 {
		idx.date = 0
	}
	return idx
}

func maxColumns(rows [][]string) int {
	max := 0
	for _, row := range rows {
		if len(row) > max {
			max = len(row)
		}
	}
	return max
}

func bestIndex(stats []struct {
	dateCount   int
	amountCount int
	negCount    int
	posCount    int
	totalLen    int
	samples     int
}, scorer func(*struct {
	dateCount   int
	amountCount int
	negCount    int
	posCount    int
	totalLen    int
	samples     int
}) int, exclude ...int) int {
	excluded := make(map[int]struct{})
	for _, e := range exclude {
		excluded[e] = struct{}{}
	}
	best := -1
	bestScore := 0
	for i := range stats {
		if _, ok := excluded[i]; ok {
			continue
		}
		score := scorer(&stats[i])
		if score > bestScore {
			bestScore = score
			best = i
		}
	}
	return best
}

func indexesAbove(stats []struct {
	dateCount   int
	amountCount int
	negCount    int
	posCount    int
	totalLen    int
	samples     int
}, scorer func(*struct {
	dateCount   int
	amountCount int
	negCount    int
	posCount    int
	totalLen    int
	samples     int
}) int, threshold int) []int {
	res := []int{}
	for i := range stats {
		if scorer(&stats[i]) > threshold {
			res = append(res, i)
		}
	}
	return res
}

func excludeIndexes(idx columnIndexes) []int {
	res := []int{}
	for _, v := range []int{idx.date, idx.amount, idx.debit, idx.credit} {
		if v >= 0 {
			res = append(res, v)
		}
	}
	return res
}

func parseTabular(rows [][]string) []model.Transaction {
	idx, start := detectColumns(rows)
	transactions := make([]model.Transaction, 0, len(rows)-start)

	for i := start; i < len(rows); i++ {
		row := rows[i]
		dateVal := fieldAt(row, idx.date)
		parsedDate, err := parseDate(dateVal)
		if err != nil {
			debugf("skipping row %d: invalid date '%s'", i+1, dateVal)
			continue
		}

		label := cleanLabel(fieldAt(row, idx.label))
		if label == "" && len(row) > 1 {
			label = cleanLabel(strings.Join(row, " "))
		}
		if isNonTransactionLabel(label) {
			debugf("skipping row %d: non-transaction label '%s'", i+1, label)
			continue
		}

		amount, ok := extractAmount(row, idx)
		if !ok {
			debugf("skipping row %d: unable to determine amount", i+1)
			continue
		}

		transactions = append(transactions, model.Transaction{Date: parsedDate, Label: label, Amount: amount})
	}
	return transactions
}

func parseLooseLines(lines []string) []model.Transaction {
	var txs []model.Transaction
	for i, line := range lines {
		dateVal, datePos, ok := findDate(line)
		if !ok {
			continue
		}
		parsedDate, err := parseDate(dateVal)
		if err != nil {
			continue
		}

		amtVal, amtPos, amtOk := findAmount(line)
		if !amtOk {
			continue
		}
		amount, err := parseAmount(amtVal)
		if err != nil {
			continue
		}

		label := cleanLabel(removeSpan(line, datePos, amtPos))
		if label == "" {
			label = cleanLabel(strings.ReplaceAll(removeSpan(line, amtPos, datePos), "  ", " "))
		}
		if label == "" {
			label = cleanLabel(line)
		}
		if isNonTransactionLabel(label) {
			debugf("skipping line %d: non-transaction label '%s'", i+1, label)
			continue
		}

		txs = append(txs, model.Transaction{Date: parsedDate, Label: label, Amount: amount})
	}
	return txs
}

func filterLooseCandidates(lines []string) []string {
	res := make([]string, 0, len(lines))
	for _, line := range lines {
		separators := strings.Count(line, ";") + strings.Count(line, ",")
		if separators >= 2 {
			continue
		}
		res = append(res, line)
	}
	return res
}

func findDate(line string) (string, [2]int, bool) {
	lower := strings.ToLower(line)
	for _, re := range dateRegexes {
		loc := re.FindStringIndex(lower)
		if loc != nil {
			return line[loc[0]:loc[1]], [2]int{loc[0], loc[1]}, true
		}
	}
	return "", [2]int{}, false
}

func findAmount(line string) (string, [2]int, bool) {
	matches := amountRegex.FindAllStringIndex(line, -1)
	if len(matches) == 0 {
		return "", [2]int{}, false
	}
	for i := len(matches) - 1; i >= 0; i-- {
		pos := matches[i]
		token := line[pos[0]:pos[1]]
		if isDateToken(token) || isEmbeddedInDate(line, pos) {
			continue
		}
		return token, [2]int{pos[0], pos[1]}, true
	}
	return "", [2]int{}, false
}

func isEmbeddedInDate(line string, pos []int) bool {
	start := pos[0]
	end := pos[1]
	if start > 0 {
		prev := line[start-1]
		if prev == '/' || prev == '-' {
			return true
		}
	}
	if end < len(line) {
		next := line[end]
		if next == '/' || next == '-' {
			return true
		}
	}
	return false
}

func removeSpan(line string, a, b [2]int) string {
	start := minInt(a[0], b[0])
	end := maxInt(a[1], b[1])
	if start == 0 && end == 0 {
		return line
	}
	var builder strings.Builder
	builder.WriteString(line[:start])
	builder.WriteString(" ")
	if end < len(line) {
		builder.WriteString(line[end:])
	}
	return builder.String()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mergeTransactions(existing, more []model.Transaction) []model.Transaction {
	if len(more) == 0 {
		return existing
	}
	seen := make(map[string]struct{}, len(existing)+len(more))
	seenAbs := make(map[string]struct{}, len(existing)+len(more))
	merged := make([]model.Transaction, 0, len(existing)+len(more))
	for _, tx := range existing {
		key := fmt.Sprintf("%s|%s|%.2f", tx.Date.Format(time.DateOnly), labelKey(tx.Label), tx.Amount)
		seen[key] = struct{}{}
		seenAbs[fmt.Sprintf("%s|%s|%.2f", tx.Date.Format(time.DateOnly), labelKey(tx.Label), math.Abs(tx.Amount))] = struct{}{}
		merged = append(merged, tx)
	}
	for _, tx := range more {
		key := fmt.Sprintf("%s|%s|%.2f", tx.Date.Format(time.DateOnly), labelKey(tx.Label), tx.Amount)
		if _, ok := seen[key]; ok {
			continue
		}
		absKey := fmt.Sprintf("%s|%s|%.2f", tx.Date.Format(time.DateOnly), labelKey(tx.Label), math.Abs(tx.Amount))
		if _, ok := seenAbs[absKey]; ok {
			continue
		}
		seen[key] = struct{}{}
		seenAbs[absKey] = struct{}{}
		merged = append(merged, tx)
	}
	return merged
}

func cleanLabel(label string) string {
	if label == "" {
		return ""
	}
	replacer := strings.NewReplacer(";", " ", ",", " ", "|", " ", "\t", " ")
	label = replacer.Replace(label)
	label = strings.TrimSpace(strings.Join(strings.Fields(label), " "))
	return label
}

func labelKey(label string) string {
	cleaned := cleanLabel(label)
	tokens := strings.Fields(cleaned)
	filtered := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		lowerTok := strings.ToLower(tok)
		if isDateToken(lowerTok) || amountRegex.MatchString(lowerTok) {
			continue
		}
		filtered = append(filtered, lowerTok)
	}
	if len(filtered) == 0 {
		return strings.ToLower(cleaned)
	}
	return strings.Join(filtered, " ")
}

func isDateToken(token string) bool {
	for _, re := range dateRegexes {
		if re.MatchString(token) {
			return true
		}
	}
	return false
}

func fieldAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func debugf(format string, args ...interface{}) {
	log.Printf("DEBUG: "+format, args...)
}
