package ws

// WARNING: i need to rewrite this code, license status is unclear, source: https://github.com/golang-samples/websocket

import (
	"golang.org/x/net/websocket"
	"log"
)

// Chat server.
type Server struct {
	//   pattern   string
	//   messages  []string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan string
	doneCh    chan bool
	errCh     chan error
	//   updateCh  chan string
}

// Create new chat server.
func NewServer() *Server {
	//   messages := []string{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan string)
	doneCh := make(chan bool)
	errCh := make(chan error)
	//   updateCh = updateCh

	return &Server{
		//     pattern,
		//     messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
		//     updateCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg string) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

// func (s *Server) sendPastMessages(c *Client) {
//   for _, msg := range s.messages {
//     c.Write(msg)
//   }
// }

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

	client := NewClient(ws, s)
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
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			//       s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			//       s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
