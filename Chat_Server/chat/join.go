package chat

import (
	"fmt"
	"zirc/helpers"

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
		log.Error().Msg("Client not passed to JOIN command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	if len(cmd_params) != 0 {

		// TODO - check if chan exists before doing stuff

		client_details := fmt.Sprintf("%s@%s!%s", client.Nick(), client.User(), client.Ip())

		msg := fmt.Sprintf(":%s JOIN :%s", client_details, cmd_params)

		_task_runner := params["task_runner"]
		task_runner := _task_runner.(chan []*Task)
		_channel_list := params["channel_list"]
		channel_list := _channel_list.(*[]*Channel)
		if channel_list == nil {
			log.Error().Msg("Channel list is null!!")
			return ERR_UNKNOWNERROR("")
		}

		chan_prefix := string(cmd_params[0])
		var duration string
		if chan_prefix == "#" {
			duration = "persistent"
		} else if chan_prefix == "&" || chan_prefix == "+" {
			// temporary channels and modeless channels are transient in nature
			// when the last client leaves they are deleted
			duration = "temporary"
		} else {
			return ERR_BADCHANMASK("")
		}

		channel := NewChannel(
			cmd_params,
			"",
			"",
			"",
			"active",
			duration,
		)

		// Add channel to channel list
		// Add channel to client?
		// should there just be a single source of truth?
		// when do we need to have both lists kept up to date?
		// is it just a burden to update both?
		// AND updating the client's list itself?
		(*channel_list) = append((*channel_list), channel)

		// TODO - update current user as the channel operator

		// server_task := NewTask(t.SERVER, fmt.Sprintf("create chan %s", cmd_params), 0.0)

		_server_metadata := params["server_metadata"]
		server_metadata := _server_metadata.(map[string]string)

		server_name := server_metadata["name"]

		_rpl_topic := RPL_TOPIC("", "")
		_rpl_topic_msg := helpers.FormatResponse(server_name, _rpl_topic.Code(), _rpl_topic.Msg())
		_rpl_namreply := RPL_NAMREPLY("")
		_rpl_namreply_msg := helpers.FormatResponse(server_name, _rpl_namreply.Code(), _rpl_namreply.Msg())
		_rpl_endofnames := RPL_ENDOFNAMES("")
		_rpl_endofnames_msg := helpers.FormatResponse(server_name, _rpl_endofnames.Code(), _rpl_endofnames.Msg())

		join_msg := NewTask(MULTICAST, msg, 0.0, nil)
		rpl_topic := NewTask(UNICAST, string(_rpl_topic_msg), 0.0, client.ClientConn)
		rpl_namreply := NewTask(UNICAST, string(_rpl_namreply_msg), 0.0, client.ClientConn)
		rpl_endofnames := NewTask(UNICAST, string(_rpl_endofnames_msg), 0.0, client.ClientConn)

		task_runner <- []*Task{
			join_msg,
			rpl_topic,
			rpl_namreply,
			rpl_endofnames,
		}
	}

	return EMPTY_RESPONSE()
}
