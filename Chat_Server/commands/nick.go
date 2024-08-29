package commands

import (
	"fmt"
	c "zirc/client"

	"github.com/phuslu/log"
)

func nick(params map[string]interface{}) Response {
	log.Debug().Msg("Running NICK...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*c.Client)

	nick, ok := params["params"]
	if !ok || len(nick.(string)) == 0 {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	// TODO - check whether nick is already taken

	// no response from NICK signals success
	res := EMPTY_RESPONSE()

	current_client_nick := client.Nick()
	// check if client is registered, i.e. they have nick AND user set
	// if they haven't completed registration, there isn't a need to send them
	// this special update of their old nickname, and especially send it to other
	if len(current_client_nick) != 0 && len(client.User()) != 0 {
		// when modifying old nickname
		res.MsgOverride(fmt.Sprintf(":%s!%s@%s NICK :%s \r\n", current_client_nick, client.User(), client.Ip(), nick))
		// :oldnickname!username@hostname NICK :newnickname

		// server then broadcasts :Alice!alice@192.0.2.1 NICK :Alicia
		// to any channel this client is part of
		// it is broadcasted to all private converstaions this client is in
		// as well as all channels this client is in
	}

	client.SetNick(nick.(string))

	return res
}
