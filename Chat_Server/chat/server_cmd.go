package chat

import "github.com/phuslu/log"

// In the case of a server connection, this command can be used for server-to-server communications (typically not used by clients).
func server(params map[string]interface{}) Response {
	msg := "Running SERVER..."
	log.Info().Msg(msg)
	res := Reply{
		code: "333",
		msg:  msg,
	}

	// parse message
	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to SERVER command!!")
		return ERR_UNKNOWNERROR("")
	}
	log.Debug().Msg("CMD(SERVER): before client cast...")
	client := _client.(*Client) // technically a client object, but this is a server connection

	_params, ok := params["params"]
	args := _params.(string)
	if len(args) == 0 || !ok {
		log.Error().Msg("Params not in map!!")
		return ERR_NEEDMOREPARAMS("")
	}

	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): args: %v", args)

	split_msg := cmd_re.FindStringSubmatch(args)
	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): split_msg: %s", split_msg)

	// 1. add this server to internal routing table with metadata
	// g_Server.RoutingTable.Servers[]

	// 2. forward SERVER command to other servers, but increment the hopcount

	return &res
}
