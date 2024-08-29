package commands

import (
	c "zirc/client"
	"zirc/helpers"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

// pass - Used to set a connection password before registration.
// Needed when the server requires a password
func pass(params map[string]interface{}) Response {
	log.Info().Msg("Running PASS...")

	server_password := helpers.GetEnv("IRC_SERVER_PASSWORD", "")
	if len(server_password) == 0 {
		// ignore the PASS command when no password is configured
		return EMPTY_RESPONSE()
	}

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}

	client := _client.(*c.Client)
	if client.GetState("server_password") == "accepted" {
		return ERR_ALREADYREGISTRED("")
	}

	password, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!!")
		return ERR_NEEDMOREPARAMS("")
	}

	err := bcrypt.CompareHashAndPassword([]byte(server_password), []byte(password.(string)))
	if err != nil {
		log.Error().Msg(err.Error())
		return ERR_PASSWDMISMATCH("")
	}

	client.UpdateState("server_password", "accepted")
	log.Info().Msg("Server password accepted")
	// IRC doesn't return anything to the client if the password is correct
	return EMPTY_RESPONSE()
}
