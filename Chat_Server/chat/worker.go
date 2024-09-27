package chat

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type Worker struct {
	id           uuid.UUID
	tasks        []*Task
	current_load float64
	quit         chan int
	frozen       bool
	mu           sync.Mutex
}

func NewWorker() *Worker {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &Worker{
		id:           uid,
		tasks:        make([]*Task, 0),
		current_load: 0.0,
		quit:         make(chan int, 1),
	}
}

func (w *Worker) MarshalObject(e *log.Entry) {
	e.Str("id", w.id.String()).Float64("current_load", w.current_load)
}

var workerActionMap = map[string]map[string]bool{
	"create": map[string]bool{
		"chan": true,
		"user": true,
	},
	"delete": map[string]bool{
		"chan": true,
	},
	"update": map[string]bool{
		"chan": true,
	},
}

func (w *Worker) unicast(message_manager *MessageManager, task *Task) {
	// target := (*task.Target)
	// if target != nil {
	// 	// If we reach here, we are sending a message to both some other 1 client + the source client
	// 	// but something else needs to be sent to another client which is stored in Task.Target

	// 	// we only expect this to be for Client Targets...
	// 	// we should never get a Channel target here

	// 	// if task.FindTarget != "" {
	// 	// 	client_list := (*message_manager).ClientList
	// 	// 	for _, c := range *client_list {

	// 	// 	}
	// 	// }

	// 	client := (target).(*Client)
	// 	conn := *(client.ClientConn)
	// 	if conn == nil {
	// 		log.Error().EmbedObject(client).Msgf("Conn is nil!!")
	// 		return
	// 	}
	// 	if _, err := conn.Write([]byte(task.Task)); err != nil {
	// 		log.Error().EmbedObject(client).Msgf("Error writing to client: %s", err.Error())
	// 		return
	// 	}
	// }

	// send a response to the client performing the action
	c := (*task.ClientConn)
	if _, err := c.Write([]byte(task.Task)); err != nil {
		log.Error().Msgf("Error writing to client: %s", err.Error())
	}
}

func (w *Worker) broadcast(message_manager *MessageManager, task *Task) {

}

func (w *Worker) multicast(message_manager *MessageManager, task *Task) {
	if task == nil {
		log.Debug().EmbedObject(w).Msgf("Multicast task is nil!")
		return
	}
	the_target := (*task).Target
	if the_target == nil {
		log.Debug().EmbedObject(w).Msgf("Multicast channel is nil!")
		return
	}

	// TODO
	// cast Target to Channel or Client
	target_type := reflect.TypeOf(the_target).Name()
	var channel *Channel
	var client *Client
	if target_type == "Channel" {
		log.Debug().Msgf("Target type is Channel")
		channel = (*the_target).(*Channel)
		channel_user_map := (*channel).UserList
		for _, c := range channel_user_map {
			if c == nil {
				log.Debug().EmbedObject(w).Msgf("Client is nil - trying next")
				continue
			}

			// TODO some modes allow you to send messages to unregistered clients
			if c.Registered {
				log.Debug().EmbedObject(w).Msgf("Staring task: %s", task.Id.String())
				if c.ClientConn == nil {
					log.Error().EmbedObject(w).Msg("Client connection is nil!!")
					continue
				}
				conn := *(c.ClientConn)
				if _, err := conn.Write([]byte(task.Task)); err != nil {
					log.Error().EmbedObject(c).Msgf("Error writing to client: %s", err.Error())
					continue
				}
			} else {
				// TODO some modes allow you to send messages to unregistered clients
				log.Info().EmbedObject(w).Msg("Client not registered")
			}
		}
	} else if target_type == "Client" {
		// Is there a MULTICAST Client option here??
		// probably not?
		// maybe something a Server can only do...
		// need to figure out
		// i.e. why wouldn't we just send a UNICAST?
		log.Debug().Msgf("Target type is Client")
		client = (*the_target).(*Client)
		conn := *(client.ClientConn)
		if _, err := conn.Write([]byte(task.Task)); err != nil {
			log.Error().EmbedObject(client).Msgf("Error writing to client: %s", err.Error())
			return
		}
	} else {
		log.Error().Msgf("Unknown target type %s", target_type)
		return
	}
}

// Work
//
//	job: jobs received from the MessageManager
//	results: any results that are returned back to the MessageManager can be sent back to the client if needed
//
// TODO - this func may only need the MessageManager's ClientList and ServerList
func (w *Worker) Work(tasks chan []*Task, results chan string, message_manager *MessageManager) {
	for {
		select {
		case task := <-tasks:
			log.Info().EmbedObject(w).Msgf("Tasks received")
			message_manager.Debug()

			if len(task) == 0 {
				log.Warn().EmbedObject(w).Msgf("No tasks to run!")
				continue
			}

			if message_manager.ClientList == nil {
				msg := "message_manager client list is null! idk how we got to this point..."
				log.Error().EmbedObject(w).Msg(msg)
				panic(msg)
			}
			log.Debug().EmbedObject(w).Msgf("Client List Len: %d", len(*message_manager.ClientList))

			for _, task := range task {
				if task == nil {
					log.Warn().Msg("Task is null!")
					continue
				}

				if task.Type == UNICAST {
					log.Debug().Msgf("%s task", UNICAST)
					// TODO - need an efficient way to get the single client's connection info
					// unicast examples:
					// 	server to client (as in a response message)
					//	client to client (privmsg to single person)

					w.unicast(message_manager, task)
				} else if task.Type == MULTICAST {
					log.Debug().Msgf("%s task", MULTICAST)
					// multicast only purpose is to send client msgs to a specific channel
					w.multicast(message_manager, task)
				} else if task.Type == BROADCAST {
					log.Debug().Msgf("%s task", BROADCAST)
					// broadcast examples:
					//	server admin broadcast to all users (e.g. server going down for maintenance)
					w.broadcast(message_manager, task)
				} else if task.Type == SERVER {
					log.Debug().Msgf("%s task", SERVER)
					// split_task := strings.Split(task.Task, " ")

					// _, ok := workerActionMap[split_task[0]][split_task[1]]
					// if !ok {
					// 	log.Debug().Msgf("Not a valid worker action")
					// 	continue
					// }

				} else {
					log.Warn().Msgf("Unknown task type: %s", task.Type)

				}
			}

		case <-w.quit:
			log.Info().EmbedObject(w).Msg("quitting")
			return
		}
	}
}
