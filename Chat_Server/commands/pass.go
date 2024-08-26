package commands

import (
	"zirc/helpers"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

// pass - Used to set a connection password before registration.
// Needed when the server requires a password
func pass(params []string) string {
	msg := "Running PASS..."
	log.Info().Msg(msg)

	if len(params) < 1 {
		// not enough params
		// 461 ERR_NEEDMOREPARAMS
		return "461 ERR_NEEDMOREPARAMS"
	}

	server_password := helpers.GetEnv("IRC_SERVER_PASSWORD", "")
	if len(server_password) == 0 {
		return "Server does not require password"
	}

	password := params[0]
	err := bcrypt.CompareHashAndPassword([]byte(server_password), []byte(password))
	if err != nil {
		log.Error().Msg(err.Error())
		return "ERROR :You need to send your password before registering"
	}

	return ""
}
