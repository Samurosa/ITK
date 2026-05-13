package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// Реплика БД (имитация)
func dbReplica(name string, in <-chan int) {
	defer wg.Done()
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

func tee(in <-chan int, replicas []chan int) {

	for n := range in {
		for _, channel := range replicas {
			channel <- n
		}
	}

	for _, channel := range replicas {
		close(channel)
	}
}

func main() {
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}
	data := [6]int{1, 2, 3, 4, 5, 6}

	for i, replica := range replicas {

		wg.Add(1)
		go dbReplica(fmt.Sprintf("replica %d", i+1),
			replica,
		)
	}

	go tee(input, replicas)

	go func() {
		defer close(input)
		for _, i := range data {
			input <- i
		}
	}()

	wg.Wait()
}
