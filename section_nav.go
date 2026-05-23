package main

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Чистая стоимость активов (NAV)
type NAVSection struct {
	StockStart decimal.Decimal
	StockEnd   decimal.Decimal
}

func (s Sections) parseNAVSection(rows [][]string) (NAVSection, error) {
	var section NAVSection
	var hasStock bool

	for _, row := range rows {
		if len(row) < 6 {
			continue
		}
		if row[0] != "Data" {
			continue
		}

		assetClass := strings.TrimSpace(row[1])
		if assetClass != "Акция" {
			continue
		}

		stockStart, err := parseDecimal(row[2])
		if err != nil {
			return NAVSection{}, fmt.Errorf("parse stock start: %w", err)
		}

		stockEnd, err := parseDecimal(row[5])
		if err != nil {
			return NAVSection{}, fmt.Errorf("parse stock end: %w", err)
		}

		section.StockStart = stockStart
		section.StockEnd = stockEnd
		hasStock = true
	}

	if !hasStock {
		return NAVSection{}, fmt.Errorf("missing NAV row: Акция")
	}

	return section, nil
}
