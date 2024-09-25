package chat

import (
	"fmt"

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
	user := p
	// user := p_split[0]
	// mode := p_split[1]
	// unused := p_split[2]
	// realname := p_split[3]

	client := _client.(*Client)
	client.SetUser(user)

	res := EMPTY_RESPONSE()

	if len(client.Nick()) != 0 && !client.Registered {
		log.Info().Msg("Registering the USER...")

		client.Registered = true

		client_nick := client.Nick()
		// _server_metadata := params["server_metadata"]
		// server_metadata := _server_metadata.(map[string]string)

		server_name := g_Server._MessageManager.Name
		server_version := g_Server.Version
		server_usermodes := g_Server.Config["IRC_USER_MODES"]
		server_channelmodes := g_Server.Config["IRC_CHANNEL_MODES"]
		server_creation_date := g_Server.CreationDate.String()

		client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())

		log.Info().Msg("Before accessing ClientMap...")
		svr_mang := g_Server._MessageManager
		clint_map := svr_mang.ClientMap
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

		_task_runner := params["task_runner"]
		task_runner := _task_runner.(chan []*Task)
		task_runner <- responses
	}

	return res
}
