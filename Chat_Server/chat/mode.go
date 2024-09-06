package chat

import (
	"fmt"
	"strings"

	"github.com/phuslu/log"
)

// Channel Modes
//Visibility and Access:
// +p (Private): The channel is not visible in the channel list.
// +s (Secret): The channel is hidden from public view and channel lists.
// +i (Invite-Only): Users must be invited to join the channel.
// Moderation:
// +m (Moderated): Only users with voice (+v) or operator status can speak.
// +n (No External Messages): Prevents users outside the channel from sending messages to it.
// +t (Topic Protection): Only operators can change the topic.
// User Limits and Restrictions:
// +l (Limit): Sets a maximum number of users allowed in the channel.
// +k (Keyed): Requires a password (key) to join the channel.
// Bans and Exemptions:
// +b (Ban List): Bans specific users or masks from the channel.
// +e (Ban Exemption): Exempts specific users from being affected by a ban.
// +P persistent: to persist the channel after all users leave
// +v Voice
// +I invite exception Allows specific users (using a mask) to join an invite-only (+i) channel without needing an invitation.
// +r Registered Channel (with ChanServ) Indicates that the channel is registered, typically with services like ChanServ. Some servers use this mode to show that a channel has been formally registered.
// +z Secure only (SSL/TLS) connections only
// +R Registerd users only (through NickServ)
// +M (Moderated for Unregistered Users) Only registered users (e.g., users identified by NickServ) can send messages to the channel. Unregistered users can join but cannot talk.
// +C (No CTCP) Blocks CTCP (Client-To-Client Protocol) messages, which are often used for things like requesting information from another client or performing actions like /me.
// +c (No Color) Prevents users from using colored text in the channel. This is often used to reduce spam or unwanted formatting.
// +a (Admin) Grants admin status to a user, typically a level between operator (+o) and owner (+q). Admins have significant control but might not have all the powers of the channel owner.
// +q (Owner/Quiet) On some IRC networks, +q indicates channel ownership, giving the user ultimate control over the channel. On others, +q is used to prevent a user from sending messages to the channel (they can see the messages but cannot contribute).
//

// MODE <chan> <+|-><see above> [params to mode]
// + means add - means remove
func mode(params map[string]interface{}) Response {
	log.Info().Msgf("MODE start")
	tmp := "Running MODE..."
	log.Info().Msg(tmp)

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

	split_p := strings.Split(cmd_params, " ")

	if len(split_p) == 0 {
		log.Debug().Msgf("no params passed to MODE")
		return ERR_NEEDMOREPARAMS("")
	}

	// split_mode := re.FindAllStringSubmatch(cmd_params, -1)
	// if split_mode == nil {
	// 	log.Debug().Msgf("no params passed to MODE")
	// 	return ERR_NEEDMOREPARAMS("")
	// }

	// param_channel := split_mode[0]
	param_channel := split_p[0]
	log.Debug().Msgf("param_channel: %s", param_channel)
	// channel, valid_channel := (*channel_map)[param_channel[0]]
	channel, valid_channel := (*channel_map)[param_channel]
	if !valid_channel {
		log.Debug().Msgf("So such chan!")
		return ERR_NOSUCHCHANNEL("", cmd_params)
	}

	// pair off channel
	split_p = split_p[1:]

	// log.Debug().Msgf("split_mode before len(split_mode) == 1: %s", split_mode)
	log.Debug().Msgf("split_p before len(split_p) == 1: %s", split_p)

	// if len(split_mode) == 1 {
	if len(split_p) == 0 {
		// e.g. regular users can run MODE #chan w/o being an operator
		// only when they pass modes are they denied
		return RPL_CHANNELMODEIS("", channel.Name, channel.ChannelModes)
	}

	// get params after channel
	// sm := split_mode[1]
	// log.Debug().Msgf("sm: %s", sm)
	// at this point split_params should have <modes> [params]

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to MODE command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	// check if they are channel operator
	// only chan operators can apply or remove modes
	client_nick := client.Nick()
	if _, channel_operator := channel.Operators[client_nick]; !channel_operator {
		log.Error().Msg("Client isn't channel operator")
		return ERR_NOTONCHANNEL("", channel.Name)
	}

	_server_metadata := params["server_metadata"]
	server_metadata := _server_metadata.(map[string]string)
	_task_runner := params["task_runner"]
	task_runner := _task_runner.(chan []*Task)

	// KEEP THIS
	// PARSE modes / params
	// res := parseModes(split_mode[1], server_metadata["channelmodes"], task_runner)
	// log.Info().Msgf("split_mode[0]: %s, split_mode[1]: %s", split_mode[0], split_mode[1])

	// _mode := split_p[0]
	// if

	// DESIGN:
	// get all modes presented, could be anything - store in variable
	// pair off modes from other params. now at the front of params will be the actual params
	// iterate through paired off modes and have an if block
	// the if block has a toggle that checks to see if we are in adding mode or removing mode
	// have two funcs that add or remove modes, but allow you to pass in params to them
	// during the loop, check to see if valid mode, check to see if + or -, if + or minus encountered
	//		then change the `if` toggle.
	// TODO - what about +vk+ke+jd - is this valid w/ multiple +?

	// var mode_re = regexp.MustCompile(`^\S*`) // captures until first space
	supported_channel_modes := server_metadata["channelmodes"]
	mode_task := []*Task{}
	for _, m := range split_p {
		action := string(m[0])
		log.Debug().Msgf("_mode: %s, action: %s", m, action)
		if len(action) == 0 || (action != "+" && action != "-") || len(m) == 1 {
			log.Debug().Msgf("action empty, action != + or -, m len == 1")
			// if no action is specified, we don't know whether to add or remove the modes
			return ERR_NEEDMOREPARAMS("")
		}

		for _, mode := range m {
			str_mode := string(mode)
			if !strings.Contains(supported_channel_modes, str_mode) {
				unknown_mode_task := NewTask(UNICAST, ERR_UNKNOWNMODE("", str_mode).Msg(), 0.0, client.ClientConn, nil)
				mode_task = append(mode_task, unknown_mode_task)
				continue
			}

			if str_mode == "+" || str_mode == "-" {
				action = str_mode
				continue
			}

			if action == "+" {
				// TODO addMode to user if applicable?
				addMode(str_mode, &channel.ChannelModes)
			} else if action == "-" {
				removeMode(str_mode, &channel.ChannelModes)
			}

			// TODO
			// I think these mode updates to channels are multicasts
			// a lot of them are, but are ALL of them?

			// e.g.
			// :Bob!bob@host MODE #example +i
			client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())
			msg := fmt.Sprintf(":%s MODE %s %s%s", client_details, channel.Name, action, str_mode)
			valid_mode_task := NewTask(MULTICAST, msg, 0.0, client.ClientConn, channel)
			mode_task = append(mode_task, valid_mode_task)
		}
	}

	task_runner <- mode_task

	return EMPTY_RESPONSE()
}

// adds the mode to the channel
// and to the user?
func addMode(mode string, channel_modes *string) {
	*channel_modes += mode
}

// removes the mode from the channel
// and from the user?
func removeMode(mode string, channel_modes *string) {
	*channel_modes = strings.Replace(*channel_modes, mode, "", 1)
}

////// MODES

// this then makes two sources of truth for modes - the ServerManager or IrcServer and this one
// annoying to update both (also not like modes change frequently, though)
// var mode_fns = map[string]func(params *map[string]interface{}){
// 	"p": p,
// 	"s": s,
// 	"i": i,
// 	"m": m,
// 	"n": n,
// 	"t": t,
// 	"l": l,
// 	"k": k,
// 	"b": b,
// 	"e": e,
// 	"P": P,
// 	"v": v,
// 	"I": I,
// 	"r": r,
// 	"R": R,
// 	"z": z,
// 	"M": M,
// 	"c": c,
// 	"C": C,
// 	"a": a,
// 	"q": q,
// }

// +p (Private): The channel is not visible in the channel list.
func p(params *map[string]interface{}) {

}

// +s (Secret): The channel is hidden from public view and channel lists.
func s(params *map[string]interface{}) {

}

// +i (Invite-Only): Users must be invited to join the channel.
func i(params *map[string]interface{}) {

}

// Moderation:
// +m (Moderated): Only users with voice (+v) or operator status can speak.
func m(params *map[string]interface{}) {

}

// +n (No External Messages): Prevents users outside the channel from sending messages to it.
func n(params *map[string]interface{}) {

}

// +t (Topic Protection): Only operators can change the topic.
func t(params *map[string]interface{}) {

}

// User Limits and Restrictions:
// +l (Limit): Sets a maximum number of users allowed in the channel.
func l(params *map[string]interface{}) {

}

// +k (Keyed): Requires a password (key) to join the channel.
func k(params *map[string]interface{}) {

}

// Bans and Exemptions:
// +b (Ban List): Bans specific users or masks from the channel.
func b(params *map[string]interface{}) {

}

// +e (Ban Exemption): Exempts specific users from being affected by a ban.
func e(params *map[string]interface{}) {

}

// +P persistent: to persist the channel after all users leave
func P(params *map[string]interface{}) {

}

// +v Voice
func v(params *map[string]interface{}) {

}

// +I invite exception Allows specific users (using a mask) to join an invite-only (+i) channel without needing an invitation.
func I(params *map[string]interface{}) {

}

// +r Registered Channel (with ChanServ) Indicates that the channel is registered, typically with services like ChanServ. Some servers use this mode to show that a channel has been formally registered.
func r(params *map[string]interface{}) {

}

// +z Secure only (SSL/TLS) connections only
func z(params *map[string]interface{}) {

}

// +R Registerd users only (through NickServ)
func R(params *map[string]interface{}) {

}

// +M (Moderated for Unregistered Users) Only registered users (e.g., users identified by NickServ) can send messages to the channel. Unregistered users can join but cannot talk.
func M(params *map[string]interface{}) {

}

// +C (No CTCP) Blocks CTCP (Client-To-Client Protocol) messages, which are often used for things like requesting information from another client or performing actions like /me.
func C(params *map[string]interface{}) {

}

// +c (No Color) Prevents users from using colored text in the channel. This is often used to reduce spam or unwanted formatting.
func c(params *map[string]interface{}) {

}

// +a (Admin) Grants admin status to a user, typically a level between operator (+o) and owner (+q). Admins have significant control but might not have all the powers of the channel owner.
func a(params *map[string]interface{}) {

}

// +q (Owner/Quiet) On some IRC networks, +q indicates channel ownership, giving the user ultimate control over the channel. On others, +q is used to prevent a user from sending messages to the channel (they can see the messages but cannot contribute).
func q(params *map[string]interface{}) {

}

//
