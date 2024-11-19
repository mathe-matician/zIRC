package helpers

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"os"
	"syscall"
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

func FormatResponse(response_args ...string) []byte {
	response := ":"
	for _, val := range response_args {
		response += val
		response += " "
	}
	response = response[:len(response)-1]
	response += "\r\n"
	return []byte(response)
}

// CreateTLSConfig creates a tls.Config
// using the passed crt and key paths
func CreateTLSConfig(crt_path, key_path string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(crt_path, key_path)
	if err != nil {
		panic(err)
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

// https://stackoverflow.com/a/73029344
func IsNetConnClosedErr(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}
