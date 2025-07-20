package main

import (
	"fmt"
	"os"

	"github.com/meh-hackathon/meh/auth"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: auth <password>")
		fmt.Println("  Hashes the provided password and prints the result")
		os.Exit(1)
	}

	hash, err := auth.HashPassword(os.Args[1])
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(hash)
}
