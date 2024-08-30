package commands

import (
	"fmt"

	c "zirc/client"
	t "zirc/task"

	"github.com/phuslu/log"
)

func join(params map[string]interface{}) Response {
	log.Info().Msgf("Join sttart")
	msg := "Running JOIN..."
	log.Info().Msg(msg)

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("no params passed to join")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	log.Info().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*c.Client)

	if len(cmd_params) != 0 {

		// TODO - check if chan exists before doing stuff

		client_details := fmt.Sprintf("%s@%s!%s", client.Nick(), client.User(), client.Ip())

		// TODO - there are different channel prefixs
		chan_prefix := "#"

		msg := fmt.Sprintf(":%s JOIN :%s%s", client_details, chan_prefix, cmd_params)

		_task_runner := params["task_runner"]
		task_runner := _task_runner.(chan []*t.Task)

		// channel := t.NewChannel("#", cmd_params)
		_rpl_topic := RPL_TOPIC("", "fake-topic")
		_rpl_namreply := RPL_NAMREPLY("")
		_rpl_endofnames := RPL_ENDOFNAMES("")

		// src := t.NewSrc(nil)

		join_msg := t.NewTask(t.MULTICAST, msg, 0.0)
		rpl_topic := t.NewTask(t.UNICAST, _rpl_topic.Msg(), 0.0)
		rpl_namreply := t.NewTask(t.UNICAST, _rpl_namreply.Msg(), 0.0)
		rpl_endofnames := t.NewTask(t.UNICAST, _rpl_endofnames.Msg(), 0.0)

		task_runner <- []*t.Task{
			join_msg,
			rpl_topic,
			rpl_namreply,
			rpl_endofnames,
		}
	}

	return EMPTY_RESPONSE()
}
