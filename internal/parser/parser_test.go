package parser

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParsePDFSeparatorAndHeader(t *testing.T) {
	pdfData := buildPDF(t, []string{
		"Date;Libellé;Montant",
		"2024-01-01;Netflix;13,49",
		"2024-02-01;Spotify;9.99",
	})

	txs, err := ParsePDF(bytes.NewReader(pdfData))
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

func TestParsePDFWithDebitCreditAndInvalidRows(t *testing.T) {
	pdfData := buildPDF(t, []string{
		"Date,Description,Debit,Credit",
		"not-a-date,Ignored Row,10,",
		"2024-01-10,Virement interne,0,100",
		"2024-01-12,Electricity,-55.5,",
		"2024-01-15,Payout,,75.25",
	})

	txs, err := ParsePDF(bytes.NewReader(pdfData))
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

func TestParsePDFWithoutHeader(t *testing.T) {
	pdfData := buildPDF(t, []string{
		"2024-01-01 ServiceA 15.00",
		"2024-02-01 ServiceB 15,00",
	})

	txs, err := ParsePDF(bytes.NewReader(pdfData))
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

func TestParsePDFWithShuffledColumns(t *testing.T) {
	pdfData := buildPDF(t, []string{
		"Libellé;Crédit;Date;Débit",
		"Prime Cashback;15,00;01/02/2024;",
		"Netflix;;2024-03-01;13.49",
	})

	txs, err := ParsePDF(bytes.NewReader(pdfData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txs))
	}
	if txs[0].Amount != 15.00 || txs[1].Amount != -13.49 {
		t.Fatalf("unexpected amounts: %+v", txs)
	}
}

func TestParsePDFLooseLines(t *testing.T) {
	pdfData := buildPDF(t, []string{
		"Some header text",
		"Netflix -13,49 prélèvement 01/01/2024",
		"02-02-2024 Spotify 9.99 paiement",
		"ignored line without amount",
	})

	txs, err := ParsePDF(bytes.NewReader(pdfData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(txs))
	}
	if txs[0].Label == "" || txs[1].Label == "" {
		t.Fatalf("labels should not be empty: %+v", txs)
	}
}

func buildPDF(t *testing.T, lines []string) []byte {
	t.Helper()
	var builder strings.Builder
	builder.WriteString("%PDF-1.4\n")
	builder.WriteString("1 0 obj <<>> endobj\n")
	builder.WriteString("2 0 obj <<>> stream\n")
	builder.WriteString("BT /F1 12 Tf 72 712 Td\n")
	for i, line := range lines {
		builder.WriteString(fmt.Sprintf("(%s) Tj\n", line))
		if i < len(lines)-1 {
			builder.WriteString("T*\n")
		}
	}
	builder.WriteString("ET\nendstream\nendobj\ntrailer <<>>\n%%EOF")

	return []byte(builder.String())
}
