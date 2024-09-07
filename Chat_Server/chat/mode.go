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

	param_channel := split_p[0]
	log.Debug().Msgf("param_channel: %s", param_channel)
	channel, valid_channel := (*channel_map)[param_channel]
	if !valid_channel {
		log.Debug().Msgf("So such chan!")
		return ERR_NOSUCHCHANNEL("", cmd_params)
	}

	// pair off channel
	split_p = split_p[1:]

	log.Debug().Msgf("split_p before len(split_p) == 1: %s", split_p)

	if len(split_p) == 0 {
		// e.g. regular users can run MODE #chan w/o being an operator
		// only when they pass modes are they denied
		return RPL_CHANNELMODEIS("", channel.Name, channel.FmtModes())
	}

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
		log.Error().Msg("Client isnt channel operator")
		return ERR_NOTONCHANNEL("", channel.Name)
	}

	_server_metadata := params["server_metadata"]
	server_metadata := _server_metadata.(map[string]string)
	_task_runner := params["task_runner"]
	task_runner := _task_runner.(chan []*Task)

	supported_channel_modes := server_metadata["channelmodes"]
	mode_task := []*Task{}

	current_offset := 1
	seen := map[int]bool{}

	for mode_param_offset, m := range split_p {
		if seen[mode_param_offset] {
			// TODO
			// out of bounds issue here possibly
			current_offset = mode_param_offset + 1
			continue
		}

		action := string(m[0])
		log.Debug().Msgf("_mode: %s, action: %s", m, action)
		if len(action) == 0 || (action != "+" && action != "-") || len(m) == 1 {
			log.Debug().Msgf("action empty, action != + or -, m len == 1")
			// if no action is specified, we don't know whether to add or remove the modes
			return ERR_NEEDMOREPARAMS("")
		}

		for _, mode := range m {
			str_mode := string(mode)
			if str_mode == "+" || str_mode == "-" {
				action = str_mode
				continue
			}

			if !strings.Contains(supported_channel_modes, str_mode) {
				unknown_mode_task := NewTask(UNICAST, ERR_UNKNOWNMODE("", str_mode).Msg()+" \r\n", 0.0, client.ClientConn, nil)
				mode_task = append(mode_task, unknown_mode_task)
				continue
			}

			fn_params := map[string]interface{}{
				"params": "",
			}

			current_mode_params := ""

			// if command requires params
			// these are the only modes that do
			if str_mode == "b" || str_mode == "k" || str_mode == "v" || str_mode == "l" || str_mode == "e" || str_mode == "I" || str_mode == "o" || str_mode == "q" {
				// check the offset index to see if there is a corresponding param for it.
				// For example:
				//					                       co
				//				                      mpo
				//  			   0     1	  2	  3    4    5
				// MODE #channel +nbov Alice Bob John +no Pickle
				// +n is applied since it requires no params
				// +o then checks if there is a corresponding param for it
				// e.g. we are on
				log.Debug().Msg("Mode func expects params!")
				n := mode_param_offset + current_offset

				if action == "-" && (str_mode == "l" || str_mode == "k") {
					log.Debug().Msgf("Removing mode %s, no need to parse params", str_mode)
				} else {
					if len(split_p) < current_offset {
						log.Debug().Msgf("NEED MORE PARAMS")
						return ERR_NEEDMOREPARAMS("")
					}

					current_mode_params = split_p[n]
					log.Debug().Msgf("Params for func %s: %s", str_mode, current_mode_params)
					if len(current_mode_params) == 0 {
						log.Error().Msgf("Somehow the modes params are empty")
						return ERR_NEEDMOREPARAMS("")
					}
				}
				seen[n] = true
				current_offset++
			}

			fn := mode_fns[str_mode]
			if fn == nil {
				log.Error().Msgf("Unknown mode %s - no mapping in mode map", str_mode)
				return ERR_UNKNOWNMODE("", str_mode)
			}
			fn_params["params"] = current_mode_params
			response := fn(&fn_params)

			if response != nil {
				return *response
			}

			// TODO
			// need to run the command each time because:
			// -o Alice will remove Alice as an operator
			// -k will remove password entierly?
			// yes not all modes that require params will require them when using `-`

			if action == "+" {
				curr_mode := NewMode(str_mode, current_mode_params)
				channel.AddMode(*curr_mode)
			} else if action == "-" {
				channel.RemoveMode(str_mode, current_mode_params)
			}

			// TODO
			// I think these mode updates to channels are multicasts
			// a lot of them are, but are ALL of them?

			// TODO
			// we can't just send these responses without running the actual
			// corresponding MODE command function below.
			// this is because the actual command could fail for whatever reason

			// TODO
			// how to get args for each corresponding command?

			// e.g.
			// :Bob!bob@host MODE #example +i
			client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())
			msg := fmt.Sprintf(":%s MODE %s %s%s \r\n", client_details, channel.Name, action, str_mode)
			valid_mode_task := NewTask(MULTICAST, msg, 0.0, client.ClientConn, channel)
			mode_task = append(mode_task, valid_mode_task)
		}
	}

	task_runner <- mode_task

	return EMPTY_RESPONSE()
}

////// MODES

// this then makes two sources of truth for modes - the ServerManager or IrcServer and this one
// annoying to update both (also not like modes change frequently, though)
var mode_fns = map[string]func(params *map[string]interface{}) *Response{
	"p": p,
	"o": o,
	"s": s,
	"i": i,
	"m": m,
	"n": n,
	"t": t,
	"l": l,
	"k": k,
	"b": b,
	"e": e,
	"P": P,
	"v": v,
	"I": I,
	"r": r,
	"R": R,
	"z": z,
	"M": M,
	"c": c,
	"C": C,
	"a": a,
	"q": q,
}

// +o (operator): add user as an operator to the channel
func o(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode o, params: %s", l_params)
	return nil
}

// +p (Private): The channel is not visible in the channel list.
func p(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode p")
	return nil
}

// +s (Secret): The channel is hidden from public view and channel lists.
func s(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode s")
	return nil
}

// +i (Invite-Only): Users must be invited to join the channel.
func i(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode i")
	return nil
}

// Moderation:
// +m (Moderated): Only users with voice (+v) or operator status can speak.
func m(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode m")
	return nil
}

// +n (No External Messages): Prevents users outside the channel from sending messages to it.
func n(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode n")
	return nil
}

// +t (Topic Protection): Only operators can change the topic.
func t(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode t")
	return nil
}

// User Limits and Restrictions:
// +l (Limit): Sets a maximum number of users allowed in the channel.
// Example: /mode #channel +l 50 (limits the channel to 50 users)
func l(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode l, params: %s", l_params)
	return nil
}

// +k (Keyed): Requires a password (key) to join the channel.
// Example: /mode #channel +k secretpassword
func k(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode k, params: %s", l_params)
	return nil
}

// Bans and Exemptions:
// +b (Ban List): Bans specific users or masks from the channel.
// Example: /mode #channel +b *!*@example.com (bans all users from example.com)
func b(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode b, params: %s", l_params)
	return nil
}

// +e (Ban Exemption): Exempts specific users from being affected by a ban.
// Example: /mode #channel +e *!*@trusted.com
func e(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode e, params: %s", l_params)
	return nil
}

// +P persistent: to persist the channel after all users leave
func P(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode P")
	return nil
}

// +v Voice
// Example: /mode #channel +v username
func v(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode v, params: %s", l_params)
	return nil
}

// +I invite exception Allows specific users (using a mask) to join an invite-only (+i) channel without needing an invitation.
// /mode #channel +I *!*@friend.com
func I(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode I, params: %s", l_params)
	return nil
}

// +r Registered Channel (with ChanServ) Indicates that the channel is registered, typically with services like ChanServ. Some servers use this mode to show that a channel has been formally registered.
func r(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode r")
	return nil
}

// +z Secure only (SSL/TLS) connections only
func z(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode z")
	return nil
}

// +R Registerd users only (through NickServ)
func R(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode R")
	return nil
}

// +M (Moderated for Unregistered Users) Only registered users (e.g., users identified by NickServ) can send messages to the channel. Unregistered users can join but cannot talk.
func M(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode M")
	return nil
}

// +C (No CTCP) Blocks CTCP (Client-To-Client Protocol) messages, which are often used for things like requesting information from another client or performing actions like /me.
func C(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode C")
	return nil
}

// +c (No Color) Prevents users from using colored text in the channel. This is often used to reduce spam or unwanted formatting.
func c(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode c")
	return nil
}

// +a (Admin) Grants admin status to a user, typically a level between operator (+o) and owner (+q). Admins have significant control but might not have all the powers of the channel owner.
func a(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode a")
	return nil
}

// +q (Owner/Quiet) On some IRC networks, +q indicates channel ownership, giving the user ultimate control over the channel. On others, +q is used to prevent a user from sending messages to the channel (they can see the messages but cannot contribute).
func q(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode q, params: %s", l_params)
	return nil
}
