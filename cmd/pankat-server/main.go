package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocraft/web"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"pankat"
	"pankat-server/ws"
	"pankat/db"
	"path/filepath"
	"time"
)

type Context struct{}

func onArticleChange(registry *ws.Registry) func(string, string) {
	return func(dstFileName string, RenderedArticle string) {
		fmt.Println("onArticleChange: ", dstFileName)
		registry.OnArticleChange(dstFileName, RenderedArticle)
	}
}

// FIXME for debugging only
func fsWriter() {
	time.Sleep(1 * time.Second)
	fPath := pankat.Config().DocumentsPath
	os.WriteFile(filepath.Join(fPath, "fsWriter.mdwn"), []byte("fsWriter test"), 0644)

	var i uint = 0
	for true {
		s := fmt.Sprintf("[[!draft]]\nfsWriter test\n %d", i)
		os.WriteFile(filepath.Join(fPath, "fsWriter.mdwn"), []byte(s), 0644)
		time.Sleep(10 * time.Minute)
		i += 1
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

	go fsWriter() // FIXME for debugging only

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
		drafts, _ := db.Instance().Drafts()
		articleQueryName := req.URL.Query().Get("article")
		if articleQueryName == "" {
			var draftList string
			draftList += "<p>this is a list of all drafts. Click on one to edit!</p>"
			draftList += "<ul>"
			for _, article := range drafts {
				if article.Draft == true {
					draftList += "<li><a href=\"/draft?article=" + article.SrcFileName + "\">" + article.SrcFileName + "</a></li>"
				}
			}
			draftList += "</ul>"
			var article db.Article
			article.Title = "drafts"
			article.SpecialPage = true
			navTitleArticleHTML := pankat.GenerateNavTitleArticleSource(article, draftList)
			standalonePageContent := pankat.GenerateStandalonePage(article, navTitleArticleHTML)
			rw.Write([]byte(standalonePageContent))
		} else {
			for _, draft := range drafts {
				if draft.SrcFileName == filepath.FromSlash(articleQueryName) {
					newArticle, errCreate := pankat.CreateArticleFromFilesystemMarkdown(draft.SrcFileName)
					if errCreate != nil {
						fmt.Println(errCreate)
						db.Instance().Del(draft.SrcFileName)
						break
					}
					dbArticle, _, errSet := db.Instance().Set(newArticle)
					if errSet != nil {
						fmt.Println(errSet)
						break
					}
					body := pankat.Render(*dbArticle)
					navTitleArticleHTML := pankat.GenerateNavTitleArticleSource(*dbArticle, body)
					standalonePageContent := pankat.GenerateStandalonePage(*dbArticle, navTitleArticleHTML)
					rw.Write([]byte(standalonePageContent))
					return
				}
			}
			http.Redirect(rw, req.Request, "/draft", http.StatusFound)
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
