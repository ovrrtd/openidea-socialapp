package main

import (
	"fmt"
	"socialapp/cmd"

	_ "github.com/lib/pq"
)

func main() {
	err := cmd.Server()

	if err != nil {
		fmt.Printf("Server error: %s", err.Error())
	}
}
