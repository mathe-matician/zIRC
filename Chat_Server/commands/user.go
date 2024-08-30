package commands

import (
	c "zirc/client"

	"github.com/phuslu/log"
)

func user(params map[string]interface{}) Response {
	log.Info().Msg("Running USER...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to USER command!!")
		return ERR_UNKNOWNERROR("")
	}

	user, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	client := _client.(*c.Client)
	client.SetUser(user.(string))

	res := EMPTY_RESPONSE()

	// if len(client.Nick()) != 0 && !client.Registered {
	// 	client.Registered = true
	// 	responses := []Response{
	// 		RPL_WELCOME(""),
	// 		RPL_YOURHOST(""),
	// 		RPL_CREATED(""),
	// 		RPL_MYINFO(""),
	// 	}
	// 	WriteMultipleResponses(responses, client.ClientConn)
	// }

	return res
}
