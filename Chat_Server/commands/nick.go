package commands

import (
	"fmt"
	c "zirc/client"
	"zirc/helpers"
	t "zirc/task"

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
	if nick == nil || !ok || len(nick.(string)) == 0 {
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

	if len(client.User()) != 0 && !client.Registered {
		client.Registered = true

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
		_server_metadata := params["server_metadata"]
		server_metadata := _server_metadata.(map[string]string)

		server_name := server_metadata["name"]
		server_version := server_metadata["version"]
		server_usermodes := server_metadata["usermodes"]
		server_channelmodes := server_metadata["channelmodes"]

		user_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())
		_001 := string(helpers.FormatResponse(server_name, "001", client_nick, RPL_WELCOME("", client_nick).Msg(), user_details))
		_002 := string(helpers.FormatResponse(server_name, "002", client_nick, RPL_YOURHOST("", server_name, server_version).Msg()))
		_003 := string(helpers.FormatResponse(server_name, "003", client_nick, RPL_CREATED("", server_metadata["date"]).Msg()))
		_rpl_myinfo := RPL_MYINFO("", client_nick, server_name, server_version, server_usermodes, server_channelmodes).Msg()
		_004 := string(helpers.FormatResponse(server_name, "004", client_nick, _rpl_myinfo))

		rpl_welcome := t.NewTask(t.UNICAST, _001, 0.0)
		rpl_yourhost := t.NewTask(t.UNICAST, _002, 0.0)
		rpl_created := t.NewTask(t.UNICAST, _003, 0.0)
		rpl_myinfo := t.NewTask(t.UNICAST, _004, 0.0)

		responses := []*t.Task{
			rpl_welcome,
			rpl_yourhost,
			rpl_created,
			rpl_myinfo,
		}

		_task_runner := params["task_runner"]
		task_runner := _task_runner.(chan []*t.Task)
		task_runner <- responses
	}

	return res
}
