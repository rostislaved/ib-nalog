package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestParseMoneyMovementSectionBuildsSummaries(t *testing.T) {
	rows := [][]string{
		{"Header", "Валюта", "Дата отчета", "Дата транзакции", "Описание", "Дебет", "Кредит", "Остаток"},
		{"Data", "Отчет по базовой валюте", "2025-01-01", "", "Начальный остаток", "", "", "60.115408656"},
		{"Data", "Отчет по базовой валюте", "2025-01-27", "2025-01-27", "Электронный перевод средств", "", "3147.6", "3207.715408656"},
		{"Data", "Отчет по базовой валюте", "2025-01-28", "2025-01-28", "Покупка 214 SPDR S&P 500 UCITS ETF ACC", "-3140.648928", "", "67.066480656"},
		{"Data", "Отчет по базовой валюте", "2025-12-31", "", "Конечный остаток", "-3140.648928", "3147.6", "67.066480656"},
	}

	section, err := (Sections{}).parseMoneyMovementSection(rows)
	if err != nil {
		t.Fatalf("parse money movement section: %v", err)
	}

	summaries := section.buildSummaries()
	summary, ok := summaries["Отчет по базовой валюте"]
	if !ok {
		t.Fatal("missing base currency summary")
	}

	assertDecimalEqual(t, summary.InitialBalance, "60.115408656")
	assertDecimalEqual(t, summary.TotalCredited, "3147.6")
	assertDecimalEqual(t, summary.TotalDebited, "3140.648928")
	assertDecimalEqual(t, summary.FinalBalance, "67.066480656")
}

func assertDecimalEqual(t *testing.T, actual decimal.Decimal, expected string) {
	t.Helper()

	expectedDecimal, err := decimal.NewFromString(expected)
	if err != nil {
		t.Fatalf("parse expected decimal: %v", err)
	}

	if !actual.Equal(expectedDecimal) {
		t.Fatalf("expected %s, got %s", expectedDecimal, actual)
	}
}
