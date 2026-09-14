package chat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phuslu/log"
)

type HandshakeArgs struct {
	Name        string
	Description string
	HopCount    int
}

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

	split_msg := cmd_re.FindStringSubmatch(args)

	// required numb of args
	if len(split_msg) != 3 {
		log.Error().Msg("CMD(SERVER) not enought args")
		return ERR_NEEDMOREPARAMS("")
	}

	serverNode := g_Server.Servers.GetServerByConn(client.ClientConn)
	task_runner := g_Server._MessageManager.Task_runner

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

	handshakeArgs := HandshakeArgs{
		Name:        serverName,
		HopCount:    hopCount,
		Description: serverDescription,
	}

	// if serverNode == nil && g_Server.Servers.GetPendingByConn(client.ClientConn) == nil {
	// 	//
	// 	return handleServerHandshake(client, task_runner, handshakeArgs)
	// }

	if serverNode == nil {
		//
		return handleServerHandshake(client, task_runner, handshakeArgs)
	}

	// if this server IS in our server tree
	// then we are getting some update from the server
	// generally it is the server giving us its network topology during BURST_SEND
	// though, it can probably be any update to the network topology
	return handleServerNetworkTopology(client, serverNode, task_runner, handshakeArgs)
}

// handleServerHandshake at a high level does the following:
// 1. Since the server didn't exist in its server tree, it add the server
// 2. Sends a CAPAB command back to the sender to introduce itself
// 2. Sends a PROTOCTL command back to the sender to introduce itself
// 2. Sends a SERVER command back to the sender to introduce itself
func handleServerHandshake(client *Client, task_runner chan []*Task, serverArgs HandshakeArgs) Response {
	log.Debug().Msgf("Starting Server handshake")

	// add server to our server table
	// since we are handshaking this server is directly reaching out to us
	// this means that this server is their parent and we will be directly connected
	serverName := serverArgs.Name
	hopCount := serverArgs.HopCount
	serverDescription := serverArgs.Description

	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): storing name: %s, hopcount: %d, description: %s", serverName, hopCount, serverDescription)
	// 1. Add the server to server tree
	svrNode, ok := g_Server.Servers.Pending[client.ClientConn]
	if ok {
		// if found in pending servers
		// remove it as we will now put it into the server Tree since we have the SID and name
		delete(g_Server.Servers.Pending, client.ClientConn)
		svrNode.SID = "" // TODO update SID
		svrNode.Name = serverName
	} else {
		svrNode = NewServerNode(
			WithSID(""),
			WithName(serverName),
			WithDirectlyConnected(true),
			WithIsMe(false),
			WithDescription(serverDescription),
			WithHopCount(hopCount),
			WithServers(nil),
			WithPingPongChan(),
			WithConn(&ServerConn{client.ClientConn, HANDSHAKING}),
			WithParent(g_Server.Servers.Tree[g_Server.Name]),
		)
	}

	g_Server.Servers.Tree[serverName] = svrNode
	g_Server.Servers.Tree[g_Server.Name].Servers = append(g_Server.Servers.Tree[g_Server.Name].Servers, svrNode)
	log.Debug().Msgf("Existing map: %v", g_Server.Servers.Tree)

	// 2. TODO Reply with CAPAB / PROTOCTL
	// 3. Reply with this server's SERVER command
	//   You must traverse this server's graph in order and send it as a reply to the other server
	//   i.e. starting with self (hopcount 0) send SERVER me

	// 4. Reply with all of my known servers via SERVER commands
	//    these must be sent in tree order
	// this server is now going to be receving the burst from us
	svrNode.Conn.State = BURST_RECV
	g_Server.Servers.ServerBurstSend(g_Server.Servers.Tree[g_Server.Name], client.ClientConn)
	// --------

	// PING is what "signals" that this server has completed its burst phase
	// at most this will block this connection's goroutine for whatever the pingpong timeout duration is set to
	ping_s2s(svrNode, true)

	log.Debug().Msg("handshake complete")
	return EMPTY_RESPONSE()
}

func handleServerNetworkTopology(client *Client, server *ServerNode, task_runner chan []*Task, serverArgs HandshakeArgs) Response {
	// otherwise this server already exists
	// and this is the other _type_ of SERVER command
	// and we should be receiving the topology of the sending server (i.e. the server graph)
	log.Debug().Msg("handleServerNetworkTopology()")
	log.Debug().EmbedObject(client).Msgf("CMD(SERVER): server already exists in list: %v", server.Name)
	// i.e. this is our REPLY back to the sender of the SERVER command
	//      as the sender is already registered as a server we know
	// CAPAB
	// BURST SERVERs
	// UID of users
	// etc

	if server.Conn.State == REGISTERED {
		log.Debug().Msg("Server is already s2s registered. This must be an update...")
		// if we are already registered with this server, i.e. we already went through connection / handshake phase
		// then either:
		//   - something went wrong and we have stale state for this server
		//   - this server is sending an update SERVER command after the handshake

		// idk if we need to do something special at this point like set the state to
		// something new or reuse BURST_RECV
	} else {
		server.Conn.State = BURST_SEND
	}

	svrNode := NewServerNode(
		WithSID(""),
		WithName(serverArgs.Name),
		WithDirectlyConnected(false),
		WithIsMe(false),
		WithDescription(serverArgs.Description),
		WithHopCount(serverArgs.HopCount),
		WithServers(nil),
		// WithConn(&ServerConn{nil, HANDSHAKING}),
		WithParent(g_Server.Servers.Tree[g_Server.Name]),
	)

	g_Server.Servers.Tree[serverArgs.Name] = svrNode
	// TODO
	// infer what node to add svrNode as a child?

	for _, myitem := range g_Server.Servers.Tree {
		msg := fmt.Sprintf("SERVER %s %d :%s \r\n", myitem.Server.DnsName, myitem.Server.HopCount, myitem.Server.Description)

		log.Debug().EmbedObject(client).Msg("CMD(SERVER): sending my servers!")
		if server.Conn.Conn == client.ClientConn {
			// do not send this server command back to the same server that just sent it to us to prevent SERVER loops in the graph
			log.Debug().EmbedObject(client).Msg("CMD(SERVER): skipping - server is either the conn that initiated the command or is self")
			continue
		}

		log.Debug().Msgf("Sending: %s to %s", msg, myitem.Server.DnsName)
		if server.Conn.Conn == nil {
			log.Error().Msgf("%s connection is nil for message: %s", myitem.Server.DnsName, msg)
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
