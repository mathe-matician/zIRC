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

	split_msg := cmd_re.FindStringSubmatch(cmd_params)
	log.Debug().Msgf("PRIVMSG: split_msg: %s", split_msg)

	split_msg_len := len(split_msg)
	if split_msg_len == 0 {
		return ERR_NEEDMOREPARAMS("")
	}

	// remove match that contains the entire string
	split_msg = split_msg[1:]
	log.Debug().Msgf("PRIVMSG: trimmed split_msg: %s, len: %d", split_msg, len(split_msg))

	if len(split_msg) < 2 || split_msg[0] == "" || split_msg[1] == "" {
		log.Debug().Msgf("PRIVMSG: not enough params")
		return ERR_NEEDMOREPARAMS("")
	}

	target := split_msg[0]
	log.Debug().Msgf("PRIVMSG: msg target: %s", target)

	first_char := string(target[0])
	var channel *Channel
	isChan := false
	if first_char == GENERAL_CHAN_PREFIX || first_char == LOCAL_CHAN_PREFIX || first_char == MODELESS_CHAN_PREFIX {
		log.Debug().Msgf("Message is for channel as prefix is: %s", first_char)
		isChan = true
		// only need to get channel map if sending to channel
		channel_map := g_Server._MessageManager.ChannelMap
		if channel_map == nil {
			log.Error().Msg("Channel list is null!!")
			return ERR_UNKNOWNERROR("")
		}

		log.Debug().Msgf("param_channel: %s", target)
		var valid_channel bool
		channel, valid_channel = (*channel_map)[target]

		if !valid_channel {
			log.Debug().Msgf("So such chan!")
			return ERR_NOSUCHCHANNEL("", cmd_params)
		}

		// TODO
		// do other channel stuff
	}

	msg := split_msg[1]
	msg = strings.TrimLeft(msg, " ") // rm space between channel/nick and message if it exists
	if len(msg) != 0 && string(msg[0]) != ":" {
		log.Debug().Msgf("PRIVMSG: len(msg) != 0 or msg doesn't have : prefix")
		return ERR_NOTEXTTOSEND("", target)
	}
	// trim off single space between target and msg as well as ':' prefix
	msg = msg[1:]

	log.Debug().Msgf("PRIVMSG: target: %s, msg: %s", target, msg)

	// TODO
	// check if they have JOINed channel already
	// TODO
	// rate limit messages

	task_runner := g_Server._MessageManager.Task_runner
	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to PRIVMSG command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	client_details := fmt.Sprintf("%s!%s@%s", client.Nick(), client.User(), client.Ip())
	// TODO
	// channel.Name could be a user name too
	original_msg := msg
	msg = fmt.Sprintf(":%s PRIVMSG %s :%s\r\n", client_details, target, msg)
	// :zach!cloak.z.irc PRIVMSG #test :hi there buddy\r\n
	// :zzirc.comPRIVMSG:marmar@marmar!cloak.z.irc JOIN :#test
	// client_prefix := fmt.Sprintf("%s!%s@%s", client.Nick(), client.User(), g_Server.DnsName)

	var _task *Task
	if isChan {
		log.Debug().Msgf("PRIVMSG: Creating new task for channel: %s", channel.Name)
		_task = NewTask(MULTICAST, msg, 0.0, client.ClientConn, channel, false, client_details)
	} else {
		log.Debug().Msgf("PRIVMSG: Creating new task for USER")

		// have the server find the target by name
		// TODO
		// this doesn't make sense in a multi-server setup
		// TODO
		// check clientservermap first
		// if not on this server
		clientServer, ok := g_Server.ClientServerMap[target]
		if !ok {
			log.Error().Msgf("Could not find client '%s' in client-server map", target)
			return ERR_NOSUCHNICK("")
		}

		var dest *Client
		var remoteServerTask Target
		if clientServer == g_Server.DnsName {
			dest := g_Server._MessageManager.GetClientByNick(target)
			if dest == nil {
				return ERR_NOSUCHNICK("")
			}
		} else {
			// create 'dummy' client obj to use as ref to remote client
			remoteServerTask = &RemoteTask{
				clientServer,
				target,
				msg,
			}
			dest = &Client{nick: target, server: clientServer}
		}

		// TODO
		// handle CTCP
		if strings.HasPrefix(original_msg, ctcpDelimiter) && strings.HasSuffix(original_msg, ctcpDelimiter) {
			log.Debug().Msgf("PRIVMSG: this is a CTCP Msg!!")
			// then this is a CTCP message
			parseCTCP()
		}

		_task = NewTask(UNICAST, msg, 0.0, (*dest).ClientConn, remoteServerTask, false, client_details)
	}

	privmsg_task := []*Task{}
	privmsg_task = append(privmsg_task, _task)

	task_runner <- privmsg_task

	return EMPTY_RESPONSE()
}
