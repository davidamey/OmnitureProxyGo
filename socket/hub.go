package socket

import (
	"log"
	"net/http"
)

type Hub interface {
	Run()
	ServeWS(w http.ResponseWriter, r *http.Request)
	RemoveClient(client *WSClient)
	SendToAll(msg interface{})
}

type hub struct {
	inbound chan interface{}

	clients []*WSClient

	register   chan *WSClient
	unregister chan *WSClient
}

func NewHub() Hub {
	return &hub{
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

func (h *hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("[Hub] new client")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &WSClient{hub: h, conn: conn, send: make(chan interface{}, 256)}
	c.hub.register <- c
	go c.writePump()
	c.readPump()
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients = append(h.clients, c)

		case c := <-h.unregister:
			h.RemoveClient(c)
		}
	}
}

func (h *hub) RemoveClient(client *WSClient) {
	for i, c := range h.clients {
		if c == client {
			close(c.send)

			// delete from slice
			copy(h.clients[i:], h.clients[i+1:])
			h.clients[len(h.clients)-1] = nil
			h.clients = h.clients[:len(h.clients)-1]
			break
		}
	}
}

func (h *hub) SendToAll(msg interface{}) {
	for _, c := range h.clients {
		select {
		case c.send <- msg:
		default:
			h.RemoveClient(c)
		}
	}
}
