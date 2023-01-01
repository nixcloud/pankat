package main

import (
	"github.com/gocraft/web"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"pankat-server/ws"
)

// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/08.2.html
// https://github.com/golang-samples/websocket/blob/master/websocket-chat/src/chat/client.go

type Context struct {
	//     HelloCount int
}

// inotify events
// 2016/07/20 16:54:38 http: Accept error: accept tcp 127.0.0.1:8080: accept4: too many open files; retrying in 1s

func inotifyWatchDir(server *ws.Server, d string) {
	watcher, err := NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Watch(d)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case ev := <-watcher.Event:
			// send updats to client if changes happen
			server.SendAll("reload")
			log.Println("event:", ev)
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}

func main() {
	//   updateCh := make(chan string)
	server := ws.NewServer()

	go server.Listen()

	go inotifyWatchDir(server, "../output/posts/") // FIXME hardcoded path

	router := web.New(Context{}). // Create your router
					Middleware(web.LoggerMiddleware).     // Use some included middleware
					Middleware(web.ShowErrorsMiddleware). // ...
		//Middleware(web.StaticMiddleware("../output")).
		Middleware(web.StaticMiddleware("../output")). // FIXME hardcoded path
		Get("/websocket", func(rw web.ResponseWriter, req *web.Request) {
			websocket.Handler(server.OnConnected).ServeHTTP(rw, req.Request)
		}).
		Get("/", redirectTo("/index.html"))

	http.ListenAndServe("localhost:8080", router) // Start the server!
}

func redirectTo(to string) func(web.ResponseWriter, *web.Request) {
	return func(rw web.ResponseWriter, req *web.Request) {
		http.Redirect(rw, req.Request, to, http.StatusFound)
	}
}
