package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Device interface {
	UpdateOS(version string) error
	GetInfo() string
}

type Smartfone struct {
	OSVersion string
	Model     string
}

type Laptop struct {
	OSVersion string
	Model     string
}

type SmartWatch struct {
	OSVersion string
	Model     string
}

var (
	ErrUnsupported = errors.New("обновление недоступно")
)

func (d *Smartfone) UpdateOS(version string) error {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(d.OSVersion)

	curentVersion, _ := strconv.ParseFloat(match, 64)

	if curentVersion >= 12.0 {
		return ErrUnsupported
	}

	d.OSVersion = version
	return nil
}

func (d *Laptop) UpdateOS(version string) error {
	IsCorectVersion := strings.Contains(version, "Windows")

	if !IsCorectVersion {
		return ErrUnsupported
	}

	d.OSVersion = version
	return nil
}

func (d *SmartWatch) UpdateOS(version string) error {
	count := len([]rune(version))

	if count < 5 {
		return ErrUnsupported
	}

	d.OSVersion = version
	return nil
}

func (d Smartfone) GetInfo() string {

	return getinfo(d.OSVersion, d.Model)

}

func (d Laptop) GetInfo() string {

	return getinfo(d.OSVersion, d.Model)

}

func (d SmartWatch) GetInfo() string {

	return getinfo(d.OSVersion, d.Model)

}

func getinfo(os string, model string) string {
	t := fmt.Sprintf("OS version: %s Model name: %s", os, model)
	return t
}
func main() {

	devices := []Device{
		&Smartfone{OSVersion: "Android 12.0", Model: "Samsung"},
		&Laptop{OSVersion: "Linux Ubuntu", Model: "Lenovo 123123"},
		&SmartWatch{OSVersion: "Fenix 5 pro", Model: "Fuu"},
	}

	newVers := []string{
		"Android 12.3",
		"Window 11",
		"hhhh",
	}

	for i, p := range devices {
		fmt.Println(p.GetInfo())

		err := p.UpdateOS(newVers[i])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(p.GetInfo())
	}

}
