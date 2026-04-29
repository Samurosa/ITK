package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width float64
	Hight float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) Perimeter() float64 {
	return math.Pi * 2 * c.Radius
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Hight
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Hight)
}

func Calculate(s Shape) {
	fmt.Println(s.Area())
	fmt.Println(s.Perimeter())
}

func main() {
	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 14, Hight: 5}

	Calculate(circle)
	Calculate(rectangle)
}
