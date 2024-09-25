package chat

import "github.com/phuslu/log"

func notice(params map[string]interface{}) Response {
	log.Debug().Msgf("Start of NOTICE command")

	_params, ok := params["params"]
	if _params == nil || !ok {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	p := _params.(string)
	log.Debug().Msgf("WHO Params: %s", p)

	return nil
}
