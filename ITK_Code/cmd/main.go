package main

import (
	"fmt"
)

type pair struct {
	word   string
	amount int
}

var pairs []pair

// FilterByValue возвращает новую map, содержащую только элементы,
// значения которых присутствуют в allowedValues.
func FilterByValue(m map[int]string, allowedValues []string) map[int]string {
	// Преобразовать allowedValues в set для быстрой проверки
	whiteList := make(map[string]struct{})
	for _, a := range allowedValues {
		whiteList[a] = struct{}{}
	}
	// Создать новую map и заполнить её подходящими элементами
	n := make(map[int]string)

	for key, value := range m {
		if _, ok := whiteList[value]; ok {
			n[key] = value
		}
	}
	return n
}

// InvertMap меняет ключи и значения местами.
// Если значения исходной map не уникальны, возвращает ошибку.
func InvertMap(m map[string]int) (map[int]string, error) {
	// Проверять уникальность значений
	// При обнаружении дубликата вернуть ошибку с описанием конфликта
	invetredMap := make(map[int]string)
	for key, value := range m {
		if _, ok := invetredMap[value]; ok {
			return nil, fmt.Errorf("duplicate value: %d", value)
		}

		invetredMap[value] = key
	}

	return invetredMap, nil
}
