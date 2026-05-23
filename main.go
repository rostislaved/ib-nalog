package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shopspring/decimal"
)

type Statement struct {
	navSection           NAVSection
	navDeltaSection      NAVDeltaSection
	moneySection         MoneySection
	moneyMovementSection MoneyMovementSection
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: ibkr-form <statement.csv>")
		os.Exit(2)
	}

	sections, err := readSections(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	statement, err := parseStatement(sections)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	print(statement)
}

func readSections(path string) (Sections, error) {
	file, err := os.Open(path)
	if err != nil {
		return Sections{}, fmt.Errorf("open csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	sections := make(map[string][][]string)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return Sections{}, fmt.Errorf("read csv: %w", err)
		}
		if len(record) == 0 {
			continue
		}

		key := strings.TrimSpace(record[0])
		if key == "" {
			continue
		}

		values := make([]string, len(record)-1)
		copy(values, record[1:])

		sections[key] = append(sections[key], values)
	}

	return Sections{sections: sections}, nil
}

func parseStatement(sections Sections) (Statement, error) {
	navSection, err := sections.GetNavSection()
	if err != nil {
		return Statement{}, err
	}

	navDeltaSection, err := sections.GetNavDeltaSection()
	if err != nil {
		return Statement{}, err
	}

	moneySection, err := sections.GetMoneySection()
	if err != nil {
		return Statement{}, err
	}

	moneyMovementSection, err := sections.GetMoneyMovementSection()
	if err != nil {
		return Statement{}, err
	}

	statement := Statement{
		navSection:           navSection,
		navDeltaSection:      navDeltaSection,
		moneySection:         moneySection,
		moneyMovementSection: moneyMovementSection,
	}

	return statement, nil
}

func print(statement Statement) {
	for _, item := range statement.moneyMovementSection.sortedSummaries() {
		fmt.Printf(
			"%s\n| %v | %v | %v | %v |\n\n",
			item.Currency,
			item.Summary.InitialBalance.StringFixed(2),
			item.Summary.TotalCredited.StringFixed(2),
			item.Summary.TotalDebited.StringFixed(2),
			item.Summary.FinalBalance.StringFixed(2),
		)
	}

	fmt.Println("Акции")
	fmt.Printf(
		"| %s | %s | %s | %s |\n",
		statement.navSection.StockStart.StringFixed(2),
		statement.moneySection.BaseCurrencyReport.BuyTrades.Total.Abs().Add(statement.navDeltaSection.MarketRevaluation.Abs()).StringFixed(2),
		statement.moneySection.BaseCurrencyReport.SellTrades.Total.Abs().Add(statement.moneySection.BaseCurrencyReport.CurrencyConversionProfitLoss.Total.Abs()).StringFixed(2),
		statement.navSection.StockEnd.StringFixed(2),
	)
}

func parseDecimal(value string) (decimal.Decimal, error) {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")

	parsed, err := decimal.NewFromString(cleaned)
	if err != nil {
		return decimal.Zero, err
	}

	return parsed, nil
}
