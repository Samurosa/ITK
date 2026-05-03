package main

import (
	"errors"
	"fmt"
)

var (
	ErrEngineAlreadyRunning = errors.New("двигатель уже работает")
	ErrEngineOff            = errors.New("двигатель не запущен")
	ErrLowBattery           = errors.New("низкий заряд батареи")
)

type Vehicle interface {
	StartEngine() error
	StopEngine() error
	GetInfo() string
}

type Car struct {
	Brand    string
	EngineOn bool
}

func (c Car) Honk() string {
	return "Beep beep!"
}

func (c Car) GetEngineStatus() bool {
	return c.EngineOn
}

func (c *Car) StartEngine() error {
	if c.EngineOn == true {
		return ErrEngineAlreadyRunning
	}
	c.EngineOn = true
	return nil
}
func (c *Car) StopEngine() error {
	if c.EngineOn == false {
		return ErrEngineOff
	}
	c.EngineOn = false
	return nil
}
func (c Car) GetInfo() string {
	return fmt.Sprintf("Марка машины: %s Состояние двигателя: %t", c.Brand, c.EngineOn)
}

type Truck struct {
	Car
	CargoCapacity float64
}

func (t Truck) Honk() string {
	return "Honk Honk!"
}

func (t Truck) GetCargoCapacity() float64 {
	return t.CargoCapacity
}

type ElectricCar struct {
	Car
	BatteryLevel int
}

func (ecar *ElectricCar) StartEngine() error {
	if ecar.BatteryLevel < 5 {
		return ErrLowBattery
	}
	ecar.EngineOn = true
	return nil
}

func (ecar ElectricCar) GetBatteryLevel() int {
	return ecar.BatteryLevel
}

func main() {

	// ===== Car =====
	car := &Car{Brand: "Lada"}

	err := car.StartEngine()
	fmt.Println("Car start:", err == nil, car.GetEngineStatus())

	err = car.StartEngine()
	fmt.Println("Car double start error:", err == ErrEngineAlreadyRunning)

	err = car.StopEngine()
	fmt.Println("Car stop:", err == nil, !car.GetEngineStatus())

	err = car.StopEngine()
	fmt.Println("Car double stop error:", err == ErrEngineOff)

	fmt.Println("Car honk:", car.Honk())

	// ===== Truck =====
	truck := &Truck{
		Car:           Car{Brand: "Sitrak"},
		CargoCapacity: 20,
	}

	fmt.Println("Truck honk:", truck.Honk())
	fmt.Println("Truck capacity:", truck.GetCargoCapacity())

	// ===== ElectricCar =====
	tesla := &ElectricCar{
		Car:          Car{Brand: "Tesla"},
		BatteryLevel: 3,
	}

	err = tesla.StartEngine()
	fmt.Println("Electric low battery error:", err == ErrLowBattery)

	tesla.BatteryLevel = 10
	err = tesla.StartEngine()
	fmt.Println("Electric start success:", err == nil, tesla.GetEngineStatus())

	// ===== Полиморфизм =====
	vehicles := []Vehicle{
		&Car{Brand: "BMW"},
		&Truck{Car: Car{Brand: "MAN"}},
		&ElectricCar{Car: Car{Brand: "Tesla"}, BatteryLevel: 10},
	}

	for _, v := range vehicles {
		err := v.StartEngine()
		fmt.Println("Start via interface:", err == nil)
		fmt.Println(v.GetInfo())
	}
}
