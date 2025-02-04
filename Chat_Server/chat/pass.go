package chat

import (
	"net"

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
	addr, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Error().Msgf("PASS: Error extracting IP: %s", err)
		return ERR_UNKNOWNERROR("")
	}

	// TODO
	// use config it makes it easier

	if (G_Config.S2S.Enable_tls && port == G_Config.S2S.Tls_port) || port == G_Config.S2S.Port {
		log.Info().Msgf("Server connection attempted by %s:%s", addr, port)
		// if the client that is connecting is communicating on the S2S port, 7000
		// we need to check whether it is valid to do so
		if !g_Server._ServerManager.whiteListedServerIp(addr) {
			// if this isn't a whitelisted server
			// pass back a terminate msg to terminate the connection
			return TERMINATE()
		}

		s2s_password := G_Config.S2S.Password_file
		if len(s2s_password) == 0 {
			// ignore the PASS command when no password is configured
			return TERMINATE()
		}
	} else if (G_Config.Server.Enable_tls && port == G_Config.Server.Tls_port) || port == G_Config.Server.Port {
		log.Info().Msgf("Client connection attempted by %s:%s", addr, port)
		// TODO
		// make sure this is the port we want it to be
		//
		// else this is a normal client trying to communicate on the normal client port
	}

	server_password := G_Config.Server.Password_file
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
	err = bcrypt.CompareHashAndPassword([]byte(server_password), []byte(password.(string)))
	if err != nil {
		log.Error().Msg(err.Error())
		return ERR_PASSWDMISMATCH("")
	}

	client.UpdateState("server_password", "accepted")
	log.Info().Msg("Server password accepted")
	// IRC doesn't return anything to the client if the password is correct
	return EMPTY_RESPONSE()
}
