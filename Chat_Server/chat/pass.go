package chat

import (
	"net"
	"strings"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	S2S_PASS_ARG_COUNT = 1
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
	passargs, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!!")
		return ERR_NEEDMOREPARAMS("")
	}

	log.Info().Msg("before client cast...")

	client := _client.(*Client)
	if client.GetState("server_password") == "accepted" {
		return ERR_ALREADYREGISTRED("")
	}

	conn := client.ClientConn
	remoteAddr := conn.RemoteAddr().String()
	addr, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Error().Msgf("PASS: Error extracting IP: %s", err)
		return ERR_UNKNOWNERROR("")
	}

	if client.IsServer {
		log.Info().Msgf("Server connection attempted by %s:%s", addr, port)
		// if the client that is connecting is communicating on the S2S port, 7000
		// we need to check whether it is valid to do so
		if !g_Server._ServerManager.whiteListedServerIp(addr) {
			// if this isn't a whitelisted server
			// pass back a terminate msg to terminate the connection
			log.Error().Msgf("PASS(s2s): %s not a whitelisted IP. Terminating connection", remoteAddr)
			return TERMINATE()
		}

		s2s_password := G_Config.S2S.Password_file
		// TODO
		// even if no password
		// we should still do something with the extra args?
		if len(s2s_password) == 0 {
			// ignore the PASS command when no password is configured
			// this SHOULDNT happen with s2s, but possible when testing
			log.Warn().Msgf("PASS(s2s): S2S password not set! This is highly irregular!!")
			return EMPTY_RESPONSE()
		}

		// TODO
		// s2s PASS usually contains more args <protocol version> <flags>
		// A set of capabilities or features supported by the connecting server.
		// Often vendor-specific, depending on the IRC daemon.

		passArgs := strings.Split(passargs.(string), " ")
		passArgsLen := len(passArgs)
		if passArgsLen < S2S_PASS_ARG_COUNT {
			log.Error().Msgf("PASS(s2s): server %s missing required args. Want >= %d, Got %d", remoteAddr, S2S_PASS_ARG_COUNT, passArgsLen)
			return TERMINATE()
		}

		passwordArg := passArgs[0]
		protocolVersionArg := passArgs[1]
		supportedProtos := G_Config.Server.Supported_protocol_versions
		if !strings.Contains(supportedProtos, protocolVersionArg) {
			log.Error().Msgf("PASS(s2s): %s using unsupported protocol version %s. I only support: %s", remoteAddr, protocolVersionArg, supportedProtos)
			return TERMINATE()
		}

		var flagsArg []string
		if passArgsLen > S2S_PASS_ARG_COUNT {
			flagsArg = passArgs[1:]
		}

		// TODO
		// do something with these extra flags
		// Example flags:
		// IRC+ao → May indicate certain features (e.g., +a for account support, +o for oper permissions).
		// TS6 → Used in timestamp-based networks for Nick/Channel collision resolution.
		// NOOPER → Disallows automatic oper permissions.
		// SID → Used for unique Server IDs in some networks.

		log.Debug().Msgf("PASS(s2s): Extra flags passed: %v", flagsArg)

		log.Debug().Msg("before password...")
		err = bcrypt.CompareHashAndPassword([]byte(s2s_password), []byte(passwordArg))
		if err != nil {
			log.Error().Msg(err.Error())
			return ERR_PASSWDMISMATCH("")
		}

		// TODO
		// what should be updated here for s2s state? if any?
		client.AddState("server_password", "accepted")
		log.Info().Msg("S2S Server password accepted")
	} else if (G_Config.Server.Enable_tls && port == G_Config.Server.Tls_port) || port == G_Config.Server.Port {
		log.Info().Msgf("Client connection attempted by %s", remoteAddr)
		// TODO
		// make sure this is the port we want it to be
		// else this is a normal client trying to communicate on the normal client port
		server_password := G_Config.Server.Password_file
		if len(server_password) == 0 {
			// ignore the PASS command when no password is configured
			return EMPTY_RESPONSE()
		}

		log.Info().Msg("before password...")
		err = bcrypt.CompareHashAndPassword([]byte(server_password), []byte(passargs.(string)))
		if err != nil {
			log.Error().Msg(err.Error())
			return ERR_PASSWDMISMATCH("")
		}

		client.AddState("server_password", "accepted")
		log.Info().Msg("Client Server password accepted")
	}

	// IRC doesn't return anything to the client if the password is correct
	return EMPTY_RESPONSE()
}
