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

type Context struct{}

func onArticleChange(registry *ws.Registry) func(string, string) {
	return func(srcFileName string, RenderedArticle string) {
		fmt.Println(srcFileName)
		registry.OnArticleChange(srcFileName, RenderedArticle)
	}
}

func main() {
	fmt.Println(color.GreenString("pankat-server"), "starting!")
	pankat.Init()
	pankat.UpdateBlog()
	registry := ws.NewRegistry()

	server := ws.NewServer(registry)
	go server.Listen()

	ona := onArticleChange(registry)
	pankat.OnArticleChange(ona)
	go fsNotifyWatchDocumentsDirectory(pankat.GetConfig().DocumentsPath)
	router := web.New(Context{}). // Create your router
					Middleware(web.LoggerMiddleware).
					Middleware(web.ShowErrorsMiddleware).
		//Middleware(web.StaticMiddleware("../output")).
		Middleware(web.StaticMiddleware(pankat.GetConfig().DocumentsPath)).
		Get("/websocket", func(rw web.ResponseWriter, req *web.Request) {
			websocket.Handler(server.OnConnected).ServeHTTP(rw, req.Request)
		}).
		Get("/", redirectTo("/index.html"))
	http.ListenAndServe(pankat.GetConfig().ListenAndServe, router) // wait until ctrl+c
}

func redirectTo(to string) func(web.ResponseWriter, *web.Request) {
	return func(rw web.ResponseWriter, req *web.Request) {
		http.Redirect(rw, req.Request, to, http.StatusFound)
	}
}
