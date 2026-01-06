package parser

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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

type columnIndexes struct {
	date   int
	label  int
	amount int
	debit  int
	credit int
}

// ParseCSV reads a CSV bank statement and converts it to a slice of transactions.
// The parser automatically detects separators (; or ,), identifies date/label/amount
// or debit/credit columns, tolerates invalid lines, and logs ignored rows at debug level.
func ParseCSV(r io.Reader) ([]model.Transaction, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(data) == 0 {
		return nil, errors.New("empty csv")
	}

	sep := detectSeparator(data)
	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.TrimLeadingSpace = true
	reader.Comma = sep
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse csv: %w", err)
	}
	if len(rows) == 0 {
		return nil, errors.New("empty csv")
	}

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

		label := strings.TrimSpace(fieldAt(row, idx.label))
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

	if len(transactions) == 0 {
		return nil, errors.New("no valid transactions found")
	}

	return transactions, nil
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
	if !hasHeader {
		if len(header) >= 3 {
			return columnIndexes{date: 0, label: 1, amount: 2, debit: -1, credit: -1}, 0
		}
		return columnIndexes{date: 0, label: 1, amount: -1, debit: -1, credit: -1}, 0
	}
	return idx, 1
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

func fieldAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func debugf(format string, args ...interface{}) {
	log.Printf("DEBUG: "+format, args...)
}
