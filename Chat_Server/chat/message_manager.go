package chat

import (
	"fmt"
	"strconv"
	"zirc/helpers"

	"github.com/phuslu/log"
)

type MessageManager struct {
	Name               string
	ClientList         *[]*Client
	ClientMap          *map[string]*Client
	ChannelMap         *map[string]*Channel
	ServerList         *[]*IrcServer
	WorkerPool         map[string]*Worker
	worker_tasks       chan []*Task
	Task_runner        chan []*Task
	results            chan string
	decreasing_workers bool
}

func NewMessageManager(client_list *[]*Client, server_list *[]*IrcServer) *MessageManager {
	init_worker_count, err := strconv.Atoi(helpers.GetEnv("IRC_SERVER_INIT_WORKER_COUNT", "3"))
	if err != nil {
		panic("can't init workers...")
	}

	if client_list == nil || server_list == nil {
		// the app's functionality requires the server manager being setup correctly
		// if its not, panic
		panic("message_manager: client_list or server_list is null! This cannot be!")
	}

	channel_list := make(map[string]*Channel)
	client_map := make(map[string]*Client)

	sm := &MessageManager{
		Name:         "",
		Task_runner:  make(chan []*Task),
		worker_tasks: make(chan []*Task),
		results:      make(chan string),
		ClientList:   client_list,
		ClientMap:    &client_map,
		ServerList:   server_list,
		ChannelMap:   &channel_list,
		WorkerPool:   make(map[string]*Worker), // TODO - do we even need to keep track of workers in the pool? !only if we want to scale them down by name - otherwise sending 'quit' to any arbitrary worker will kill it
	}

	for range init_worker_count {
		w := NewWorker()
		sm.WorkerPool[w.id.String()] = w
		go w.Work(sm.worker_tasks, sm.results, sm)
	}

	return sm
}

// Run starts the MessageManager which manages the Worker pool
func (sm *MessageManager) Run() {
	log.Info().Msg("MessageManager started")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	for {
		select {
		case client_task := <-sm.Task_runner:
			log.Info().Msgf("client_tasks: %v, len: %d", client_task, len(client_task))
			sm.worker_tasks <- client_task
		}
	}
}

// TODO - is this fn even needed?
func (sm *MessageManager) StartTask(task string, client_chan chan string) {
	if len(task) == 0 {
		return
	}

	if len(sm.WorkerPool) == 0 {
		panic("server manager has no workers!!")
	}

	lowest_worker := ""
	lowest_count := 0

	// TODO - make this smarter?
	for _, w := range sm.WorkerPool {
		if w == nil || w.frozen {
			log.Error().Msg("found worker is null or frozen - trying next")
			continue
		}

		if len(w.tasks) >= lowest_count {
			lowest_worker = w.id.String()
		}
	}

	log.Debug().Msgf("Lowest worker: %s, count: %d", lowest_worker, lowest_count)
}

// TODO
// this is NOT functional
// ScaleWorkerPool allows the message_manager to increase or delete workers from the pool
// if force is true, delete the workers immediately
// MUST be run in a go routine
func (sm *MessageManager) ScaleWorkerPool(by int, force bool) {
	if sm.decreasing_workers {
		log.Info().Msg("In the process of decreasing worker count")
		return
	}

	if by == 0 {
		log.Info().Msg("Must pass a number to increase or decrease by")
		return
	}

	if by > 0 {
		for range by {
			w := NewWorker()
			sm.WorkerPool[w.id.String()] = w
		}
	} else {
		worker_pool_len := len(sm.WorkerPool)
		if by > worker_pool_len {
			by = worker_pool_len
		}

		// run decrease workers
		sm.decreasing_workers = true
		decrease_by := by
		wrkrs := make([]string, 0)

		// TODO - make this smarter instead of taking the first x
		for _, v := range sm.WorkerPool {
			if decrease_by == 0 {
				break
			}
			v.frozen = true
			wrkrs = append(wrkrs, v.id.String())
		}

		for {
			// TODO - wait for these workers to terminate
			for _, w := range wrkrs {
				if force {
					// TODO - check if this value is in the map
					sm.WorkerPool[w] = nil
				}

				// else wait for work to decrease to 0 then "delete" the worker
			}

			break
		}

		sm.decreasing_workers = false
	}
}

func (sm *MessageManager) GetClientByNick(nick string) *Client {
	log.Debug().Msgf("GetClientByNick: %s", nick)
	client_map := sm.ClientMap
	client, ok := (*client_map)[nick]
	if !ok {
		return nil
	}

	return client
}

func (sm *MessageManager) Debug() {
	log.Info().Msgf("MessageManager Debug!")
}

func (sm *MessageManager) Route(msg string) {
	// route msg across spanning tree irc network

}
