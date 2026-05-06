package main

import (
	"fmt"
)

func Level1() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Паника обработана на уровне 1: ошибка в Level3")
		}
	}()
	Level2()
}

func Level2() {
	Level3()
	defer fmt.Println("Завершаем Level2")
}

func Level3() {
	panic("ошибка в Level3")
}

func main() {

	Level1()
}
