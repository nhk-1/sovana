package parser

import (
	"strings"
	"testing"
)

func TestParseCSVSeparatorAndHeader(t *testing.T) {
	csvData := "Date;Libellé;Montant\n2024-01-01;Netflix;13,49\n2024-02-01;Spotify;9.99"
	txs, err := ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txs))
	}
	if txs[0].Amount != 13.49 {
		t.Fatalf("expected amount 13.49, got %v", txs[0].Amount)
	}
	if txs[1].Label != "Spotify" {
		t.Fatalf("expected label Spotify, got %s", txs[1].Label)
	}
}

func TestParseCSVWithDebitCreditAndInvalidRows(t *testing.T) {
	csvData := `Date,Description,Debit,Credit
not-a-date,Ignored Row,10,
2024-01-10,Virement interne,0,100
2024-01-12,Electricity,-55.5,
2024-01-15,Payout,,75.25`

	txs, err := ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("expected 2 valid transactions, got %d", len(txs))
	}
	if txs[0].Amount != -55.5 {
		t.Fatalf("expected debit -55.5, got %v", txs[0].Amount)
	}
	if txs[1].Amount != 75.25 {
		t.Fatalf("expected credit 75.25, got %v", txs[1].Amount)
	}
}

func TestParseCSVWithoutHeader(t *testing.T) {
	csvData := "2024-01-01,Service A,15.00\n2024-02-01,Service B,15,00"
	txs, err := ParseCSV(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txs))
	}
	if txs[1].Amount != 15.00 {
		t.Fatalf("expected amount 15.00, got %v", txs[1].Amount)
	}
}
