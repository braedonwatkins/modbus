package main

import (
	"log"
)

func main() {
	err := RunServer()
	if err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
