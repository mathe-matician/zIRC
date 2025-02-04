package chat

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/phuslu/log"
)

var SUPPORTED_AUTH_TYPES = G_Config.Server.Supported_auth_types

func ValidAuthNType(authType string) bool {
	authTypes := strings.Split(SUPPORTED_AUTH_TYPES, ",")
	for _, i := range authTypes {
		if authType == i {
			return true
		}
	}
	return false
}

func ValidAuthZType() {

}

// authenticate uses sasl cap
// allows a user to reserve a nick IF they authenticate
// if other users try to connect and use that nick for a registered user
// then they can't unless they can authenticate
func authenticate(params map[string]interface{}) Response {
	msg := "Running AUTHENTICATE..."
	log.Info().Msg(msg)

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("no params passed to join")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	log.Info().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to JOIN command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)
	client_Auth := &client.Auth.AuthenticationState
	_, valid_authType := (*client_Auth)[cmd_params]

	if !valid_authType {
		log.Error().Msgf("Not a valid auth type: %s", cmd_params)
		return ERR_NOSASL("")
	}

	var res Response
	if cmd_params == "PLAIN" {
		res = plain(client_Auth, cmd_params, client)
	} else if cmd_params == "SCRAM-SHA-256" {
		res = scramSha265(client_Auth, cmd_params)
	} else if cmd_params == "OAUTHBEARER" {
		res = oAuthBearer(client_Auth, cmd_params)
	} else if cmd_params == "EXTERNAL" {
		res = external(client_Auth, cmd_params)
	} else {
		// just return the available auth methods
		return ERR_NOSASL(fmt.Sprintf(":%s", SUPPORTED_AUTH_TYPES))
	}

	return res
}

// Authentication Identity: Typically the username or identifier of the user logging in.
// Authorization Identity: Optional and usually left empty (often represented as ""), which means the user is authorizing as themselves.
// Password: The password associated with the username.
// <authzid>\0<authcid>\0<password>
// <authzid> is the authorization identity (optional, usually empty).
//
//	used when a user can "assume" another identity. useful for admins / bots
//
// <authcid> is the authentication identity (username/nick).
// <password> is the password.

var auth_step_count = map[string]int{
	"PLAIN":         5,
	"SCRAM-SHA-256": 5,
	"OAUTHBEARER":   5,
}

func plain(auth_step *map[string]int, params string, client *Client) Response {
	log.Debug().Msg("Running PLAIN authentication")
	var res Response

	if (*auth_step)["PLAIN"] == 1 {
		(*auth_step)["PLAIN"]++
		res = &Reply{
			code: "",
			msg:  "AUTHENTICATE +",
		}
	} else if (*auth_step)["PLAIN"] == 2 {
		// base64 decode
		decodedBytes, err := base64.StdEncoding.DecodeString(params)
		if err != nil {
			log.Error().Msgf("Error decoding PLAIN auth: %s", err)
			return ERR_SASLFAIL("SASL Failed")
		}

		authString := string(decodedBytes)
		log.Debug().Msgf("PLAIN auth string decoded: %s", authString)

		// split by "\0"
		splitAuthString := strings.Split(authString, "\\0")
		if len(splitAuthString) < 3 {
			log.Error().Msgf("Len of auth string is less than 3")
			return ERR_SASLFAIL("SASL Failed")
		}

		// get authzid, authcid, and password
		authzid := splitAuthString[0]
		authcid := splitAuthString[1]
		// password := splitAuthString[2]

		if authzid == "" {
			// if the user isn't authorizing as someone else
			// they usually just send an empty string, which then authorizes them as their passed user
			authzid = authcid
		}

		// check if valid authcid
		// var _username, _password string
		// _, err = plain_auth_stmt.QueryOne(pg.Scan(&_username, &_password), authcid)

		// if err != nil {
		// 	log.Error().Msgf("Error executing PLAIN auth query: %s", err)
		// 	return ERR_SASLFAIL("SASL Failed")
		// }

		// check if they can assume this role authzid

		// finally if we get here, the client's authentication/authorization process was successful
		client.Auth.IsAuthenticated = true
		client.Auth.AuthorizationType = AuthZType(authzid)  // (assumed role)
		client.Auth.AuthenticationType = AuthNType(authcid) // (assumed role)
	}

	return res
}

func scramSha265(auth_step *map[string]int, params string) Response {
	log.Debug().Msg("Running SCRAM-SHA-256 authentication")
	return ERR_NOSASL(fmt.Sprintf(":%s", SUPPORTED_AUTH_TYPES))
}

func oAuthBearer(auth_step *map[string]int, params string) Response {
	log.Debug().Msg("Running OAUTHBEARER authentication")
	return ERR_NOSASL(fmt.Sprintf(":%s", SUPPORTED_AUTH_TYPES))
}

// conn.ConnectionState().PeerCertificates[0]
func external(auth_step *map[string]int, params string) Response {
	log.Debug().Msg("Running EXTERNAL authentication")
	return ERR_NOSASL(fmt.Sprintf(":%s", SUPPORTED_AUTH_TYPES))
}
