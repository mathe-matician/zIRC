package chat

import (
	"errors"
	"fmt"

	"github.com/phuslu/log"
)

type MessageManager struct {
	Name       string
	ClientList *[]*Client
	ClientMap  *map[string]*Client
	ChannelMap *map[string]*Channel
	// ServerList         *[]*IrcServer
	WorkerPool         map[string]*Worker
	worker_tasks       chan []*Task
	Task_runner        chan []*Task
	results            chan string
	decreasing_workers bool
	ready              bool
}

func NewMessageManager(client_list *[]*Client) *MessageManager {
	init_worker_count := G_Config.Server.Init_worker_count

	if client_list == nil {
		// the app's functionality requires the server manager being setup correctly
		// if its not, panic
		panic("message_manager: client_list is null! This cannot be!")
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
		// ServerList:   server_list,
		ChannelMap: &channel_list,
		WorkerPool: make(map[string]*Worker), // TODO - do we even need to keep track of workers?
	}

	// TODO
	// refactor this to not have a static go routine worker count, but just run tasks in go routines
	// go routines are cheap
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
	sm.ready = true
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

func (sm *MessageManager) GetClientByNick(nick string) *Client {
	log.Debug().Msgf("GetClientByNick: %s", nick)
	client_map := sm.ClientMap
	client, ok := (*client_map)[nick]
	if !ok {
		return nil
	}

	return client
}

func (mm *MessageManager) ClientMapInsert(client *Client) error {
	cm := mm.ClientMap
	client_uid := client.UID.String()
	if _, ok := (*cm)[client_uid]; ok {
		errmsg := fmt.Sprintf("client %s already exists in client map", client_uid)
		log.Info().Msg(errmsg)
		return errors.New(errmsg)
	}

	return nil
}

func (mm *MessageManager) ClientMapPatch(client *Client) error {
	cm := mm.ClientMap
	client_uid := client.UID.String()
	if _, ok := (*cm)[client_uid]; !ok {
		errmsg := fmt.Sprintf("client %s doesnt exists in client map", client_uid)
		log.Info().Msg(errmsg)
		return errors.New(errmsg)
	}

	return nil
}

func (mm *MessageManager) ClientMapDrop(client_uid string) error {
	cm := mm.ClientMap
	if _, ok := (*cm)[client_uid]; !ok {
		errmsg := fmt.Sprintf("client %s doesnt exists in client map", client_uid)
		return errors.New(errmsg)
	}

	return nil
}

func (sm *MessageManager) Debug() {
	log.Info().Msgf("MessageManager Debug!")
}

func (sm *MessageManager) Route(msg string) {
	// route msg across spanning tree irc network

}
