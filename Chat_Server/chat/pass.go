package chat

import (
	"net"
	"zirc/helpers"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

// pass - Used to set a connection password before registration.
// Needed when the server requires a password
func pass(params map[string]interface{}) Response {
	log.Info().Msg("Running PASS...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}
	log.Info().Msg("before client cast...")

	client := _client.(*Client)
	if client.GetState("server_password") == "accepted" {
		return ERR_ALREADYREGISTRED("")
	}

	// TODO
	// check server whitelist
	// if connection is from valid server ip

	conn := client.ClientConn
	remoteAddr := (*conn).RemoteAddr().String()
	_, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Error().Msgf("PASS: Error extracting IP: %s", err)
		return ERR_UNKNOWNERROR("")
	}

	if port == "7000" {

	}

	if !g_Server._ServerManager.whiteListedServerIp((*conn).RemoteAddr().String()) {
		// if this isn't a whitelisted server
		// pass back a terminate msg to terminate the connection
		return TERMINATE()
	}

	server_password := helpers.GetEnv("IRC_SERVER_PASSWORD", "")
	if len(server_password) == 0 {
		// ignore the PASS command when no password is configured
		return EMPTY_RESPONSE()
	}

	log.Info().Msg("before params...")

	password, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!!")
		return ERR_NEEDMOREPARAMS("")
	}

	log.Info().Msg("before password...")
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
