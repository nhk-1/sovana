package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"sovana/internal/model"
)

var dateLayouts = []string{
	"2006-01-02",
	"02/01/2006",
	"01/02/2006",
	"2006/01/02",
}

// ParseCSV reads a CSV bank statement and converts it to a slice of transactions.
// Expected columns: date, label/description, amount. A header row is optional.
// Amounts can use a comma or dot as decimal separator and negative values
// represent debits.
func ParseCSV(r io.Reader) ([]model.Transaction, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(rows) == 0 {
		return nil, errors.New("empty csv")
	}

	startIdx := 0
	if isHeader(rows[0]) {
		startIdx = 1
	}

	transactions := make([]model.Transaction, 0, len(rows)-startIdx)
	for i := startIdx; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			return nil, fmt.Errorf("row %d: expected at least 3 columns", i+1)
		}
		date, err := parseDate(row[0])
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+1, err)
		}
		label := strings.TrimSpace(row[1])
		amount, err := parseAmount(row[2])
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+1, err)
		}

		transactions = append(transactions, model.Transaction{Date: date, Label: label, Amount: amount})
	}

	return transactions, nil
}

func isHeader(row []string) bool {
	if len(row) < 3 {
		return false
	}
	joined := strings.ToLower(strings.Join(row, ","))
	return strings.Contains(joined, "date") && (strings.Contains(joined, "label") || strings.Contains(joined, "description"))
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
	amount, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %s", value)
	}
	return amount, nil
}
