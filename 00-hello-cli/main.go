package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Hello %s!\n", os.Args[1])
		return
	}
	fmt.Println("Hello World!")
}