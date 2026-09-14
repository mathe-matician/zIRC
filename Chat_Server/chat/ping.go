package chat

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

// ping is when we receive a PING from s2s communication
func ping(params map[string]interface{}) Response {
	log.Debug().Msg("Running PING...")

	_client, ok := params["client"]
	if _client == nil || !ok {
		log.Error().Msg("Client not passed to PASS command!!")
		return ERR_UNKNOWNERROR("")
	}
	client := _client.(*Client)

	pingParams, ok := params["params"]
	if pingParams == nil || !ok || len(pingParams.(string)) == 0 {
		log.Error().Msg("Params not in map!")
		return ERR_NEEDMOREPARAMS("")
	}

	pingToken := pingParams.(string)

	log.Debug().Msgf("PING token: %v", pingToken)

	sendingNode := g_Server.Servers.GetServerByConn(client.ClientConn)
	if sendingNode.Conn.State == BURST_SEND {
		log.Debug().Msgf("PING: Conn.State == BURST_SEND")
		// this means that this connection is currently in BURST_SEND mode, so it is signaling
		// that this is the end of the burst.
		// and that we can now burst to them our network topology.
		sendingNode.Conn.SetStateBurstRecv()
		pongCmd := fmt.Sprintf("PONG :%s", pingToken)

		task_runner := g_Server._MessageManager.Task_runner

		_task := NewTask(UNICAST, pongCmd, 0.0, sendingNode.Conn.Conn, nil, false, "")
		pongTask := []*Task{}
		pongTask = append(pongTask, _task)

		task_runner <- pongTask
	} else if sendingNode.Conn.State != REGISTERED {
		// if we are receiving a ping from a non registered conn
		// and it isn't currently in burst_send
		// then I'm not sure what is going on - this isn't right
		log.Error().Msgf("sendingNode.Conn.State != REGISTERED and sendingNode.Conn.State != BURST_SEND")
		return ERR_UNKNOWNERROR("")
	} else {
		log.Debug().Msgf("PING: otherwise send this to the go routine that is handling pingpong stuff?")
		// otherwise send this to the go routine that is handling pingpong stuff?
		// g_Server.Servers.Tree[g_Server.Name].PingPongChan <- pingParams.(string)
	}

	return EMPTY_RESPONSE()
}

// ping_send_s2s is primarily used for initial s2s communication to single the end of BURST_SEND
func s2sPingToken(s *ServerNode) (string, *Task) {
	ping_token, err := uuid.NewV7()
	pingToken := ping_token.String()
	if err != nil {
		log.Error().EmbedObject(s).Msgf("Could not generate ping_token uuid!!!: %s", err.Error())
		pingToken = "TOKEN123"
	}

	cmdmsg := fmt.Sprintf("PING :%s", pingToken)

	// task_runner := g_Server._MessageManager.Task_runner

	task := NewTask(UNICAST, cmdmsg, 0.0, s.Conn.Conn, nil, false, "")

	return pingToken, task
	// ping_task := []*Task{}
	// ping_task = append(ping_task, _task)

	// task_runner <- ping_task

	// select {
	// case pong_token := <-s.PingPongChan:
	// 	_pongToken := pong_token[1:]
	// 	pongToken := strings.Trim(_pongToken, " ")
	// 	if pongToken != pingToken {
	// 		log.Warn().EmbedObject(s).Msgf("%s != %s", pongToken, pingToken)
	// 		g_Server.Servers.Remove(s)
	// 	}
	// case <-time.After(time.Duration(5) * time.Second):
	// 	return errors.New("S2S never received PONG, closing connection")
	// }
	// return nil
}

// ping_send is more primarily used for automatic ping keepalives
// which are sent upon client registration
// i.e. after NICK and USER are sent by client the server initiates the PING back to the client
func ping_send(c *Client) {
	// log.Debug().EmbedObject(c).Msgf("Running PING keepalive...")

	for {
		ping_token, err := uuid.NewV7()
		pingToken := ping_token.String()
		if err != nil {
			log.Error().EmbedObject(c).Msgf("Could not generate ping_token uuid!!!: %s", err.Error())
			pingToken = "TOKEN123"
		}

		conn := c.ClientConn
		cmdmsg := fmt.Sprintf("PING :%s", pingToken)

		task_runner := g_Server._MessageManager.Task_runner

		_task := NewTask(UNICAST, cmdmsg, 0.0, conn, nil, false, "")
		ping_task := []*Task{}
		ping_task = append(ping_task, _task)

		task_runner <- ping_task
		// TODO
		// add all parts of response (server, etc)
		// ping_msg := helpers.FormatResponse(cmdmsg)
		// if _, err := (*conn).Write([]byte(ping_msg)); err != nil {
		// 	log.Error().EmbedObject(c).Msgf("Error writing to client: %s", err.Error())
		// 	return
		// }

		timeout := G_Config.Server.Ping_pong_timeout

		_duration := G_Config.Server.Ping_pong_timeout_duration
		var duration time.Duration
		if _duration == "minute" {
			duration = time.Minute
		} else if _duration == "second" {
			duration = time.Second
		}

		select {
		case pong_token := <-c.PingPongChan:
			pongToken := pong_token[1:]
			if pongToken != pingToken {
				log.Warn().EmbedObject(c).Msgf("%s != %s", pongToken, pingToken)
				conn.Close()
			}
			// log.Debug().EmbedObject(c).Msgf("PONG successful")
		case <-time.After(time.Duration(timeout) * duration):
			log.Info().EmbedObject(c).Msgf("Server never received PONG, closing connection")
			conn.Close()
		}

		time.Sleep(duration)
	}
}

func ping_s2s(s *ServerNode, waitOnce bool) {
	log.Debug().Msgf("ping_s2s start")
	for {
		ping_token, err := uuid.NewV7()
		pingToken := ping_token.String()
		if err != nil {
			log.Error().EmbedObject(s).Msgf("Could not generate ping_token uuid!!!: %s", err.Error())
			pingToken = "TOKEN123"
		}

		conn := s.Conn.Conn
		cmdmsg := fmt.Sprintf("PING :%s", pingToken)

		task_runner := g_Server._MessageManager.Task_runner

		_task := NewTask(UNICAST, cmdmsg, 0.0, conn, nil, false, "")
		ping_task := []*Task{}
		ping_task = append(ping_task, _task)

		task_runner <- ping_task

		timeout := G_Config.Server.Ping_pong_timeout

		_duration := G_Config.Server.Ping_pong_timeout_duration
		var duration time.Duration
		if _duration == "minute" {
			duration = time.Minute
		} else if _duration == "second" {
			duration = time.Second
		}

		select {
		case pong_token := <-s.PingPongChan:
			log.Debug().Msgf("PING-PONG S2S pong recv")
			pongToken := pong_token[1:]
			if pongToken != pingToken {
				log.Warn().EmbedObject(s).Msgf("%s != %s", pongToken, pingToken)
				g_Server.Servers.Remove(s)
				return
			}
		case <-time.After(time.Duration(timeout) * duration):
			log.Warn().EmbedObject(s).Msgf("Server never received PONG, closing connection")
			g_Server.Servers.Remove(s)
			return
		}

		if waitOnce {
			log.Debug().Msgf("PING s2s loop: waitOnce == true. breaking")
			return
		}

		time.Sleep(duration)
	}
}
