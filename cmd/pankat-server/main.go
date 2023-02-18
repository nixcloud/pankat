package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocraft/web"
	"golang.org/x/net/websocket"
	"net/http"
	"pankat"
	"pankat-server/ws"
	"path/filepath"
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

	onArticleChangeFunction := onArticleChange(registry)
	pankat.OnArticleChange(onArticleChangeFunction)
	go fsNotifyWatchDocumentsDirectory(pankat.GetConfig().DocumentsPath)

	router := web.New(Context{})
	router.Middleware(web.LoggerMiddleware)
	router.Middleware(web.ShowErrorsMiddleware)
	router.Middleware(web.StaticMiddleware(pankat.GetConfig().DocumentsPath))
	router.Get("/websocket", func(rw web.ResponseWriter, req *web.Request) {
		websocket.Handler(server.OnConnected).ServeHTTP(rw, req.Request)
	})
	router.Get("/draft", func(rw web.ResponseWriter, req *web.Request) {
		articles := pankat.GetTargets(".")
		articleQueryName := req.URL.Query().Get("article")
		if articleQueryName == "" {
			var draftList string
			draftList += "<p>this is a list of all drafts. Click on one to edit!</p>"
			draftList += "<ul>"
			for _, article := range articles {
				if article.Draft == true {
					aname := filepath.Clean(article.SrcDirectoryName + "/" + article.SrcFileName)
					draftList += "<li><a href=\"/draft?article=" + aname + "\">" + aname + "</a></li>"
				}
			}
			draftList += "</ul>"
			var article pankat.Article
			article.Title = "drafts"
			article.SpecialPage = true
			navTitleArticleHTML := pankat.GenerateNavTitleArticleSource(articles, article, draftList)
			standalonePageContent := pankat.GenerateStandalonePage(articles, article, navTitleArticleHTML)
			rw.Write([]byte(standalonePageContent))
		} else {
			for _, article := range articles {
				if filepath.Clean(article.SrcDirectoryName+"/"+article.SrcFileName) == articleQueryName {
					article.WebsocketSupport = true
					navTitleArticleHTML := pankat.GenerateNavTitleArticleSource(articles, *article, article.Render())
					standalonePageContent := pankat.GenerateStandalonePage(articles, *article, navTitleArticleHTML)
					rw.Write([]byte(standalonePageContent))
				}
			}
		}
	})
	router.Get("/", redirectTo("/index.html"))
	http.ListenAndServe(pankat.GetConfig().ListenAndServe, router) // wait until ctrl+c
}

func redirectTo(to string) func(web.ResponseWriter, *web.Request) {
	return func(rw web.ResponseWriter, req *web.Request) {
		http.Redirect(rw, req.Request, to, http.StatusFound)
	}
}
