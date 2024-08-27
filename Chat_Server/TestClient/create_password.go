package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, err := bcrypt.GenerateFromPassword([]byte("1235"), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	fmt.Printf("Hash: %s", string(hash))
}
