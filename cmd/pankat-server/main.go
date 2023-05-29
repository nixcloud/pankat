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
		fmt.Println("onArticleChange: ", srcFileName)
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
	go fsNotifyWatchDocumentsDirectory(pankat.Config().DocumentsPath)

	router := web.New(Context{})
	router.Middleware(web.LoggerMiddleware)
	router.Middleware(web.ShowErrorsMiddleware)
	router.Middleware(web.StaticMiddleware(pankat.Config().DocumentsPath))
	router.Get("/websocket", func(rw web.ResponseWriter, req *web.Request) {
		websocket.Handler(server.OnConnected).ServeHTTP(rw, req.Request)
	})

	// used to check from webpage if this is served by pankat-server
	router.Get("/pankat-server", func(rw web.ResponseWriter, req *web.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	router.Get("/draft", func(rw web.ResponseWriter, req *web.Request) {
		articles := pankat.GetArticles(".")
		articleQueryName := req.URL.Query().Get("article")
		if articleQueryName == "" {
			var draftList string
			draftList += "<p>this is a list of all drafts. Click on one to edit!</p>"
			draftList += "<ul>"
			for _, article := range articles {
				if article.Draft == true {
					draftList += "<li><a href=\"/draft?article=" + article.SrcFileName + "\">" + article.SrcFileName + "</a></li>"
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
				if article.SrcFileName == filepath.FromSlash(articleQueryName) {
					article.WebsocketSupport = true
					article.SourceReference = true
					navTitleArticleHTML := pankat.GenerateNavTitleArticleSource(articles, *article, article.Render())
					standalonePageContent := pankat.GenerateStandalonePage(articles, *article, navTitleArticleHTML)
					rw.Write([]byte(standalonePageContent))
				}
			}
		}
	})
	router.Get("/", redirectTo("/index.html"))
	http.ListenAndServe(pankat.Config().ListenAndServe, router) // wait until ctrl+c
}

func redirectTo(to string) func(web.ResponseWriter, *web.Request) {
	return func(rw web.ResponseWriter, req *web.Request) {
		http.Redirect(rw, req.Request, to, http.StatusFound)
	}
}
