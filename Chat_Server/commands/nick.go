package commands

import (
	c "zirc/client"

	"github.com/phuslu/log"
)

func nick(params map[string]interface{}) *map[string]string {
	log.Info().Msg("Running NICK...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}

	// when modifying old nickname
	// :oldnickname!username@hostname NICK :newnickname

	// server then broadcasts :Alice!alice@192.0.2.1 NICK :Alicia
	// to any channel this client is part of
	// it is broadcasted to all private converstaions this client is in
	// as well as all channels this client is in

	nick, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	client := _client.(*c.Client)
	client.SetNick(nick.(string))

	// no response from NICK signals success
	return EMPTY_RESPONSE()
}
