package main

import (
	"github.com/shopspring/decimal"
)

// Отчет о денежных средствах
type MoneySection struct {
	m                  map[currency]map[string]oneRow
	BaseCurrencyReport BaseCurrencyReport
}

type BaseCurrencyReport struct {
	InitialAmount                oneRow
	DepositAmount                oneRow
	Dividends                    oneRow
	BrokerFeePaidAndReceived     oneRow
	SellTrades                   oneRow
	BuyTrades                    oneRow
	EndPeriodAmount              oneRow
	Commissions                  oneRow
	WithholdingTax               oneRow
	CurrencyConversionProfitLoss oneRow
	FinalSettlementAmount        oneRow
}

func (s Sections) parseMoneySection(rows [][]string) (MoneySection, error) {
	// m := make(map[string]map[string]oneRow)
	mms := MoneySection{
		m: make(map[string]map[string]oneRow),
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		type_ := row[2]

		mm, ok := mms.m[type_]
		if !ok {
			mm = make(map[string]oneRow)
		}

		total, err := parseDecimal(row[3])
		if err != nil {
			return MoneySection{}, err
		}

		securities, err := parseDecimal(row[4])
		if err != nil {
			return MoneySection{}, err
		}

		futures, err := parseDecimal(row[5])
		if err != nil {
			return MoneySection{}, err
		}

		r := oneRow{
			Total:      total,
			Securities: securities,
			Futures:    futures,
		}

		name := row[1]
		mm[name] = r

		mms.m[type_] = mm
	}

	for currency, value := range mms.m {
		switch currency {
		case "Отчет по базовой валюте":

			bcr := BaseCurrencyReport{
				InitialAmount:                value["Начальная сумма средств"],
				DepositAmount:                value["Внесение средств"],
				Dividends:                    value["Дивиденды"],
				BrokerFeePaidAndReceived:     value["Ставка брокера: уплачено и получено"],
				SellTrades:                   value["Сделки (продажа)"],
				BuyTrades:                    value["Сделки (покупка)"],
				EndPeriodAmount:              value["Остаток средств на конец периода"],
				Commissions:                  value["Комиссии"],
				WithholdingTax:               value["Удерживаемый налог"],
				CurrencyConversionProfitLoss: value["Прибыль/убытки при пересчете валюты"],
				FinalSettlementAmount:        value["Конечная расчетная сумма средств"],
			}

			mms.BaseCurrencyReport = bcr
		default:
			continue

		}

	}

	return mms, nil
}

type oneRow struct {
	Total      decimal.Decimal
	Securities decimal.Decimal
	Futures    decimal.Decimal
}
