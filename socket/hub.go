package socket

import (
	"log"
	"net/http"
	"time"
)

type Hub interface {
	Run()
	ServeWS(w http.ResponseWriter, r *http.Request)
	removeClient(client *WSClient)
	SendToAll(msg interface{})
}

type hub struct {
	clients map[*WSClient]struct{}

	register   chan *WSClient
	unregister chan *WSClient
}

func NewHub() Hub {
	return &hub{
		clients:    make(map[*WSClient]struct{}),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

func (h *hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &WSClient{hub: h, conn: conn, send: make(chan interface{}, 256)}
	c.hub.register <- c

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go c.writePump()
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}

		case c := <-h.unregister:
			h.removeClient(c)
		}
	}
}

func (h *hub) removeClient(c *WSClient) {
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
}

func (h *hub) SendToAll(msg interface{}) {
	for c := range h.clients {
		select {
		case c.send <- msg:
		default:
			h.removeClient(c)
		}
	}
}
