package chat

import (
	"fmt"
	"strings"

	"github.com/phuslu/log"
)

func privmsg(params map[string]interface{}) Response {
	msgg := "Running PRIVMSG..."
	log.Info().Msg(msgg)

	_cmd_params := params["params"]
	if _cmd_params == nil {
		log.Debug().Msgf("no params passed to MODE")
		return ERR_NEEDMOREPARAMS("")
	}

	cmd_params := _cmd_params.(string)
	log.Info().Msgf("cmd_params: %s, len: %d", cmd_params, len(cmd_params))

	if len(cmd_params) == 0 {
		log.Debug().Msgf("no params passed to MODE")
		return ERR_NEEDMOREPARAMS("")
	}

	_channel_map := params["channel_map"]
	channel_map := _channel_map.(*map[string]*Channel)
	if channel_map == nil {
		log.Error().Msg("Channel list is null!!")
		return ERR_UNKNOWNERROR("")
	}

	chan_prefix := string(cmd_params[0])
	if chan_prefix != GENERAL_CHAN_PREFIX && chan_prefix != LOCAL_CHAN_PREFIX && chan_prefix != MODELESS_CHAN_PREFIX {
		return ERR_BADCHANMASK("")
	}

	// split_p := strings.Split(cmd_params, " ")
	split_msg := cmd_re.FindStringSubmatch(cmd_params)
	log.Debug().Msgf("PRIVMSG: split_msg: %s", split_msg)

	split_msg_len := len(split_msg)
	if split_msg_len == 0 {
		return ERR_NEEDMOREPARAMS("")
	}

	// remove match that contains the entire string
	split_msg = split_msg[1:]
	log.Debug().Msgf("PRIVMSG: trimmed split_msg: %s, len: %d", split_msg, len(split_msg))

	// what SHOULD be left is the channel and the msg

	if len(split_msg) < 2 || split_msg[0] == "" || split_msg[1] == "" {
		log.Debug().Msgf("PRIVMSG: not enough params")
		return ERR_NEEDMOREPARAMS("")
	}

	// TODO
	// could be channel OR user
	param_channel := split_msg[0]
	msg := split_msg[1]
	msg = strings.TrimLeft(msg, " ") // rm space between channel/nick and message if it exists
	if len(msg) != 0 && string(msg[0]) != ":" {
		log.Debug().Msgf("PRIVMSG: len(msg) != 0 or msg doesn't have : prefix")
		return ERR_NEEDMOREPARAMS("")
	}
	// trim off single space between channel and msg as well as ':' prefix
	msg = msg[1:]

	log.Debug().Msgf("PRIVMSG: channel: %s, msg: %s", param_channel, msg)

	log.Debug().Msgf("param_channel: %s", param_channel)
	channel, valid_channel := (*channel_map)[param_channel]

	if !valid_channel {
		log.Debug().Msgf("So such chan!")
		return ERR_NOSUCHCHANNEL("", cmd_params)
	}

	// TODO
	// check if they have JOINed channel already
	// TODO
	// rate limit messages

	_task_runner := params["task_runner"]
	task_runner := _task_runner.(chan []*Task)
	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to MODE command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	client_details := fmt.Sprintf("%s@%s!%s", client.Nick(), client.User(), client.Ip())
	// TODO
	// channel.Name could be a user name too
	msg = fmt.Sprintf(":%s PRIVMSG %s :%s \r\n", client_details, channel.Name, msg)

	_task := NewTask(MULTICAST, msg, 0.0, client.ClientConn, channel)
	privmsg_task := []*Task{}
	privmsg_task = append(privmsg_task, _task)

	task_runner <- privmsg_task

	return EMPTY_RESPONSE()
}
