package main

import (
	"fmt"
)

type Sections struct {
	sections map[string][][]string
}

type currency = string

func (s Sections) getSectionRows(sectionName string) ([][]string, error) {
	section, ok := s.sections[sectionName]
	if !ok {
		return nil, fmt.Errorf("В отчете нет секции: %v", sectionName)
	}

	return section, nil
}

func (s Sections) GetNavSection() (NAVSection, error) {
	rows, err := s.getSectionRows("Чистая стоимость активов (NAV)")
	if err != nil {
		return NAVSection{}, err
	}

	section, err := s.parseNAVSection(rows)
	if err != nil {
		return NAVSection{}, err
	}

	return section, nil
}

func (s Sections) GetNavDeltaSection() (NAVDeltaSection, error) {
	rows, err := s.getSectionRows("Изменение в NAV")
	if err != nil {
		return NAVDeltaSection{}, err
	}

	section, err := s.parseNAVDeltaSection(rows)
	if err != nil {
		return NAVDeltaSection{}, err
	}

	return section, nil
}

func (s Sections) GetMoneySection() (MoneySection, error) {
	rows, err := s.getSectionRows("Отчет о денежных средствах")
	if err != nil {
		return MoneySection{}, err
	}

	section, err := s.parseMoneySection(rows)
	if err != nil {
		return MoneySection{}, err
	}

	return section, nil
}

func (s Sections) GetMoneyMovementSection() (MoneyMovementSection, error) {
	rows, err := s.getSectionRows("Отчет о движении денежных средств")
	if err != nil {
		return MoneyMovementSection{}, err
	}

	section, err := s.parseMoneyMovementSection(rows)
	if err != nil {
		return MoneyMovementSection{}, err
	}

	return section, nil
}
