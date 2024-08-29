package server

import (
	"io"
	"net"
	"strconv"
	"strings"

	"zirc/chat"
	c "zirc/client"
	"zirc/helpers"
	rc "zirc/remote_conn"

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
	Addr          string
	Role          string
	Listener      *net.Listener
	_ServerManger *ServerManager
	Servers       []*IrcServer
	Clients       []*c.Client
	Config        map[string]string
}

type Task struct {
	id          uuid.UUID
	weight      float64
	task        string
	client_chan chan<- string
}

type Worker struct {
	id           uuid.UUID
	tasks        []*Task
	current_load float64
	quit         chan int
	frozen       bool
	// recv         chan<- *Task
}

type ServerManager struct {
	Name               string
	ClientList         *[]*c.Client
	ServerList         *[]*IrcServer
	WorkerPool         map[string]*Worker
	worker_tasks       chan string
	Task_runner        chan string
	results            chan string
	decreasing_workers bool
	// send        chan string
}

// task_weight_mapping := map[string]float64 {

// }

// Run starts the ServerManager which manages the Worker pool
func (sm *ServerManager) Run() {
	log.Info().Msg("ServerManager started")

	for {
		select {
		case client_task := <-sm.Task_runner:
			log.Info().Msgf("client_task: %s", client_task)
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

func (w *Worker) MarshalObject(e *log.Entry) {
	e.Str("id", w.id.String()).Float64("current_load", w.current_load)
}

func (sm *ServerManager) Debug() {
	log.Info().Msgf("ServerManager Debug!")
}

// Work
//
//	job: jobs received from the ServerManager
//	results: any results that are returned back to the ServerManager can be sent back to the client if needed
//
// TODO - this func may only need the ServerManager's ClientList and ServerList
func (w *Worker) Work(job chan string, results chan string, server_manager *ServerManager) {
	for {
		select {
		case j := <-job:
			log.Info().EmbedObject(w).Msgf("Job received: %s", j)
			server_manager.Debug()

			if server_manager.ClientList == nil {
				msg := "servermanager client list is null! idk how we got to this point..."
				log.Error().EmbedObject(w).Msg(msg)
				panic(msg)
			}

			log.Debug().EmbedObject(w).Msgf("Client List Len: %d", len(*server_manager.ClientList))
			for _, c := range *server_manager.ClientList {
				if c == nil {
					log.Debug().EmbedObject(w).Msgf("Client is null - trying next")
					continue
				}

				if c.Registered {
					log.Debug().EmbedObject(w).Msgf("Staring job: %s", j)
					if c.ClientConn == nil {
						log.Error().EmbedObject(w).Msg("Client connection is nil!!")
						break
					}
					conn := *(c.ClientConn)
					if _, err := conn.Write([]byte(j)); err != nil {
						log.Error().EmbedObject(c).Msgf("Error writing to client: %s", err.Error())
						break
					}
				} else {
					log.Info().EmbedObject(w).Msg("Client not registered")
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
		Task_runner:  make(chan string),
		worker_tasks: make(chan string),
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

func NewTask(task string) *Task {
	uid, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &Task{
		id:     uid,
		weight: 0.0,
		task:   task,
	}
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

func NewIrcServer(dns_name string, addr string, server_role string, server_list *[]*IrcServer, client_list *[]*c.Client, config *map[string]string) *IrcServer {
	if len(dns_name) == 0 {
		dns_name = helpers.GetEnv("IRC_SERVER_DNS_NAME", "localhost")
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
		}
		config = &conf
	}

	// TODO - parse whether it should run w/ TLS or not which will determine the port used

	is := IrcServer{
		DnsName:       dns_name,
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
	recv_buf := make([]byte, max_buffer_size)

	for {
		// TODO - clear buffers so no extra data is sent?
		// TODO - send periodic PING commands
		//		  if no PONG is received, terminate the connection
		//		  used to determine dead connections
		_, err := (*conn).Read(recv_buf)
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
			return
		}

		response := chat.ProcessMessage(&recv_buf, client, is._ServerManger.Task_runner)

		if len(response) == 0 {
			// e.g. sometimes the server doesn't send anything back to the client
			// 		as in the case of correct password via PASS
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
