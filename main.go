package main

import (
	"fmt"
	"go-controller/core"
)

func main() {
	fmt.Println("12313213")
	writer := core.NewWriter()
	reader := core.NewReader()
	controller := core.NewController(writer, reader)
	controller.Run()
}
