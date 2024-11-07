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

	split_p := strings.Split(cmd_params, " ")

	if len(split_p) == 0 {
		log.Debug().Msgf("no params passed to MODE")
		return ERR_NEEDMOREPARAMS("")
	}

	// TODO
	// check to see if channel
	// if not channel must be user, so check user
	target := split_p[0]
	log.Debug().Msgf("target: %s", target)
	if len(target) == 0 {
		log.Debug().Msgf("no target")
		return ERR_NEEDMOREPARAMS("")
	}

	first_char := string(target[0])
	is_chan := false
	var channel *Channel
	var valid_channel bool
	var target_client *Client
	var valid_target_client bool

	if first_char == GENERAL_CHAN_PREFIX || first_char == LOCAL_CHAN_PREFIX || first_char == MODELESS_CHAN_PREFIX {
		// then this is a channel
		is_chan = true
		if !IsValidChanName(target) {
			log.Debug().Msgf("Invalid chan name")
			return ERR_BADCHANMASK("")
		}
		// channel_map := (*g_Server)._MessageManager.ChannelMap
		channel, valid_channel = (*channel_map)[target]
		if !valid_channel {
			log.Debug().Msgf("So such chan!")
			return ERR_NOSUCHCHANNEL("", cmd_params)
		}
	} else {
		// else this is a user
		if !IsValidNickName(target) {
			log.Debug().Msgf("Invalid nick!")
			return ERR_NOSUCHNICK("")
		}

		client_map := (*g_Server)._MessageManager.ClientMap
		target_client, valid_target_client = (*client_map)[target]
		if !valid_target_client {
			log.Debug().Msgf("So such nick!")
			return ERR_NOSUCHNICK("")
		}
	}

	// pair off channel
	split_p = split_p[1:]

	log.Debug().Msgf("split_p before len(split_p) == 1: %s", split_p)

	if len(split_p) == 0 || (len(split_p) == 1 && len(split_p[0]) == 0) {
		// e.g. regular users can run MODE #chan w/o being an operator
		// only when they pass modes are they denied
		if is_chan {
			log.Debug().Msgf("Sending RPL_CHANNELMODEIS")
			return RPL_CHANNELMODEIS("", target, channel.FmtModes())
		} else {
			log.Debug().Msgf("Sending RPL_UMODEIS")
			return RPL_UMODEIS("", target, target_client.FmtModes())
		}
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
	if is_chan {
		if _, channel_operator := channel.Operators[client_nick]; !channel_operator {
			log.Error().Msg("Client isnt channel operator")
			return ERR_NOTONCHANNEL("", channel.Name)
		}
	}

	_server_metadata := params["server_metadata"]
	server_metadata := _server_metadata.(map[string]string)
	_task_runner := params["task_runner"]
	task_runner := _task_runner.(chan []*Task)

	supported_channel_modes := server_metadata["channelmodes"]
	mode_task := []*Task{}

	current_offset := 1
	seen := map[int]bool{}

	res_modes := ""
	res_final_params := ""
	res_action := ""
	prev_action := ""

	log.Debug().Msgf("Before MODE loop")
	for mode_param_offset, m := range split_p {
		log.Debug().Msgf("Starting inner MODE loop. mode_param_offset: %d, m: %s", mode_param_offset, m)
		if seen[mode_param_offset] {
			// TODO
			// out of bounds issue here possibly
			log.Debug().Msgf("Already seen index: %d, value: %s", mode_param_offset, m)
			current_offset = mode_param_offset
			continue
		}

		action := string(m[0])
		log.Debug().Msgf("_mode: %s, action: %s", m, action)
		pre_modes := string(m[1:])
		var valid_mode bool
		if is_chan {
			_, valid_mode = channel_modes[pre_modes]
		} else {
			_, valid_mode = user_modes[pre_modes]
		}
		if len(pre_modes) == 1 && !valid_mode {
			log.Debug().Msgf("Not a valid mode. pre_modes: %s", pre_modes)
			unknown_mode_task := NewTask(UNICAST, ERR_UNKNOWNMODE("", pre_modes).Msg()+"\r\n", 0.0, client.ClientConn, nil, false)
			mode_task = append(mode_task, unknown_mode_task)
			continue
		}

		if len(action) == 0 || (action != "+" && action != "-") || len(m) == 1 {
			log.Debug().Msgf("action empty, action(%s) != + or -, m[1:](%s) len == 1", action, m[1:])
			// if no action is specified, we don't know whether to add or remove the modes
			// need_more_params := NewTask(UNICAST, ERR_NEEDMOREPARAMS("").Msg()+" \r\n", 0.0, client.ClientConn, nil, false)
			// mode_task = append(mode_task, need_more_params)
			// continue
			return ERR_NEEDMOREPARAMS("")
		}

		for _, mode := range m {
			str_mode := string(mode)
			log.Debug().Msgf("Starting inner inner for loop. Mode: %s", str_mode)
			if str_mode == "+" || str_mode == "-" {
				log.Debug().Msgf("str_mode is action: %s. Continuing loop", str_mode)
				action = str_mode
				continue
			}

			if is_chan {
				if !strings.Contains(supported_channel_modes, str_mode) {
					log.Debug().Msgf("Not a supported channel mode: %s. Continuing", str_mode)
					unknown_mode_task := NewTask(UNICAST, ERR_UNKNOWNMODE("", str_mode).Msg()+"\r\n", 0.0, client.ClientConn, nil, false)
					mode_task = append(mode_task, unknown_mode_task)
					continue
				}
			} else {
				if !supported_user_mode(str_mode) {
					log.Debug().Msgf("Not a supported user mode: %s. Continuing", str_mode)
					unknown_mode_task := NewTask(UNICAST, ERR_UNKNOWNMODE("", str_mode).Msg()+"\r\n", 0.0, client.ClientConn, nil, false)
					mode_task = append(mode_task, unknown_mode_task)
					continue
				}
			}

			// TODO
			// successful responses should include ALL modes issued in the command
			// not just a single one at a time.
			// only errors are issued in additional responses
			// e.g.

			// :YourNick!username@host MODE #channel +oXt John will give response:
			//
			// :YourNick!username@host MODE #channel +ot John
			//:irc.example.com 472 YourNick X :is unknown mode char to me

			fn_params := map[string]interface{}{
				"params": "",
			}

			current_mode_params := ""

			// if command requires params (only applicable to channel modes)
			// these are the only modes that do
			// because there is overlap between user and channel modes but they are different
			// if it is a user mode, then we don't care because they don't need params
			// if it is a channel mode, then check if it requires params
			if is_chan && mode_requires_params(str_mode) {
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
				log.Debug().Msgf("Starting n at offset: %d", n)

				if action == "-" && (str_mode == "l" || str_mode == "k") {
					log.Debug().Msgf("Removing mode %s, no need to parse params", str_mode)
				} else {
					log.Debug().Msgf("Params for func %s: %s", str_mode, split_p)

					if len(split_p) < current_offset {
						log.Debug().Msgf("NEED MORE PARAMS")
						return ERR_NEEDMOREPARAMS("")
					}

					if n > len(split_p) {
						log.Debug().Msgf("The index %d is greater than the length of the slice %s!!", n, split_p)
						return ERR_NEEDMOREPARAMS("")
					}
					log.Debug().Msgf("Before split_p[n]. n: %d. split_p: %s, len(split_p): %d", n, split_p, len(split_p))
					if n >= len(split_p) {
						log.Debug().Msgf("Somehow the mode's params are empty")
						return ERR_NEEDMOREPARAMS("")
					}
					current_mode_params = split_p[n]
				}
				seen[n] = true
				current_offset++
			}

			var fn func(params *map[string]interface{}) *Response
			if is_chan {
				fn = channel_modes[str_mode]
			} else {
				fn = user_modes[str_mode]
			}

			if fn == nil {
				log.Debug().Msgf("Unknown mode %s - no mapping in mode map", str_mode)
				return ERR_UNKNOWNMODE("", str_mode)
			}
			fn_params["params"] = current_mode_params
			response := fn(&fn_params)

			if response != nil {
				return *response
			}

			if action == "+" {
				curr_mode := NewMode(str_mode, current_mode_params)
				if is_chan {
					channel.AddMode(*curr_mode)
				} else {
					target_client.AddMode(*curr_mode)
				}
			} else if action == "-" {
				if is_chan {
					channel.RemoveMode(str_mode, current_mode_params)
				} else {
					target_client.RemoveMode(str_mode)
				}
			}

			// TODO
			// I think these mode updates to channels are multicasts
			// a lot of them are, but are ALL of them?

			// TODO
			// we can't just send these responses without running the actual
			// corresponding MODE command function below.
			// this is because the actual command could fail for whatever reason

			// e.g.
			// :Bob!bob@host MODE #example +i
			res_params := ""
			if len(current_mode_params) != 0 {
				res_params += " " + current_mode_params
			}

			var final_valid_mode bool
			if is_chan {
				_, final_valid_mode = channel_modes[str_mode]
			} else {
				_, final_valid_mode = user_modes[str_mode]
			}

			if (action == "+" || action == "-") && final_valid_mode {
				if prev_action == action {
					res_action = ""
				} else {
					res_action = action
				}

				prev_action = action

				res_modes += res_action + str_mode
				res_final_params += res_params
			}
		}
	}
	client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())
	// TODO
	// try sending target
	msg := fmt.Sprintf(":%s MODE %s %s%s\r\n", client_details, target, res_modes, res_final_params)
	var valid_mode_task *Task
	if is_chan {
		// MULTICAST since other channel members should see these server responses
		// UNLESS they have some user mode set to NOT see them.
		valid_mode_task = NewTask(MULTICAST, msg, 0.0, client.ClientConn, channel, false)
	} else {
		valid_mode_task = NewTask(UNICAST, msg, 0.0, client.ClientConn, nil, false)
	}
	mode_task = append(mode_task, valid_mode_task)

	current_offset = 1
	log.Debug().Msgf("Before sending mode_task slice off: %v", mode_task)
	task_runner <- mode_task

	return EMPTY_RESPONSE()
}

////// MODES

// this then makes two sources of truth for modes - the MessageManager or IrcServer and this one
// annoying to update both (also not like modes change frequently, though)
var channel_modes = map[string]func(params *map[string]interface{}) *Response{
	"p": p,
	"o": o_chan,
	"s": s_chan,
	"i": i_chan,
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
	"r": r_chan,
	"R": R,
	"z": z,
	"M": M,
	"c": c,
	"C": C,
	"a": a_chan,
	"q": q,
}

var user_modes = map[string]func(params *map[string]interface{}) *Response{
	"i": i_user,
	"o": o_user,
	"w": w,
	"x": x,
	"s": s_user,
	"a": a_user,
	"r": r_user,
	"D": D,
	"C": C,
	"g": g,
	"G": g,
	"O": O,
	"S": S,
	"B": B,
	"W": W,
	"H": H,
}

func supported_user_mode(mode string) bool {
	for k, _ := range user_modes {
		if k == mode {
			return true
		}
	}
	return false
}

// use `return nil` as a good thing below
// if you need to exit early, return an error

// +o (operator): add user as an operator to the channel
func o_chan(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode o chan, params: %s", l_params)

	if len(l_params) == 0 {
		res := ERR_NEEDMOREPARAMS("")
		return &res
	}

	// todo
	// add to Channel.Operators client map

	return nil
}

// or if done on a user
// (operator): This grants the user IRC operator status, giving them elevated privileges such as kicking or banning users, shutting down servers, etc.
// /MODE Bob +o   # Bob becomes an IRC operator
func o_user(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode o user, params: %s", l_params)

	if len(l_params) == 0 {
		res := ERR_NEEDMOREPARAMS("")
		return &res
	}

	return nil
}

// +p (Private): The channel is not visible in the channel list.
func p(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode p")
	return nil
}

// +s
// Channel: (Secret): The channel is hidden from public view and channel lists.
// User: (receive server notices): This mode allows the user to receive special messages from the server.
func s_chan(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode s")
	return nil
}

// Allows the user to receive server notices (like warnings or alerts from IRC servers).
// /MODE David +s   # David will receive server notices
func s_user(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode s")
	return nil
}

// +i
// Channel: (Invite-Only):Users must be invited to join the channel.
// Channel version DOES NOT require params
func i_chan(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode i chan")
	return nil
}

// User: (invisible): This hides the user from other users who are not in the same channel as them. The user’s presence is not listed in /WHO or /NAMES commands unless the querying user shares a channel with them.
// User version DOES require params
func i_user(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode i user")
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
// OR if done on user
// (registered): This flag is often used to indicate that a user is registered with the network (e.g., via NickServ).
func r_chan(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode r")
	return nil
}

// Indicates the user is registered with services like NickServ (sometimes required to join specific channels).
// /MODE Frank +r   # Frank is marked as registered
func r_user(params *map[string]interface{}) *Response {
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
//
//	Prevents the user from receiving CTCP (Client-To-Client Protocol) requests like /PING, /VERSION, etc.
//
// /MODE Isla +C   # Isla won't receive CTCP requests
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
func a_chan(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode a")
	return nil
}

// +a
// Marks the user as away, meaning they won't respond to direct messages or queries (depends on the IRC client and server).
// /AWAY I'm away now!
// /MODE Eve +a    # Eve is marked as away
func a_user(params *map[string]interface{}) *Response {
	log.Debug().Msg("Running mode a")
	return nil
}

// +q (Owner/Quiet) On some IRC networks, +q indicates channel ownership, giving the user ultimate control over the channel. On others, +q is used to prevent a user from sending messages to the channel (they can see the messages but cannot contribute).
func q(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode q, params: %s", l_params)
	return nil
}

// User mode
// (wallops): Enables the user to receive special broadcast messages called "wallops" that are typically sent by IRC operators or administrators.
func w(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode w, params: %s", l_params)
	return nil
}

// User mode
// (cloaked): Some networks allow users to cloak their real IP address to protect their privacy, making their hostname hidden.
// /MODE Grace +x   # Grace's IP and hostname are hidden
func x(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode x, params: %s", l_params)
	return nil
}

// Makes the user deaf to all channel messages except for private messages (usually network-specific).
// /MODE Henry +D   # Henry won't see channel messages
func D(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode D, params: %s", l_params)
	return nil
}

// The user can receive notices from server administrators or messages broadcast globally or locally.
// /MODE John +g   # John receives global notices
func g(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode g, params: %s", l_params)
	return nil
}

// Opts the user out of receiving messages sent to IRC operators.
// /MODE Karen -O  # Karen will not receive operator messages
func O(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode O, params: %s", l_params)
	return nil
}

// Indicates that the user is protected by a network service, like ChanServ or NickServ. This mode is often set by the services automatically.
// /MODE Luke +S   # Luke is services protected
func S(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode S, params: %s", l_params)
	return nil
}

// Marks the user as a bot. Some networks use this to identify automated IRC clients.
// /MODE BotUser +B  # Marks BotUser as a bot
func B(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode B, params: %s", l_params)
	return nil
}

// Allows the user to send messages to other IRC operators.
// /MODE Mark +W  # Mark can send messages to IRC operators
func W(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode W, params: %s", l_params)
	return nil
}

// Prevents users from seeing whether you are marked as away in /WHO results.
// /MODE Nora +H  # Nora's away status is hidden
func H(params *map[string]interface{}) *Response {
	l_params := (*params)["params"].(string)
	log.Debug().Msgf("Running mode W, params: %s", l_params)
	return nil
}

// only channel modes take params
func mode_requires_params(mode string) bool {
	if mode == "b" || mode == "k" || mode == "v" || mode == "l" || mode == "e" || mode == "I" || mode == "o" || mode == "q" {
		return true
	}

	return false
}
