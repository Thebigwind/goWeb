package main

import (
	"fmt"
	"os"

	. "github.com/xtao/goWeb"
)

func main() {
	fmt.Println("httpserver starts now...")

	err := ServerStart()
	if err != nil {
		fmt.Printf("Fail to start httpserver :%s\n", err.Error())
		os.Exit(1)
	}
}
