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
	log.Debug().Msg("CMD(SERVER): before client cast...")
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
	// something is messed up here with this regex
	// THis command:
	// SERVER ZZ.IRC 1 :super fun zzirc server
	// shows up as this in this log:
	// ZZ.IRC 1 :super fun zzirc server ZZ.IRC  1 :super fun zzirc server
	// which then throws off the index parsing of Atoi below
	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): split_msg: %s", split_msg)

	// required numb of args
	if len(split_msg) != 3 {
		log.Error().Msg("CMD(SERVER) not enought args")
		return ERR_NEEDMOREPARAMS("")
	}

	serverName := split_msg[1]
	// STOPPED
	// ZZ.IRC 1 :super fun zzirc server ZZ.IRC  1 :super fun zzirc server
	// do I need to rerun FindStringSubmatch multiple times on this to wittle down the args?
	// still hitting "{"time":"2025-07-13T22:13:03.802Z","level":"error","message":"error converting hopcount to int: strconv.Atoi: parsing \" 1 :super fun zzirc server\": invalid syntax"}"
	// where strconv.Atoi(split_msg[2]) using index 2
	// probably print out what each index is
	_hopCount := cmd_re.FindStringSubmatch(strings.Trim(split_msg[2], " "))
	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): _hopCount: %s", _hopCount)

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
	}

	task_runner := g_Server._MessageManager.Task_runner

	// 2. send individual SERVER command to THIS SERVER'S _DIRECTLY_ connected servers, and increment the hopcount
	// EXCEPT for the server that is sending this SERVER command
	hopCount++
	msg := fmt.Sprintf("SERVER %s %d :%s \r\n", serverName, hopCount, serverDescription)

	for _, server := range g_Server.Servers {
		if server.Conn == client.ClientConn {
			// do not send this server command back to the same server that just sent it to us to prevent SERVER loops in the graph
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
