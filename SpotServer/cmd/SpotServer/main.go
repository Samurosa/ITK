package main

import (
	"ITK_Code/m/v2/internal/config"
	"fmt"
)

func main() {

	cfg, err := config.Load("")
	if err != nil {
		fmt.Println("error load config: ", err)
		return
	}

	fmt.Printf("%#v\n", cfg)

}
