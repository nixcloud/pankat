package ws

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

const channelBufSize = 100

// Chat client.
type Client struct {
	ws       *websocket.Conn
	server   *Server
	ch       chan *string
	doneCh   chan bool
	registry *Registry
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server, registry *Registry) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}
	if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan *string, channelBufSize)
	doneCh := make(chan bool)

	return &Client{ws, server, ch, doneCh, registry}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *string) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client disconnected.")
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			//fmt.Println("-------------------------------------")
			//fmt.Println("Send:", *msg)
			websocket.JSON.Send(c.ws, *msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg string
			err := websocket.Message.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				fmt.Println("client.go: Receive ws registry message:", msg)
				c.registry.Add(msg, c)
			}
		}
	}
}
