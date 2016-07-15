package ws

// WARNING: i need to rewrite this code, license status is unclear, source: https://github.com/golang-samples/websocket

type Message struct {
  Author string `json:"author"`
  Body   string `json:"body"`
}

func (self *Message) String() string {
  return self.Author + " says " + self.Body
}