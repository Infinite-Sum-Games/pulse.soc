package sse

import (
	"fmt"
	"sync"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
)

type Client struct {
	Channel chan string
}

// Manage all the SSE clients connections, disconnections, braodcasts,
// and keep-live ticker every 15 seconds to avoid client-timeouts
type LiveServer struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan string
	Mutex      sync.Mutex
}

// Constructer for SSE server
func NewLiveServer() *LiveServer {
	return &LiveServer{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan string),
	}
}

// Initiate the event server as a goroutine to handle client connections,
// disconnections, broadcasts and keep-alives
func (lv *LiveServer) Start() {
	// Global ticker to send out keep-alive ticks to all clients
	keepAliveTicker := time.NewTicker(15 * time.Second)
	defer keepAliveTicker.Stop()

	for {
		select {
		case client := <-lv.Register:
			lv.Mutex.Lock()
			lv.Clients[client] = true
			lv.Mutex.Unlock()
			cmd.Log.Info(
				fmt.Sprintf("Client connected. Total clients: %d", len(lv.Clients)),
			)
		case client := <-lv.Unregister:
			lv.Mutex.Lock()
			delete(lv.Clients, client)
			close(client.Channel)
			lv.Mutex.Unlock()
			cmd.Log.Info(
				fmt.Sprintf("Client disconnected. Total clients: %d", len(lv.Clients)),
			)
		case msg := <-lv.Broadcast:
			lv.Mutex.Lock()
			// Broadcast the message to all the clients
			for client := range lv.Clients {
				select {
				case client.Channel <- msg:
				default:
					// If the client buffer is full then remove the client
					delete(lv.Clients, client)
					close(client.Channel)
				}
			}
			lv.Mutex.Unlock()
		case <-keepAliveTicker.C:
			// Broadcast keep-alives to all clients to avoid timeout
			lv.Mutex.Lock()
			for client := range lv.Clients {
				select {
				case client.Channel <- ": keep-alive\n\n":
				default:
					// If the client buffer is full then remove the client
					delete(lv.Clients, client)
					close(client.Channel)
				}
			}
			lv.Mutex.Unlock()
		}
	}
}
