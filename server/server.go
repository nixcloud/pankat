package main

import (
    "github.com/gocraft/web"
    "fmt"
    "net/http"
    "strings"
//     "golang.org/x/net/websocket"
//     "io"
    "log"
    "./ws"
)

// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/08.2.html
// https://github.com/golang-samples/websocket/blob/master/websocket-chat/src/chat/client.go

type Context struct {
    HelloCount int
}

func (c *Context) SetHelloCount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
    c.HelloCount = 51
    next(rw, req)
}

func (c *Context) SayHello(rw web.ResponseWriter, req *web.Request) {
    fmt.Fprint(rw, strings.Repeat("Hello ", c.HelloCount), "World!")
}

// inotify events
func inotifyWatchDir(d string) {
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
            // FIXME do it
            log.Println("event:", ev)
        case err := <-watcher.Error:
            log.Println("error:", err)
        }
    }
}
  
func main() {
  
  server := ws.NewServer("/websocket")
  go server.Listen()
  go http.ListenAndServe("localhost:12345", nil)   // Start the WS server!
  
  go inotifyWatchDir("output"); // FIXME hardcoded path
  
  router := web.New(Context{}).                   // Create your router
        Middleware(web.LoggerMiddleware).           // Use some included middleware
        Middleware(web.ShowErrorsMiddleware).       // ...
        //Middleware(web.StaticMiddleware("../output")).
        Middleware(web.StaticMiddleware("output")). // FIXME hardcoded path
        Middleware((*Context).SetHelloCount).       // Your own middleware!
        Get("/", (*Context).SayHello)               // Add a route
    http.ListenAndServe("localhost:8080", router)   // Start the server!
    
}

