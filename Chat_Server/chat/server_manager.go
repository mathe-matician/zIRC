package chat

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"zirc/helpers"

	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
)

var server_manager_commands = map[string]string{
	"PASS":    "",
	"CAPAB":   "",
	"SERVER":  "",
	"CAP":     "",
	"SJOIN":   "",
	"PING":    "",
	"PONG":    "",
	"SQUIT":   "",
	"CONNECT": "",
	"BURST":   "",
	"EUID":    "",
	//	b. BURST / EUID (IRCv3)
	//
	// Purpose: Synchronize state after a netsplit or during initial connection.
	// Routing Table Use:
	// The server ensures that its routing table is updated with the correct paths for users and channels.
}

type ServerConnection struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Password     string `yaml:"password"`
	Passwordfile string `yaml:"password_file"`
	Conn         net.Conn
}

type ServerManagerConfig struct {
	ServerList map[string]ServerConnection
}

type ServerManager struct {
	Addr              string
	serverIpWhitelist map[string]string
	Config            ServerManagerConfig
}

func NewServerManager() *ServerManager {
	sm_config := NewServerManagerConfig()
	if len(sm_config.ServerList) != 0 {
		// do connections
		for _, server := range sm_config.ServerList {
			// TODO
			// try to connect to these servers
			go Connect(server)
		}

		// wait for any failures or successful connections?
		// select {
		// case pong_token := <-c.PingPongChan:
		// 	pongToken := pong_token[1:]
		// 	if pongToken != pingToken {
		// 		log.Info().EmbedObject(c).Msgf("%s != %s", pongToken, pingToken)
		// 		log.Info().EmbedObject(c).Msgf("%s len: %d", pongToken, len(pongToken))
		// 		log.Info().EmbedObject(c).Msgf("%s len: %d", pingToken, len(pingToken))
		// 		(*conn).Close()
		// 	}
		// 	log.Info().EmbedObject(c).Msgf("PONG successful")
		// case <-time.After(time.Duration(timeout) * time.Second):
		// 	log.Info().EmbedObject(c).Msgf("Server never received PONG, closing connection")
		// 	(*conn).Close()
		// }
	}

	return &ServerManager{
		Addr:   G_Config.S2S.Port,
		Config: sm_config,
	}
}

type ServerConnectionPkg struct {
	ServerName       string
	ConnectedServers []ServerConnection // todo this probably shouldn't contain other server passwords
	ClientList       []Client
	ChannelList      []Channel
}

// connect to another server
func Connect(connection ServerConnection) {

	// initial server handshake
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// - server name
	// - directly connected servers
	// - user and channel lists it manages
	// - network paths to reach other servers
	// e.g.
	//       A -- B -- C
	//			  |
	//			  D
	// When B connects to A, it says "I'm connced to C and D"
	// A could then hold:
	// []Servers{
	//		[]B {C, D}
	//  }
	// If multiple paths exist, then path selection considers:
	// 		- shortest path (least hops)
	//		- link cost (if the network assigns weights to connections)

	//
	// routing table updates
	//
}

func NewServerConnection() *ServerConnection {
	return &ServerConnection{}
}

func NewServerManagerConfig() ServerManagerConfig {
	var server_list map[string]ServerConnection
	config_path := helpers.GetEnv("IRC_S2S_CONFIG_FILE", "/chat_server/s2s_config.yaml")
	config_data, err := os.ReadFile(config_path)
	if err != nil {
		// probaby doens't have to be an error since we don't _have_ to have a config here
		// as of right now it only loads connected servers
		// _could_ support more in the future, though
		log.Error().Msgf("Error reading config file %v", err)
		return ServerManagerConfig{
			ServerList: server_list,
		}
	}

	// TODO
	// check for zero length config_data
	// no need to unmarshal then as no additional servers exist
	if len(config_data) == 0 {
		log.Info().Msgf("No data in server manager config")
		return ServerManagerConfig{
			ServerList: server_list,
		}
	}

	err = yaml.Unmarshal([]byte(config_data), &server_list)
	if err != nil {
		log.Error().Msgf("ServerManagerConfig couldn't unmarshal config: %v", err)
		panic(err)
	}

	log.Debug().Msgf("Server config list: %v", server_list)

	// TODO
	// should we connect to the servers here?
	for _, server := range server_list {
		if server.Passwordfile == "" {
			continue
		}

		server_password, err := os.ReadFile(server.Passwordfile)
		if err != nil {
			log.Error().Msgf("Error reading %s password file: %v", server.Host, err)
			continue
		}

		server.Password = string(server_password)
	}

	return ServerManagerConfig{
		ServerList: server_list,
	}
}

func (sm *ServerManager) init() {

}

func (sm *ServerManager) AddToServerIpWhiteList(ip, args string) {
	sm.serverIpWhitelist[ip] = args
}

func (sm *ServerManager) RemoveFromServerIpWhiteList(ip string) {
	delete(sm.serverIpWhitelist, ip)
}

// Allows this server to Join an existing IRC network
// TODO
// should this be an init function?
// e.g. get Join config
func (sm *ServerManager) Join(ip string) {
	/*
		1. PASS
		2. CAPAB (capability announcement)
			The server initiating the connection sends a CAPAB (capability) message to the receiving server. This message lists the features or protocol extensions that the server supports.

			CAPAB :QS EX IE KLN UNKLN

			Here, the server is advertising that it supports several features:

			QS: Quiet Channel Synchronization
			EX: Extended bans
			IE: Invite exceptions
			KLN: Kill line commands
			UNKLN: Undo kill line commands

			Each capability represents a specific feature or extension that allows enhanced communication between the servers.

		3. CAPAB matching
			Once the first server announces its capabilities, the receiving server also sends back a CAPAB message indicating its own supported capabilities.

			The two servers then compare their lists of supported capabilities. Both servers will use the features they both support for further communication. If a capability is not supported by one of the servers, it is not used in future communication.

		4. Negotiate features

			If both servers support advanced features (like extended bans, invite exceptions, etc.), they will use these capabilities in communication. For instance, if both servers support EX (extended bans), they will use extended ban syntax when communicating about bans across the network.

			If one of the servers does not support a specific feature, the two servers will revert to the basic, standard IRC protocol features for that area of communication.

		5. SERVER cmd

		https://docs.inspircd.org/server/examples/connection/
	*/

}

// Checks to see whether the server trying to connect to this network
// has entered in a valid password via the PASS command
func (sm *ServerManager) ValidS2SPassword(password string) bool {
	// check db first
	// if can't connect to db check locally for conf file or env variable
	// if all fails, then don't let join

	// last fallback is env var
	return G_Config.S2S.Password_file == password
}

func (sm *ServerManager) Run() {
	// listens on port 7000 for servers joining the irc network
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	ln := GetTCPListener(
		G_Config.S2S.Enable_tls,
		G_Config.S2S.Tls_cert_path,
		G_Config.S2S.Tls_key_path,
		G_Config.S2S.Tls_port,
		G_Config.Server.Host,
		G_Config.S2S.Port,
	)

	log.Info().Msgf("ServerManager listening: %s", sm.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Msgf("Error accepting connection: %s", err.Error())
			continue
		}

		if conn.RemoteAddr().Network() != "tcp" {
			conn.Close()
			continue
		}

		if !sm.whiteListedServerIp(conn.RemoteAddr().String()) {
			conn.Close()
			continue
		}

		go handleConnection(&conn, true)

		// go sm.handleConnection(&conn)
	}
}

// Checks to see if the connected IP is an expected server IP
// i.e. do we expect this IP to be communicating with the ServerManager
// Note: arg ip is expected to be ip:port actually.
func (sm *ServerManager) whiteListedServerIp(ip string) bool {
	server_ip, _, err := net.SplitHostPort(ip)
	if err != nil {
		log.Error().Msgf("ServerManager Error extracting IP: %s", err)
		return false
	}

	if _, ok := sm.serverIpWhitelist[server_ip]; !ok {
		log.Error().Msgf("ServerManager not a whitelisted server ip!: %s", server_ip)
		return false
	}
	return true
}

func (sm *ServerManager) handleConnection(conn *net.Conn) {
	defer (*conn).Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ServerManager recovered from panic:", r)
		}
	}()
	remote_addr := (*conn).RemoteAddr()
	log.Info().Msgf("Handling ServerManager connection: %s", remote_addr)

	// TODO
	// does this size differ for S2S?
	max_buffer_size := G_Config.S2S.Max_buffer_size

	for {
		if !sm.whiteListedServerIp((*conn).RemoteAddr().String()) {
			// need to always check if this IP is still valid as we handle this connection
			// as we could update the whitelist mid connection
			break
		}

		// block on read until the buffer has at least 1 byte.
		// just a hacky way for this to block as Read() doesn't block on its own
		recv_buf := make([]byte, max_buffer_size)
		_, err := io.ReadAtLeast((*conn), recv_buf, 1)
		if err != nil {
			if err == io.EOF {
				log.Info().Msgf("Server %s disconnected", remote_addr)
			}
			return
		}

		if _, err := (*conn).Write([]byte("")); err != nil {
			break
		}
	}
}

func (sm *ServerManager) Route(msg []byte, server string) {
	log.Debug().Msgf("ServerManager::Route::msg: %s", msg)
	// TODO
	// we already have a connection to all direct servers
	// we need to now send these details via some channel
	// since those connections are running in their own go routines.
}
