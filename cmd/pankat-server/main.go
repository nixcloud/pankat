package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocraft/web"
	"golang.org/x/net/websocket"
	"net/http"
	"pankat"
	"pankat-server/ws"
)

type Context struct {
	//     HelloCount int
}

func onArticleChange(wsServer *ws.Server) func(string, string) {
	return func(srcFileName string, RenderedArticle string) {
		fmt.Println(srcFileName)
		//if srcFileName == "docker_compose_vs_nixcloud.mdwn" {
		wsServer.SendAll(RenderedArticle)
		//}
	}
}

func main() {
	fmt.Println(color.GreenString("pankat-server"), "starting!")
	pankat.Init()
	wsServer := ws.NewServer()
	ona := onArticleChange(wsServer)
	pankat.OnArticleChange(ona)
	go wsServer.Listen()
	go fsNotifyWatchDocumentsDirectory(wsServer, pankat.GetConfig().DocumentsPath)
	router := web.New(Context{}). // Create your router
					Middleware(web.LoggerMiddleware).
					Middleware(web.ShowErrorsMiddleware).
		//Middleware(web.StaticMiddleware("../output")).
		Middleware(web.StaticMiddleware(pankat.GetConfig().DocumentsPath)).
		Get("/websocket", func(rw web.ResponseWriter, req *web.Request) {
			websocket.Handler(wsServer.OnConnected).ServeHTTP(rw, req.Request)
		}).
		Get("/", redirectTo("/index.html"))
	http.ListenAndServe(pankat.GetConfig().ListenAndServe, router) // wait until ctrl+c
}

func redirectTo(to string) func(web.ResponseWriter, *web.Request) {
	return func(rw web.ResponseWriter, req *web.Request) {
		http.Redirect(rw, req.Request, to, http.StatusFound)
	}
}
