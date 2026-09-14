package chat

import (
	"fmt"
	"strings"

	"github.com/phuslu/log"
)

// USER <username> <mode> <unused> <realname>
// <username>: Required
// <mode>: Required (though usually set to 0)
// <unused>: Required (often set to *)
// <realname>: Required
func user(params map[string]interface{}) Response {
	log.Info().Msg("Running USER...")

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to USER command!!")
		return ERR_UNKNOWNERROR("")
	}

	_params, ok := params["params"]
	if _params == nil || !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	p := _params.(string)
	// TODO - USER command takes more params sent automatically by the client
	//		  not just "USER myuser"
	// TODO - can't split as realname param can contain spaces
	// if it contains spaces, it must be prefixed with `:`
	// p_split := re.FindAllStringSubmatch(p, -1)

	if len(p) == 0 {
		// if len(p) < 0 || p_split == nil {
		log.Error().Msg("not enough params")
		return ERR_NEEDMOREPARAMS("")
	}

	// user := p

	var user, mode, unused, realname string
	log.Debug().Msgf("USER: p: %s", p)
	trimmedParams := strings.Trim(p, " ")
	// TODO
	// these length checks are arbitrary
	// check if there is actually a limit for these
	if strings.Contains(trimmedParams, " ") {
		p_split := cmd_re.FindStringSubmatch(trimmedParams)

		log.Debug().Msgf("USER: p_split: %s, p_split len: %d", p_split, len(p_split))
		user = p_split[1]
		if len(user) > 100 {
			user = user[:100]
		}

		mode_split := cmd_re.FindStringSubmatch(strings.Trim(p_split[2], " "))
		log.Debug().Msgf("USER: mode_split: %s", mode_split)
		_mode := mode_split[1]
		if len(_mode) > 100 {
			_mode = _mode[:100]
		}
		mode = _mode

		unused_split := cmd_re.FindStringSubmatch(strings.Trim(mode_split[2], " "))
		log.Debug().Msgf("USER: unused_split: %s", unused_split)
		unused = "*"

		realname_trimmed := strings.Trim(unused_split[2], " ")
		realname = realname_trimmed
		if len(realname_trimmed) > 100 {
			realname = realname_trimmed[:100]
		}

		log.Debug().Msgf("USER: realname: %s", realname)
	} else {
		user = trimmedParams
	}

	client := _client.(*Client)
	client.SetUser(user)
	if len(realname) == 0 {
		realname = user
	}
	client.RealName = realname

	log.Info().EmbedObject(client).Msgf("user: %s, mode: %s, unused: %s, realname: %s", user, mode, unused, realname)

	res := EMPTY_RESPONSE()

	if !client.Registered && client.GetState("NICK") != "" {
		log.Info().Msg("Registering the USER...")

		client.Registered = true

		nickState := client.GetState("NICK")
		if nickState == "" {
			log.Error().EmbedObject(client).Msgf("NICK state empty?!")
			return ERR_UNKNOWNERROR("")
		}

		client.SetNick(nickState)
		client.RemoveState("NICK")

		client_nick := client.Nick()

		server_name := g_Server._MessageManager.Name
		server_version := g_Server.Version
		server_usermodes := G_Config.Server.User_modes
		server_channelmodes := G_Config.Server.Channel_modes
		server_creation_date := g_Server.CreationDate.String()

		client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())

		log.Info().Msg("Before accessing ClientMap...")
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

		if G_Config.Server.Ping_pong_enabled {
			go ping_send(client)
		}
	}

	return res
}
