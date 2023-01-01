package chat

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
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
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			//fmt.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			// FIXME add pandoc here
			a := strings.Replace(msg.Body, "\\n", "\n", -1)
			d1 := []byte(a)

			// FIXME use TempFile here (qknight)
			err := ioutil.WriteFile("/tmp/input", d1, 0644)
			if err != nil {
				panic(err)
			}
			//       v := strings.Split("-i /tmp/input -o /tmp/output", " ")
			// markdown
			out, err := exec.Command("/run/current-system/sw/bin/pandoc", "--toc", "-t", "html5", "--highlight-style", "kate", "-i", "/tmp/input", "-o", "/tmp/output.html").Output()
			// mediawiki
			//out, err := exec.Command("/home/joachim/.nix-profile/bin/pandoc", "--toc", "-f", "mediawiki", "-t", "html5", "-s", "--highlight-style",  "kate", "-i", "/tmp/input", "-o", "/tmp/output.html").Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("The date is %s\n", out)

			f, err1 := ioutil.ReadFile("/tmp/output.html")
			if err1 != nil {
				log.Fatal(err1)
			}
			//t := strings.Replace(string(f), "{&quot;body&quot;:&quot;","" , -1)
			//t1 := strings.Replace(string(t), "&quot;}","" , -1)
			//msg.Body = string(t1)

			//log.Println(f)
			msg.Body = string(f)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
