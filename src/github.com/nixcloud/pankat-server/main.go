package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gocraft/web"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"pankat-server/ws"
	"path/filepath"
)

type Context struct {
	//     HelloCount int
}

func fsNotifyWatchDocumentsDirectory(wsServer *ws.Server, directory string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		//fmt.Println("fsNotifyWatchDocumentsDirectory: go func()")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
				if event.Has(fsnotify.Create) {
					log.Println("created file:", event.Name)
				}
				//Remove FIXME
				//Rename FIXME
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	// Add a path.
	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("fsNotifyWatchDocumentsDirectory started")
	//for {
	//	select {
	//	case ev := <-watcher.Event:
	//		// send updats to client if changes happen
	//		wsServer.SendAll("reload")
	//		log.Println("event:", ev)
	//	case err := <-watcher.Error:
	//		log.Println("error:", err)
	//	}
	//}
	<-make(chan struct{})
}

var inputPath string
var outputPath string
var SiteURL string
var SiteTitle string

func main() {

	pflag.String("input", "documents", "input directory  ('documents'') in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files")
	pflag.String("output", "output", "output directory ('output') all generated files will be stored there and all directories like css/ js/ images and fonts/ will be rsynced there")
	pflag.String("siteURL", "https://lastlog.de/blog", "The URL of the blog, for example: 'https://example.com/blog'")
	pflag.String("siteTitle", "lastlog.de/blog", "Title which is inserted top left, for example: 'lastlog.de/blog'")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	input := viper.GetString("input")
	output := viper.GetString("output")
	SiteURL = viper.GetString("siteURL")
	SiteTitle = viper.GetString("siteTitle")

	i1, err := filepath.Abs(input)
	inputPath = i1

	o1, err := filepath.Abs(output)
	outputPath = o1

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("pankat-server starting")
	fmt.Println("input Directory: ", inputPath)
	fmt.Println("output Directory: ", outputPath)

	//   updateCh := make(chan string)
	wsServer := ws.NewServer()

	go wsServer.Listen()
	go fsNotifyWatchDocumentsDirectory(wsServer, inputPath)

	router := web.New(Context{}). // Create your router
					Middleware(web.LoggerMiddleware).
					Middleware(web.ShowErrorsMiddleware).
		//Middleware(web.StaticMiddleware("../output")).
		Middleware(web.StaticMiddleware(outputPath)).
		Get("/websocket", func(rw web.ResponseWriter, req *web.Request) {
			websocket.Handler(wsServer.OnConnected).ServeHTTP(rw, req.Request)
		}).
		Get("/", redirectTo("/index.html"))

	http.ListenAndServe("localhost:8000", router)
	// wait until ctrl+c
}

func redirectTo(to string) func(web.ResponseWriter, *web.Request) {
	return func(rw web.ResponseWriter, req *web.Request) {
		http.Redirect(rw, req.Request, to, http.StatusFound)
	}
}
