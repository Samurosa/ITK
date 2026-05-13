package main

import (
	"fmt"
)

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

func generate(nameMetric string, value float64) <-chan ServerMetric {
	serverMetricObject := ServerMetric{Name: nameMetric, Value: value}

	out := make(chan ServerMetric)

	go func() {
		defer close(out)
		out <- serverMetricObject
	}()
	return out
}

func ConvertBytesToMB(in <-chan ServerMetric) <-chan ServerMetric {
	out := make(chan ServerMetric)

	go func() {
		defer close(out)
		for n := range in {
			n.Value /= 1024 * 1024
			out <- n
		}
	}()

	return out
}

func main() {

	inputData := generate("memory_usage", 1024)

	transformResult := ConvertBytesToMB(inputData)

	for result := range transformResult {
		fmt.Printf("Metric name: %s\nMetric value: %f",
			result.Name,
			result.Value,
		)
	}

}
