package commands

import (
	c "zirc/client"
	"zirc/helpers"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

// pass - Used to set a connection password before registration.
// Needed when the server requires a password
func pass(params map[string]interface{}) string {
	log.Info().Msg("Running PASS...")

	server_password := helpers.GetEnv("IRC_SERVER_PASSWORD", "")
	if len(server_password) == 0 {
		// ignore the PASS command when no password is configured
		return ""
	}

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return "400 :Unknown error occurred" // ERR_UNKNOWNERROR
	}

	client := _client.(*c.Client)
	if client.GetState("server_password") == "accepted" {
		return "462 :You may not reregister"
	}

	password, ok := params["params"]
	if !ok {
		log.Error().Msg("")
		return "461 :Need more params" //ERR_NEEDMOREPARAMS
	}

	err := bcrypt.CompareHashAndPassword([]byte(server_password), []byte(password.(string)))
	if err != nil {
		log.Error().Msg(err.Error())
		return "ERROR :Closing Link: <username>[<hostname>] (Password incorrect)"
	}

	client.UpdateState("server_password", "accepted")
	log.Info().Msg("Server password accepted")
	// IRC doesn't return anything to the client if the password is correct
	return ""
}
