package chat

import (
	"fmt"

	"github.com/phuslu/log"
)

// TODO:
// Limit characters in NICK
// Alphanumeric Characters: Letters (A-Z, a-z) and digits (0-9) are generally allowed.
// Certain Special Characters: Commonly allowed special characters include - (hyphen), _ (underscore), and . (dot).
// Disallow: Leading/Trailing Characters: Characters like - and . might be restricted from being at the beginning or end of a nickname.
// Length restriction: Most servers have a maximum length for nicknames, commonly between 9 and 30 characters.
// Nicks are case sensitive

func nick(params map[string]interface{}) Response {
	log.Debug().Msg("Running NICK...")

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	nick, ok := params["params"]
	if nick == nil || !ok || len(nick.(string)) == 0 {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	// TODO - check if valid nickname
	// TODO - check whether nick is already taken

	// no response from NICK signals success
	res := EMPTY_RESPONSE()

	current_client_nick := client.Nick()
	// check if client is registered, i.e. they have nick AND user set
	// if they haven't completed registration, there isn't a need to send them
	// this special update of their old nickname, and especially send it to other
	if len(current_client_nick) != 0 && len(client.User()) != 0 {
		// when modifying old nickname
		new_nick := nick.(string)
		nickExists := g_Server.NickExists(new_nick)
		if nickExists {
			return ERR_NICKNAMEINUSE("")
		}

		client.SetNick(new_nick)
		res.MsgOverride(fmt.Sprintf(":%s!%s@%s NICK :%s\r\n", current_client_nick, client.User(), client.Ip(), new_nick))

		delete(g_Server.ClientServerMap, client.nick)
		g_Server.ClientServerMap[new_nick] = g_Server.DnsName

		// TODO
		// :oldnickname!username@hostname NICK :newnickname
		// server then broadcasts :Alice!alice@192.0.2.1 NICK :Alicia
		// to any channel this client is part of
		// it is broadcasted to all private converstaions this client is in
		// as well as all channels this client is in

	} else {
		// add NICK to state
		// so that when registration is complete
		// the NICK will be committed to user
		client.AddState("NICK", nick.(string))
	}

	if len(client.User()) != 0 && !client.Registered {
		client.Registered = true

		nickState := client.GetState("NICK")
		if nickState == "" {
			log.Error().EmbedObject(client).Msgf("NICK state empty?!")
			return ERR_UNKNOWNERROR("")
		}

		client.SetNick(nickState)
		client.RemoveState("NICK")

		// TODO - need to have access to the server manager to get info
		// 		  for these messages...
		//			do we pass the SM into all commands or return these responses to do further formatting
		//			within the ProcessMessage which has access to the SM?
		//			OR do we pass this off to the task_runner which is passed to all command Fns?
		// actually might be easier to pass to task_runner
		//	we can prefix with a special denotation to do special processing within workers...
		// the issue with that is that this is concurrent, so the client can technically send something
		// again before they get their response
		// it could also send them randomly out of order if sent to different workers
		// so it woud need to support an array of tasks

		client_nick := client.Nick()

		server_name := g_Server.DnsName
		server_version := g_Server.Version
		server_usermodes := G_Config.Server.User_modes
		server_channelmodes := G_Config.Server.Channel_modes
		server_creation_date := g_Server.CreationDate.String()

		client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())

		// add the now registered client to the Server's client map
		msg_mang := g_Server._MessageManager
		clint_map := msg_mang.ClientMap
		(*clint_map)[client_nick] = client

		responses := WELCOME_WRAPPER(
			client.ClientConn,
			server_name,
			server_version,
			server_creation_date,
			server_usermodes,
			server_channelmodes,
			client_nick,
			client_details,
		)

		task_runner := g_Server._MessageManager.Task_runner
		task_runner <- responses

		g_Server.ClientServerMap[client.nick] = g_Server.DnsName

		go ping(client)
	}

	return res
}
