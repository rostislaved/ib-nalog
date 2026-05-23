package main

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Изменение в NAV
type NAVDeltaSection struct {
	MarketRevaluation decimal.Decimal
}

func (s Sections) parseNAVDeltaSection(rows [][]string) (NAVDeltaSection, error) {
	var section NAVDeltaSection
	var hasMarketRevaluation bool

	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		if row[0] != "Data" {
			continue
		}

		name := strings.TrimSpace(row[1])
		if name != "Рыночная переоценка" {
			continue
		}

		value, err := parseDecimal(row[2])
		if err != nil {
			return NAVDeltaSection{}, fmt.Errorf("parse market revaluation: %w", err)
		}

		section.MarketRevaluation = value
		hasMarketRevaluation = true
	}

	if !hasMarketRevaluation {
		return NAVDeltaSection{}, fmt.Errorf("missing NAV delta row: Рыночная переоценка")
	}

	return section, nil
}
