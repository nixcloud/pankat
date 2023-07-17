package ws

import "fmt"

type Registry struct {
	//map["foo.mdwn"]={client1, client2, client3}
	relations map[string][]*Client
}

func NewRegistry() *Registry {
	relations := make(map[string][]*Client)

	return &Registry{relations}
}

func (r *Registry) Add(article string, client *Client) {
	for _, c := range r.relations[article] {
		if c == client {
			return
		}
	}
	fmt.Println("Adding client to article", article)
	r.relations[article] = append(r.relations[article], client)
}

func (r *Registry) Del(client *Client) {
	for article, _ := range r.relations {
		for i, c := range r.relations[article] {
			if c == client {
				fmt.Println("Removing client from article", article)
				r.relations[article] = append(r.relations[article][:i], r.relations[article][i+1:]...)
			}
		}
	}
}

func (r *Registry) OnArticleChange(dstFileName string, RenderedArticle string) {
	for _, c := range r.relations[dstFileName] {
		c.Write(&RenderedArticle) // sends the string via the websocket to related clients
	}
}
