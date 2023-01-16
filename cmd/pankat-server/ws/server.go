package ws

import (
	"golang.org/x/net/websocket"
	"log"
)

type Server struct {
	clients   map[*websocket.Conn]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan string
	doneCh    chan bool
	errCh     chan error
	registry  *Registry
}

func NewServer(r *Registry) *Server {
	clients := make(map[*websocket.Conn]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan string)
	doneCh := make(chan bool)
	errCh := make(chan error)
	registry := r

	return &Server{
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
		registry,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(msg string) {
	for _, c := range s.clients {
		c.Write(&msg)
	}
}

// websocket handler
func (s *Server) OnConnected(ws *websocket.Conn) {
	defer func() {
		err := ws.Close()
		if err != nil {
			s.errCh <- err
		}
	}()

	client := NewClient(ws, s, s.registry)
	s.Add(client)
	client.Listen()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.ws] = c
			log.Println("Now", len(s.clients), "clients connected.")
			//       s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			s.registry.Del(c)
			delete(s.clients, c.ws)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
