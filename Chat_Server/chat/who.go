package chat

import (
	"regexp"

	"github.com/phuslu/log"
)

// WHO [<mask>] [o]
// <mask> (optional): A mask to filter the users for whom you want to get information.
//		This can be a channel name (e.g., #example) or a wildcard pattern (e.g., *, bob*).
// o (optional): If included, it restricts the search to only IRC operators.
// example response:
// :server 352 <requesting_user> <channel> <user> <host> <server> <nick> <H|G>[*][@|+] :<hopcount> <real name>

func who(params map[string]interface{}) Response {
	log.Debug().Msgf("Start of WHO command")

	_params, ok := params["params"]
	if _params == nil || !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	p := _params.(string)
	log.Debug().Msgf("WHO Params: %s", p)

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to USER command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	if len(p) == 0 {
		// show all users visible on server!!
		message_manager := (*g_Server)._MessageManager
		client_list := *(*message_manager).ClientList
		// create the size optimistically as the current server's client list
		list_all_users_tasks := make([]*Task, len(client_list))
		for _, c := range client_list {
			if c == nil {
				continue
			}

			// if the user is registerd and isn't invisible
			if c.Registered && !c.HasUserMode("i") {
				// TODO
				// need to get the current channel a user is in and share it in this response
				user_channel := ""
				chan_operator := ""
				if len(c.Channels) != 0 {
					// race condition here if a client leaves a channel
					// doesn;t matter much
					// and some other client runs WHO and this indexes any one client's Channel list
					// TODO
					// shouldn't this list all channels for all users?
					user_channel = c.Channels[0].Name
					// if chan operator to this chan
					chan_operator = "@"
					// if voice operator
					chan_operator = "+"
				} else {
					user_channel = "*"
				}

				away_status := "H"
				if len(c.AwayMessage) != 0 {
					away_status = "G"
				}

				operator_status := ""
				if c.HasUserMode("o") {
					operator_status = "*"
				}

				// TODO
				// get hopcount between this user running the WHO command
				// and the current user
				// will need to have shortest path algo in place to get this
				hopcount := ""
				if client.Host == c.Host {
					// on the same server
					hopcount = "0"
				} else {
					hopcount = "?"
				}

				res := RPL_WHOREPLY("", c.Nick(), user_channel, c.User(), c.Ip(), c.Host, c.Nick(), away_status, operator_status, chan_operator, hopcount, c.RealName)
				task := NewTask(UNICAST, res.Msg(), 0.0, client.ClientConn, nil, false)
				list_all_users_tasks = append(list_all_users_tasks, task)
			}
		}

		message_manager.Task_runner <- list_all_users_tasks
		return EMPTY_RESPONSE()
	}

	// TODO
	// Network Behavior: Some IRC networks might have configurations or commands (e.g., /WHO **0**)
	// that can show users across the entire network, but these are not part of the default behavior of the WHO command.

	// WHO <nick>: Lists all channels that the specified user is in, provided that you can see those channels (subject to visibility settings and modes like +i).
	// WHO <channel>: Only lists the users in the specified channel.
	// WHO (no arguments): Lists users from all visible channels they are in (public channels or channels you're a member of).

	// can handle * and ?
	// wildcard_match := `.*`
	// singlechar_match := `.?`
	// e.g. split word and then add the match var
	// who_match := regexp.MustCompile(``)

	// can also add a host name!!!!
	// e.g. user@host.com, WHO ?@host.com (all users on host.com server)

	if isChannel(p) {
		chan_map := g_Server.GetChannelMap()
		// TODO
		// parse hidden channels?

		channel, ok := (*chan_map)[p]
		if channel == nil || !ok {
			return ERR_UNKNOWNERROR("")
		}

		for _, c := range (*channel).UserList {
			if c == nil || c.HasUserMode("i") {
				continue
			}

			away_status := "H"
			if len(c.AwayMessage) != 0 {
				away_status = "G"
			}

			// TODO
			// get hopcount between this user running the WHO command
			// and the current user
			// will need to have shortest path algo in place to get this
			hopcount := ""
			if client.Host == c.Host {
				// on the same server
				hopcount = "0"
			} else {
				hopcount = "?"
			}

			chan_operator := ""
			operator_status := ""
			_, ok := channel.Operators[c.Nick()]
			if ok {
				operator_status = "*"
				chan_operator = " @ "
			}

			cwho := RPL_WHOREPLY("", c.Nick(), channel.Name, c.User(), c.Ip(), c.Host, c.Nick(), away_status, operator_status, chan_operator, hopcount, c.RealName)
			task := NewTask(UNICAST, cwho.Msg(), 0.0, client.ClientConn, nil, false)
			g_Server._MessageManager.Task_runner <- []*Task{task}
		}
	} else {
		// IS CLIENT
		// NOTE: wildcards only matter for nick

		// search users
		// TODO
		// need cross server routing here

		// client_list := make([]*Client, 0)
		var who_re *regexp.Regexp

		regexString := ""
		containsRegexMatch := false
		for _, c := range p {
			if c == '*' || c == '?' {
				containsRegexMatch = true
				regexString += "."
			}
			regexString += string(c)
		}

		client_map := g_Server.GetClientMap()
		client_who_tasks := []*Task{}

		if !containsRegexMatch {
			searchedClient, ok := (*client_map)[p]
			if searchedClient == nil || !ok {
				// TODO
				// is there a better error to return here?
				return ERR_UNKNOWNERROR("")
			}
			if len(searchedClient.Channels) == 0 {
				hopcount := "0"
				away_status := "H"
				if len(searchedClient.AwayMessage) != 0 {
					away_status = "G"
				}
				cwho := RPL_WHOREPLY("", client.Nick(), "*", searchedClient.User(), searchedClient.Ip(), searchedClient.Host, searchedClient.Nick(), away_status, "", "", hopcount, searchedClient.RealName)
				task := NewTask(UNICAST, cwho.Msg(), 0.0, client.ClientConn, nil, false)
				client_who_tasks = append(client_who_tasks, task)
			} else {
				for _, ch := range searchedClient.Channels {
					if ch == nil {
						log.Debug().Msgf("WHO command, channel was nil, continuing")
						continue
					}
					away_status := "H"
					if len(searchedClient.AwayMessage) != 0 {
						away_status = "G"
					}
					chan_operator := ""
					operator_status := ""
					// TODO
					// need to route to find this channel
					_, ok := ch.Operators[searchedClient.Nick()]
					if ok {
						operator_status = "*"
						chan_operator = " @ "
					}
					hopcount := "0"

					cwho := RPL_WHOREPLY("", client.Nick(), ch.Name, searchedClient.User(), searchedClient.Ip(), searchedClient.Host, searchedClient.Nick(), away_status, operator_status, chan_operator, hopcount, searchedClient.RealName)
					task := NewTask(UNICAST, cwho.Msg(), 0.0, client.ClientConn, nil, false)
					client_who_tasks = append(client_who_tasks, task)
				}
			}
		} else {
			who_re = regexp.MustCompile(regexString)

			for key, cli := range *client_map {
				match := who_re.FindAllStringSubmatch(key, -1)
				log.Debug().Msgf("who match: %v", match)
				if len(match) != 0 {
					log.Debug().Msgf("Match successful on: %v", match)

					if len(cli.Channels) == 0 {
						// TODO
						// get real hopcount
						log.Debug().Msgf("Client hasn't JOINed any channels")
						hopcount := "0"
						away_status := "H"
						if len(cli.AwayMessage) != 0 {
							away_status = "G"
						}
						cwho := RPL_WHOREPLY("", client.Nick(), "*", cli.User(), cli.Ip(), cli.Host, cli.Nick(), away_status, "", "", hopcount, cli.RealName)
						task := NewTask(UNICAST, cwho.Msg(), 0.0, client.ClientConn, nil, false)
						client_who_tasks = append(client_who_tasks, task)
					} else {
						log.Debug().Msgf("Client has JOINed channels, looping through them...")
						for _, ch := range cli.Channels {
							if ch == nil {
								log.Debug().Msgf("WHO command, channel was nil, continuing")
								continue
							}
							away_status := "H"
							if len(cli.AwayMessage) != 0 {
								away_status = "G"
							}
							chan_operator := ""
							operator_status := ""
							// TODO
							// need to route to find this channel
							_, ok := ch.Operators[cli.Nick()]
							if ok {
								operator_status = "*"
								chan_operator = " @ "
							}
							hopcount := "0"

							cwho := RPL_WHOREPLY("", client.Nick(), ch.Name, cli.User(), cli.Ip(), cli.Host, cli.Nick(), away_status, operator_status, chan_operator, hopcount, cli.RealName)
							task := NewTask(UNICAST, cwho.Msg(), 0.0, client.ClientConn, nil, false)
							client_who_tasks = append(client_who_tasks, task)
						}
					}
				}
			}
		}

		// log.Debug().Msgf("WHO client_list: %s", client_list)
		g_Server._MessageManager.Task_runner <- client_who_tasks
	}

	// As iter through client list, if client has +i user mode,
	// they won't be returned in this command response

	return EMPTY_RESPONSE()
}
