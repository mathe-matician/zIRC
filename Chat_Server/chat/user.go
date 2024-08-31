package chat

import (
	"fmt"
	c "zirc/client"

	"github.com/phuslu/log"
)

func user(params map[string]interface{}) Response {
	log.Info().Msg("Running USER...")

	_client, ok := params["client"]
	if !ok {
		log.Error().Msg("Client not passed to USER command!!")
		return ERR_UNKNOWNERROR("")
	}

	user, ok := params["params"]
	if !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	client := _client.(*c.Client)
	client.SetUser(user.(string))

	res := EMPTY_RESPONSE()

	if len(client.Nick()) != 0 && !client.Registered {
		client.Registered = true

		client_nick := client.Nick()
		_server_metadata := params["server_metadata"]
		server_metadata := _server_metadata.(map[string]string)

		server_name := server_metadata["name"]
		server_version := server_metadata["version"]
		server_usermodes := server_metadata["usermodes"]
		server_channelmodes := server_metadata["channelmodes"]
		server_creation_date := server_metadata["date"]

		client_details := fmt.Sprintf("%s@%s!%s", client_nick, client.User(), client.Ip())

		responses := WELCOME_WRAPPER(
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
