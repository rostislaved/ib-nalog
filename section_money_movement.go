package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Отчет о движении денежных средств
type MoneyMovementSection struct {
	ByCurrency         map[currency]MoneyMovementReport
	BaseCurrencyReport MoneyMovementReport
}

type MoneyMovementReport struct {
	InitialBalance MoneyMovementRow
	FinalBalance   MoneyMovementRow
	Rows           []MoneyMovementRow
}

type MoneyMovementRow struct {
	Currency        currency
	ReportDate      *time.Time
	TransactionDate *time.Time
	Description     string
	Debit           decimal.Decimal
	Credit          decimal.Decimal
	Balance         decimal.Decimal
}

func (s Sections) parseMoneyMovementSection(rows [][]string) (MoneyMovementSection, error) {
	mms := MoneyMovementSection{
		ByCurrency: make(map[currency]MoneyMovementReport),
	}
	var hasRows bool

	for _, row := range rows {
		if len(row) < 8 {
			continue
		}

		if strings.TrimSpace(row[0]) != "Data" {
			continue
		}

		c, r, err := parseMoneyMovementRow(row)
		if err != nil {
			return MoneyMovementSection{}, err
		}

		report := mms.ByCurrency[c]
		report.Rows = append(report.Rows, r)

		switch r.Description {
		case "Начальный остаток":
			report.InitialBalance = r
		case "Конечный остаток":
			report.FinalBalance = r
		}

		mms.ByCurrency[c] = report
		hasRows = true
	}

	if !hasRows {
		return MoneyMovementSection{}, fmt.Errorf("missing money movement rows")
	}

	baseCurrencyReport, ok := mms.ByCurrency["Отчет по базовой валюте"]
	if ok {
		mms.BaseCurrencyReport = baseCurrencyReport
	}

	return mms, nil
}

func parseMoneyMovementRow(row []string) (currency, MoneyMovementRow, error) {
	c := currency(strings.TrimSpace(row[1]))

	reportDate, err := parseOptionalReportDate(row[2])
	if err != nil {
		return "", MoneyMovementRow{}, err
	}

	transactionDate, err := parseOptionalReportDate(row[3])
	if err != nil {
		return "", MoneyMovementRow{}, err
	}

	debit, err := parseOptionalDecimal(row[5])
	if err != nil {
		return "", MoneyMovementRow{}, err
	}

	credit, err := parseOptionalDecimal(row[6])
	if err != nil {
		return "", MoneyMovementRow{}, err
	}

	balance, err := parseDecimal(row[7])
	if err != nil {
		return "", MoneyMovementRow{}, err
	}

	r := MoneyMovementRow{
		Currency:        c,
		ReportDate:      reportDate,
		TransactionDate: transactionDate,
		Description:     strings.TrimSpace(row[4]),
		Debit:           debit,
		Credit:          credit,
		Balance:         balance,
	}

	return c, r, nil
}

func parseOptionalDecimal(value string) (decimal.Decimal, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return decimal.Zero, nil
	}

	d, err := parseDecimal(value)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return d, nil
}

func parseReportDate(value string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func parseOptionalReportDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	t, err := parseReportDate(value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

type MoneyMovementSummary struct {
	InitialBalance decimal.Decimal
	TotalCredited  decimal.Decimal
	TotalDebited   decimal.Decimal
	FinalBalance   decimal.Decimal
}

type MoneyMovementSummaryByCurrency struct {
	Currency currency
	Summary  MoneyMovementSummary
}

func (mms MoneyMovementSection) buildSummaries() map[currency]MoneyMovementSummary {
	summaries := make(map[currency]MoneyMovementSummary, len(mms.ByCurrency))

	for c, report := range mms.ByCurrency {
		summary := MoneyMovementSummary{
			InitialBalance: report.InitialBalance.Balance,
			FinalBalance:   report.FinalBalance.Balance,
		}

		row, ok := getRow(report.Rows, "Начальный остаток")
		if !ok {
			panic("Нет строки")
		}

		summary.InitialBalance = row.Balance

		row, ok = getRow(report.Rows, "Конечный остаток")
		if !ok {
			panic("Нет строки")
		}

		summary.TotalCredited = row.Credit
		summary.FinalBalance = row.Balance

		// summary.TotalDebited = row.Debit // Если считать так, то калькулятор на nalog.ru скажет, что не сходится
		summary.TotalDebited = summary.InitialBalance.Sub(summary.FinalBalance).Add(summary.TotalCredited)

		summaries[c] = summary
	}

	return summaries
}

func getRow(rows []MoneyMovementRow, name string) (MoneyMovementRow, bool) {
	for _, row := range rows {
		if row.Description == name {
			return row, true
		}
	}

	return MoneyMovementRow{}, false

}

func (mms MoneyMovementSection) sortedSummaries() []MoneyMovementSummaryByCurrency {
	summaries := mms.buildSummaries()
	currencies := make([]currency, 0, len(summaries))

	for c := range summaries {
		if c == "Отчет по базовой валюте" {
			continue
		}

		currencies = append(currencies, c)
	}

	sort.Strings(currencies)

	sortedSummaries := make([]MoneyMovementSummaryByCurrency, 0, len(currencies))
	for _, c := range currencies {
		sortedSummary := MoneyMovementSummaryByCurrency{
			Currency: c,
			Summary:  summaries[c],
		}

		sortedSummaries = append(sortedSummaries, sortedSummary)
	}

	return sortedSummaries
}
