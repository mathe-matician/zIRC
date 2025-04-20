package chat

import "github.com/phuslu/log"

// In the case of a server connection, this command can be used for server-to-server communications (typically not used by clients).
func server(params map[string]interface{}) Response {
	log.Info().Msg("Running SERVER...")

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

	// required numb of args
	if len(split_msg) != 3 {
		log.Error().Msg("CMD(SERVER) not enought args")
		return ERR_NEEDMOREPARAMS("")
	}

	serverName := split_msg[0]
	hopCount := split_msg[1]
	// serverDescription := split_msg[2]

	// 1. add this server to internal routing table with metadata
	// need to know "who" (what server) is sending this SERVER command, as that is considered its parent or "next hop" in the routing table
	// could be found with client.ClientConn via IP
	nextHop := "direct"
	if hopCount != "1" {
		// set it to client.server as if this client is a server, then it would have gone through this command previously
		// and we would have set client.server = serverName
		nextHop = client.server
	} else {
		client.server = serverName
	}

	g_Server.RoutingTable.Servers[serverName] = nextHop

	// 2. forward SERVER command to other servers, but increment the hopcount

	return RPL_WELCOME("", client.server)
}
