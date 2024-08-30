package server

import (
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"zirc/chat"
	c "zirc/client"
	"zirc/helpers"
	rc "zirc/remote_conn"
	t "zirc/task"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type Server interface {
	Run()
	handleConnection(conn *net.Conn)
	Stop()
}

type ServerConfig interface {
}

type IrcServer struct {
	DnsName       string
	Version       string
	CreationDate  time.Time
	Addr          string
	Role          string
	Listener      *net.Listener
	_ServerManger *ServerManager
	Servers       []*IrcServer
	Clients       []*c.Client
	Config        map[string]string
}

type Task struct {
	id     uuid.UUID
	Type   string
	weight float64
	task   string
}

type Worker struct {
	id           uuid.UUID
	tasks        []*t.Task
	current_load float64
	quit         chan int
	frozen       bool
}

type ServerManager struct {
	Name               string
	ClientList         *[]*c.Client
	ServerList         *[]*IrcServer
	WorkerPool         map[string]*Worker
	worker_tasks       chan []*t.Task
	Task_runner        chan []*t.Task
	results            chan string
	decreasing_workers bool
}

// Run starts the ServerManager which manages the Worker pool
func (sm *ServerManager) Run() {
	log.Info().Msg("ServerManager started")

	for {
		select {
		case client_task := <-sm.Task_runner:
			log.Info().Msgf("client_tasks: %s, len: %d", client_task, len(client_task))
			sm.worker_tasks <- client_task
		}
	}
}

// TODO - is this fn even needed?
func (sm *ServerManager) StartTask(task string, client_chan chan string) {
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

	log.Debug().Msgf("Lowest worker: %s, count: %s", lowest_worker, lowest_count)
}

// ScaleWorkerPool allows the servermanager to increase or delete workers from the pool
// if force is true, delete the workers immediately
// MUST be run in a go routine
func (sm *ServerManager) ScaleWorkerPool(by int, force bool) {
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

func (sm *ServerManager) Debug() {
	log.Info().Msgf("ServerManager Debug!")
}

func (w *Worker) MarshalObject(e *log.Entry) {
	e.Str("id", w.id.String()).Float64("current_load", w.current_load)
}

// task_weight_mapping := map[string]float64 {

// }

func (w *Worker) unicast(server_manager *ServerManager, task *t.Task) {

}

func (w *Worker) broadcast(server_manager *ServerManager, task *t.Task) {

}

func (w *Worker) multicast(server_manager *ServerManager, task *t.Task) {
	for _, c := range *server_manager.ClientList {
		if c == nil {
			log.Debug().EmbedObject(w).Msgf("Client is null - trying next")
			continue
		}

		if c.Registered {
			log.Debug().EmbedObject(w).Msgf("Staring task: %s", task.Id.String())
			if c.ClientConn == nil {
				log.Error().EmbedObject(w).Msg("Client connection is nil!!")
				break
			}
			conn := *(c.ClientConn)
			if _, err := conn.Write([]byte(task.Task)); err != nil {
				log.Error().EmbedObject(c).Msgf("Error writing to client: %s", err.Error())
				break
			}
			// iterating through tasks here doesn't make sense anymore
			// it is done outside of this loop
			// for _, task := range task {
			// 	if _, err := conn.Write([]byte(task.Task)); err != nil {
			// 		log.Error().EmbedObject(c).Msgf("Error writing to client: %s", err.Error())
			// 		break
			// 	}
			// }
		} else {
			log.Info().EmbedObject(w).Msg("Client not registered")
		}
	}
}

// Work
//
//	job: jobs received from the ServerManager
//	results: any results that are returned back to the ServerManager can be sent back to the client if needed
//
// TODO - this func may only need the ServerManager's ClientList and ServerList
func (w *Worker) Work(tasks chan []*t.Task, results chan string, server_manager *ServerManager) {
	for {
		select {
		case task := <-tasks:
			log.Info().EmbedObject(w).Msgf("Tasks received")
			server_manager.Debug()

			if len(task) == 0 {
				log.Warn().EmbedObject(w).Msgf("No tasks to run!")
				continue
			}

			if server_manager.ClientList == nil {
				msg := "servermanager client list is null! idk how we got to this point..."
				log.Error().EmbedObject(w).Msg(msg)
				panic(msg)
			}
			log.Debug().EmbedObject(w).Msgf("Client List Len: %d", len(*server_manager.ClientList))

			for _, task := range task {
				if task == nil {
					log.Warn().Msg("Task is null!")
					continue
				}

				if task.Type == t.UNICAST {
					log.Debug().Msgf("%s task", t.UNICAST)
					// TODO - need an efficient way to get the single client's connection info
					// unicast examples:
					// 	server to client (as in a response message)
					//	client to client (privmsg to single person)

					w.unicast(server_manager, task)
				} else if task.Type == t.MULTICAST {
					log.Debug().Msgf("%s task", t.MULTICAST)
					// multicast examples:
					//	client msg to channel

					w.multicast(server_manager, task)
				} else if task.Type == t.BROADCAST {
					log.Debug().Msgf("%s task", t.BROADCAST)
					// broadcast examples:
					//	server admin broadcast to all users (e.g. server going down for maintenance)
					w.broadcast(server_manager, task)
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

func NewServerManager(client_list *[]*c.Client, server_list *[]*IrcServer) *ServerManager {
	init_worker_count, err := strconv.Atoi(helpers.GetEnv("IRC_SERVER_INIT_WORKER_COUNT", "3"))
	if err != nil {
		panic("can't init workers...")
	}

	if client_list == nil || server_list == nil {
		// the app's functionality requires the server manager being setup correctly
		// if its not, panic
		panic("servermanager: client_list or server_list is null! This cannot be!")
	}

	sm := &ServerManager{
		Name:         "",
		Task_runner:  make(chan []*t.Task),
		worker_tasks: make(chan []*t.Task),
		results:      make(chan string),
		ClientList:   client_list, // only contains registered clients
		ServerList:   server_list,
		WorkerPool:   make(map[string]*Worker), // TODO - do we even need to keep track of workers in the pool?
	}

	for range init_worker_count {
		w := NewWorker()
		sm.WorkerPool[w.id.String()] = w
		go w.Work(sm.worker_tasks, sm.results, sm)
	}

	return sm
}

func NewWorker() *Worker {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &Worker{
		id:           uid,
		tasks:        make([]*t.Task, 0),
		current_load: 0.0,
		quit:         make(chan int, 1),
	}
}

func NewIrcServer(dns_name string, version string, addr string, server_role string, server_list *[]*IrcServer, client_list *[]*c.Client, config *map[string]string) *IrcServer {
	if len(dns_name) == 0 {
		dns_name = helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost")
	}

	if len(version) == 0 {
		version = helpers.GetEnv("IRC_SERVER_VERSION", "v99.99.99+default")
	}

	if len(server_role) == 0 {
		server_role = helpers.GetEnv("IRC_SERVER_ROLE", "leaf")
	}
	err := helpers.VerifyServerMode(server_role)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err.Error())
	}

	if server_list == nil {
		sl := make([]*IrcServer, 0)
		server_list = &sl
	}

	if client_list == nil {
		cl := make([]*c.Client, 0)
		client_list = &cl
	}

	if len(addr) == 0 {
		addr = helpers.GetEnv("IRC_HOST", "0.0.0.0") + ":" + helpers.GetEnv("IRC_PORT", "6667")
	}

	if config == nil {
		conf := map[string]string{
			"MAX_BUFFER_SIZE":       helpers.GetEnv("IRC_MAX_BUFFER_SIZE", "8192"),
			"IRC_MAX_USER_CHANNELS": helpers.GetEnv("IRC_MAX_USER_CHANNELS", "20"),
			"IRC_USER_MODES":        helpers.GetEnv("IRC_USER_MODES", "oiws"),
			"IRC_CHANNEL_MODES":     helpers.GetEnv("IRC_CHANNEL_MODES", "opsmt"),
		}
		config = &conf
	}

	// TODO - parse whether it should run w/ TLS or not which will determine the port used
	now := time.Now()

	creation_date_time := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		now.Location(),
	)

	is := IrcServer{
		DnsName:       dns_name,
		Version:       version,
		CreationDate:  creation_date_time,
		Addr:          addr,
		Role:          server_role,
		Listener:      nil,
		_ServerManger: nil,
		Servers:       *server_list,
		Clients:       *client_list,
		Config:        *config,
	}

	s_manager := NewServerManager(&is.Clients, &is.Servers)
	s_manager.Name = dns_name
	is._ServerManger = s_manager
	return &is
}

func (is *IrcServer) Run() {
	ln, err := net.Listen("tcp", is.Addr)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err.Error())
	}
	is.Listener = &ln

	log.Info().Msgf("Server started, role: %s, addr: %s, dns: %s, config: %v", is.Role, is.Addr, is.DnsName, is.Config)

	if is._ServerManger == nil {
		panic("Server's manager is null!!")
	}
	go is._ServerManger.Run()

	for {
		conn, err := (*is.Listener).Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		go is.handleConnection(&conn)
	}
}

func (is *IrcServer) Stop() {
	// TODO - probably do some shutdown stuff before just closing all connections
	(*is.Listener).Close()
}

func (is *IrcServer) handleConnection(conn *net.Conn) {
	defer (*conn).Close()

	remote_addr := (*conn).RemoteAddr()
	remote_ip, remote_port, err := net.SplitHostPort(remote_addr.String())
	if err != nil {
		log.Error().Str("remote_addr", remote_addr.String()).Msgf("Error splitting remote addr: %s", err.Error())
	}
	// TODO - resolve DNS name here for additional checks / verification
	// e.g. w/ servers and compare to server list
	remote_conn := rc.NewRemoteConn("", remote_ip, remote_port)
	client, session_timestamp, err := c.NewClient("", "", remote_conn, conn)
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		return
	}

	// add the client to the global client list
	is.Clients = append(is.Clients, client)

	log.Info().EmbedObject(client).Msgf("Client connected at %s", *session_timestamp)

	var max_buffer_size int
	max_buffer_size, err = strconv.Atoi((*is).Config["MAX_BUFFER_SIZE"])
	if err != nil {
		log.Error().EmbedObject(client).Msgf(err.Error())
		max_buffer_size = 8192
	}

	for {
		// TODO - clear buffers so no extra data is sent?
		// TODO - send periodic PING commands
		//		  if no PONG is received, terminate the connection
		//		  used to determine dead connections

		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		recv_buf := make([]byte, max_buffer_size)
		_, err := io.ReadAtLeast((*conn), recv_buf, 1)
		// _, err := (*conn).Read(recv_buf)
		if err != nil {
			if err == io.EOF {
				end_timestamp, err := client.SetSessionEndTimestamp()
				if err != nil {
					log.Error().EmbedObject(client).Msg(err.Error())
				}
				log.Info().EmbedObject(client).Msgf("Client disconnected: %s", *end_timestamp)
				// TODO - remove client state from ServerManager!!!
				// TODO - write session duration as a metric / possible analysis
			} else {
				log.Error().EmbedObject(client).Msgf("Error reading data from connection: %s", err.Error())
			}

			//if err == io.ErrShortBuffer
			return
		}

		server_metadata := map[string]string{
			"name":         is.DnsName,
			"version":      is.Version,
			"date":         is.CreationDate.String(),
			"usermodes":    is.Config["IRC_USER_MODES"],
			"channelmodes": is.Config["IRC_CHANNEL_MODES"],
		}

		response := chat.ProcessMessage(&recv_buf, client, is._ServerManger.Task_runner, server_metadata)

		if len(response) == 0 {
			// e.g. sometimes the server doesn't send anything back to the client
			// 		as in the case of correct password via PASS
			// reset the buffer
			continue
		}

		if _, err := (*conn).Write(response); err != nil {
			log.Error().EmbedObject(client).Msgf("Error writing to client: %s", err.Error())
			break
		}

		if strings.Contains(string(response), "ERROR") {
			log.Error().EmbedObject(client).Msgf("Critical error occurred. Closing client connection: %s", response)
			(*conn).Close()
			break
		}
	}
}
