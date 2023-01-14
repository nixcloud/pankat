package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocraft/web"
	"github.com/radovskyb/watcher"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"pankat"
	"pankat-server/ws"
	"path/filepath"
	"strings"
	"time"
)

type Context struct {
	//     HelloCount int
}

func fsNotifyWatchDocumentsDirectory(wsServer *ws.Server, directory string) {
	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.FileInfo.IsDir() == false {
					if strings.HasSuffix(event.Name(), ".mdwn") {
						if event.Op == watcher.Remove {
							fmt.Println("file removed:", event.Name())
						}
						if event.Op == watcher.Write || event.Op == watcher.Create {
							fmt.Println("Name", event.Name())
							//fmt.Println("Path", event.Path)
							//wsServer.SendAll("reload")
							//wsServer.SendAll(pankat.PandocMarkdown2HTML("")
							pankat.UpdateBlog(true)
						}
					}
				}
				if event.FileInfo.IsDir() == true {
					if event.Op == watcher.Remove {
						w.Remove(event.Path)
					}
					if event.Op == watcher.Create {
						w.Add(event.Path)
					}
				}
			case err := <-w.Error:
				if err == watcher.ErrWatchedFileDeleted {
					continue
				}
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	walkFunc := watchDir(w)
	if err := filepath.Walk(directory, walkFunc); err != nil {
		fmt.Println("ERROR", err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func watchDir(w *watcher.Watcher) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return w.Add(path)
		}
		return nil
	}
}

func onArticleChange(wsServer *ws.Server) func(string, string) {
	return func(srcFileName string, RenderedArticle string) {
		fmt.Println(srcFileName)
		if srcFileName == "docker_compose_vs_nixcloud.mdwn" {
			wsServer.SendAll(RenderedArticle)
		}
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
