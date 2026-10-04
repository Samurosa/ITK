package main

import (
	"ITK_Code/m/v2/internal/app"
	"flag"
	"fmt"
)

func main() {
	cfgPath := flag.String(
		"config",
		"./internal/config/local.yaml",
		"config file path",
	)

	flag.Parse()

	if err := app.Run(*cfgPath); err != nil {
		fmt.Println("server start error: ", err)
		return
	}
}
