package main

import (
    "path/filepath"
    "io/ioutil"
    "os"
    "time"
    "strings"
    "strconv"
    "fmt"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *ClientMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *ClientMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
            var filename string = filepath.Join(os.TempDir(), client.Channel + ".msg")
            content, err := ioutil.ReadFile(filename)

            if err == nil {
                var lastMessage string = string(content[:])
                result := strings.Split(lastMessage, ",")
                timeLastMessage, err := time.Parse(time.UnixDate, result[0])

                if err == nil {
                    originalMessageClientId, err := strconv.ParseInt(result[1], 10, 64)
                    if err == nil {
                        if  originalMessageClientId != client.Id {
                            // if last message was send 30 seconds later, send again
                            if (time.Now().Sub(timeLastMessage).Seconds() < 30) {
                                select {
                                case client.send <- []byte(result[2]):
                                default:
                                    close(client.send)
                                    delete(h.clients, client)
                                }
                            }
                        }
                    }
                }
            }

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
                if client.Channel == message.Channel {
                    if client.Id != message.From.Id {
                        select {
                        case client.send <- message.Body:
                        default:
                            close(client.send)
                            delete(h.clients, client)
                        }

                        var filename string = filepath.Join(os.TempDir(), message.Channel + ".msg")
                        ioutil.WriteFile(filename, []byte(time.Now().Format(time.UnixDate) + "," + fmt.Sprint(client.Id) + "," + string(message.Body)), 0755)
                    }
                }
			}
		}
	}
}
