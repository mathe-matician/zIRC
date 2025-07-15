package chat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phuslu/log"
)

// server command can be thought of as us as a server receiving a SERVER command from another server NOT sending it / syncing a new server
// In the case of a server connection, this command can be used for server-to-server communications (typically not used by clients).
//
// Note: connection whitelisting has already taken place if we are running this command
func server(params map[string]interface{}) Response {
	log.Info().Msg("Running SERVER...")

	// parse message
	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to SERVER command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client) // technically a client object, but this is a server connection

	if !client.IsServer {
		log.Error().Msg("Client isn't server")
		return ERR_UNKNOWNERROR("")
	}

	// params should include in this order:
	// 1. server name
	// 2. hopcount
	// 3. server description
	_params, ok := params["params"]
	args := _params.(string)
	if len(args) == 0 || !ok {
		log.Error().Msg("Params not in map!!")
		return ERR_NEEDMOREPARAMS("")
	}

	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): args: %v", args)

	split_msg := cmd_re.FindStringSubmatch(args)

	// required numb of args
	if len(split_msg) != 3 {
		log.Error().Msg("CMD(SERVER) not enought args")
		return ERR_NEEDMOREPARAMS("")
	}

	serverName := split_msg[1]
	_hopCount := cmd_re.FindStringSubmatch(strings.Trim(split_msg[2], " "))

	description := cmd_re.FindStringSubmatch(strings.Trim(_hopCount[2], " "))

	hopCount, err := strconv.Atoi(_hopCount[1])
	if err != nil {
		log.Error().Msgf("error converting hopcount to int: %v", err.Error())
		return ERR_UNKNOWNERROR("")
	}
	serverDescription := strings.Trim(description[2], " ")

	srvr := g_Server.GetServerByConn(client.ClientConn)
	// if no server connection is found in this server's Server list
	// then add it
	if srvr == nil {
		log.Debug().EmbedObject(client).Msgf("CMD(SERVER): storing name: %s, hopcount: %d, description: %s", serverName, hopCount, serverDescription)
		newServer := IrcServer{
			DnsName:     serverName,
			HopCount:    hopCount,
			Description: serverDescription,
			Conn:        client.ClientConn,
			Servers:     make([]*IrcServer, 0),
			Clients:     make([]*Client, 0),
		}

		g_Server.Servers = append(g_Server.Servers, &newServer)
	} else {
		// otherwise this server already exists
		log.Debug().EmbedObject(client).Msgf("CMD(SERVER): server already exists in list: %v", srvr.DnsName)
	}

	task_runner := g_Server._MessageManager.Task_runner

	// 2. send individual SERVER command to THIS SERVER'S _DIRECTLY_ connected servers, and increment the hopcount
	// EXCEPT for the server that is sending this SERVER command
	hopCount++
	msg := fmt.Sprintf("SERVER %s %d :%s \r\n", serverName, hopCount, serverDescription)

	for _, server := range g_Server.Servers {
		log.Debug().EmbedObject(client).Msg("CMD(SERVER): sending my servers!")
		if server.Conn == client.ClientConn {
			// do not send this server command back to the same server that just sent it to us to prevent SERVER loops in the graph
			log.Debug().EmbedObject(client).Msg("CMD(SERVER): skipping server that initiated SERVER cmd...")
			continue
		}

		log.Debug().Msgf("Sending: %s to %s", msg, server.DnsName)
		if server.Conn == nil {
			log.Error().Msgf("%s connection is nil", msg, server.DnsName)
			// TODO
			// clean up conn?
			// reconnect? idk
			continue
		}
		server_cmd_tsk := NewTask(SERVER, msg, 0.0, server.Conn, nil, true, "")
		task_runner <- []*Task{
			server_cmd_tsk,
		}
	}

	// TODO
	// RPL_WELCOME needs to be sent via task manager as the custom server connection code doesn't follow the
	// same return path as all the other commands
	// return RPL_WELCOME("", client.server)
	return EMPTY_RESPONSE()
}
