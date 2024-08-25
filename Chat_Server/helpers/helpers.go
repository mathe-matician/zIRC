package helpers

import (
	"errors"
	"os"
)

var server_modes = []string{
	"leaf",
	"hybrid",
	"hub",
}

// VerifyServerMode checks whether the mode passed in is a valid mode
// one of - leaf, hub, hybrid (leaf + hub):
//
//	leaf: handles client connections
//	hub: relays messages between leafs
//	hybrid: leaf that can relay messages
func VerifyServerMode(mode string) error {
	for _, v := range server_modes {
		if mode == v {
			return nil
		}
	}
	return errors.New("server mode not recognized")
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
