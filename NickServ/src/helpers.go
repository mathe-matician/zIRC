package nickserv

import (
	"crypto/tls"
	"net"
	"os"
	"strconv"

	"github.com/phuslu/log"
)

// CreateTLSConfig creates a tls.Config
// using the passed crt and key paths
func CreateTLSConfig(crt_path, key_path string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(crt_path, key_path)
	if err != nil {
		panic(err)
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// GetTCPListener checks to see if TLS is enabled for the particular listener
// if it is, it creates a TLS Listener with the provided TLS env vars
// else it returns a normal Listener
func GetTCPListener(enableTls, tls_cert_path, tls_key_path, tls_port, port string) net.Listener {
	var ln net.Listener
	var err error
	enable_tls, err1 := strconv.ParseBool(enableTls)
	if err1 != nil {
		log.Error().Msg(err1.Error())
	}
	irc_host := GetEnv("IRC_HOST", "0.0.0.0")
	if enable_tls {
		config := CreateTLSConfig(tls_cert_path, tls_key_path)
		tls_port := tls_port
		ln, err = tls.Listen("tcp", irc_host+":"+tls_port, config)
		if err != nil {
			panic(err)
		}
	} else {
		ln, err = net.Listen("tcp", irc_host+":"+port)
		if err != nil {
			log.Error().Msg(err.Error())
			panic(err.Error())
		}
	}
	return ln
}
