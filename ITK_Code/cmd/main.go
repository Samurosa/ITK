package main

import (
	"fmt"
	"sort"
	"strings"
)

type pair struct {
	word   string
	amount int
}

var pairs []pair

// WordFrequency принимает строку текста и возвращает map с частотой слов.
func WordFrequency(text string) map[string]int {
	value := strings.Fields(text)
	counter := make(map[string]int)
	for _, v := range value {
		counter[v]++
	}
	return counter
}

// PrintWordFrequency выводит частотный анализ слов, отсортированный по убыванию частоты.
func PrintWordFrequency(freqMap map[string]int) {
	// TODO: Реализуйте функцию.

	for key, count := range freqMap {
		pairs = append(pairs, pair{word: key, amount: count})

	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].amount > pairs[j].amount
	})

	for _, p := range pairs {
		fmt.Println(p.word, p.amount)
	}
}

func main() {

	text := "golang is great and golang is fast"

	PrintWordFrequency(WordFrequency(text))

}
