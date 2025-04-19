package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		panic("no args provided")
	}
	password := args[1]

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	fmt.Printf("Password: %s", password)
	fmt.Printf("Hash: %s", string(hash))
}
