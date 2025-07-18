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
		client.ClientConn.Close()
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
		client.ClientConn.Close()
		return ERR_NEEDMOREPARAMS("")
	}

	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): args: %v", args)

	split_msg := cmd_re.FindStringSubmatch(args)

	// required numb of args
	if len(split_msg) != 3 {
		log.Error().Msg("CMD(SERVER) not enought args")
		return ERR_NEEDMOREPARAMS("")
	}

	ircServer := g_Server.GetServerByConn(client.ClientConn)
	task_runner := g_Server._MessageManager.Task_runner

	if ircServer == nil {
		// if this server isn't in our server graph
		// then it is attempting to handshake with us
		return handleServerHandshake(client, ircServer, task_runner, split_msg)
	}

	// if this server IS in our server graph
	// then we are getting some update from the server
	// generally it is the server giving us its network topology
	return handleServerNetworkTopology(client, ircServer, task_runner)
}

// handleServerHandshake at a high level does the following:
// 1. Since the server didn't exist in its server graph, it add the server
// 2. Sends a CAPAB command back to the sender to introduce itself
// 2. Sends a PROTOCTL command back to the sender to introduce itself
// 2. Sends a SERVER command back to the sender to introduce itself
func handleServerHandshake(client *Client, ircServer *IrcServer, task_runner chan []*Task, split_msg []string) Response {
	serverName := split_msg[1]
	_hopCount := cmd_re.FindStringSubmatch(strings.Trim(split_msg[2], " "))

	description := cmd_re.FindStringSubmatch(strings.Trim(_hopCount[2], " "))

	hopCount, err := strconv.Atoi(_hopCount[1])
	if err != nil {
		log.Error().Msgf("error converting hopcount to int: %v", err.Error())
		client.ClientConn.Close()
		return ERR_UNKNOWNERROR("")
	}
	serverDescription := strings.Trim(description[2], " ")

	// connState := NULL
	// if ircServer.Conn.State == HANDSHAKING {
	// 	// if the server is currently handshaking, don't reset the connState back to NULL
	// 	connState = HANDSHAKING
	// }

	newServer := IrcServer{
		DnsName:      serverName,
		HopCount:     hopCount,
		Description:  serverDescription,
		Conn:         ServerConn{client.ClientConn, HANDSHAKING},
		_ServerGraph: NewServerGraph(),
		Clients:      make([]*Client, 0),
	}

	// add server to this server's internal graph
	g_Server._ServerGraph.Graph = append(g_Server._ServerGraph.Graph, &newServer)
	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): storing name: %s, hopcount: %d, description: %s", serverName, hopCount, serverDescription)

	//////////////
	// Steps
	//////////////
	// Reply with CAPAB
	// Reply with this server's SERVER command
	//   You must traverse this server's graph in order and send it as a reply to the other server
	//   i.e. starting with self (hopcount 0) send SERVER me

	// 2. send individual SERVER command to THIS SERVER'S _DIRECTLY_ connected servers, and increment the hopcount
	// EXCEPT for the server that is sending this SERVER command
	hopCount++
	msg := fmt.Sprintf("SERVER %s %d :%s \r\n", serverName, hopCount, serverDescription)

	// these need to be sent IN ORDER
	// e.g. need to traverse the graph and send servers in order
	for _, server := range g_Server._ServerGraph.Graph {
		log.Debug().EmbedObject(client).Msg("CMD(SERVER): sending my servers!")
		if server.Conn.Conn == client.ClientConn || server.Conn.Conn == g_Server.Conn.Conn {
			// do not send this server command back to the same server that just sent it to us to prevent SERVER loops in the graph
			log.Debug().EmbedObject(client).Msg("CMD(SERVER): skipping - server is either the conn that initiated the command or is self")
			continue
		}

		log.Debug().Msgf("Sending: %s to %s", msg, server.DnsName)
		if server.Conn.Conn == nil {
			log.Error().Msgf("%s connection is nil for %s", msg, server.DnsName)
			// TODO
			// clean up conn?
			// reconnect? idk
			continue
		}
		server_cmd_tsk := NewTask(SERVER, msg, 0.0, server.Conn.Conn, nil, true, "")
		task_runner <- []*Task{
			server_cmd_tsk,
		}
	}
	return EMPTY_RESPONSE()
}

func handleServerNetworkTopology(client *Client, server *IrcServer, task_runner chan []*Task) Response {
	// otherwise this server already exists
	// and this is the other _type_ of SERVER command
	// and we should be receiving the topology of the sending server (i.e. the server graph)
	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): server already exists in list: %v", server.DnsName)
	// i.e. this is our REPLY back to the sender of the SERVER command
	//      as the sender is already registered as a server we know
	// CAPAB
	// BURST SERVERs
	// UID of users
	// etc

	for _, mysvr := range g_Server._ServerGraph.Graph {
		msg := fmt.Sprintf("SERVER %s %d :%s \r\n", mysvr.DnsName, mysvr.HopCount+1, mysvr.Description)

		log.Debug().EmbedObject(client).Msg("CMD(SERVER): sending my servers!")
		if server.Conn.Conn == client.ClientConn {
			// do not send this server command back to the same server that just sent it to us to prevent SERVER loops in the graph
			log.Debug().EmbedObject(client).Msg("CMD(SERVER): skipping - server is either the conn that initiated the command or is self")
			continue
		}

		log.Debug().Msgf("Sending: %s to %s", msg, mysvr.DnsName)
		if server.Conn.Conn == nil {
			log.Error().Msgf("%s connection is nil", msg, mysvr.DnsName)
			// TODO
			// clean up conn?
			// reconnect? idk
			continue
		}
		server_cmd_tsk := NewTask(SERVER, msg, 0.0, server.Conn.Conn, nil, true, "")
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
