package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocraft/web"
	"github.com/radovskyb/watcher"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"pankat-server/ws"
	"path/filepath"
	"time"
)

type Context struct {
	//     HelloCount int
}

func fsNotifyWatchDocumentsDirectory(wsServer *ws.Server, directory string) {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	//w.SetMaxEvents(1)

	// Only notify rename and move events.
	//w.FilterOps(watcher.Create, watcher.Write)

	// Only files that match the regular expression during file listings
	// will be watched.
	//r := regexp.MustCompile("^abc$")
	//w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				wsServer.SendAll("reload")
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(directory); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	//for path, f := range w.WatchedFiles() {
	//	fmt.Printf("%s: %s\n", path, f.Name())
	//}

	//fmt.Println()

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
	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
	//<-make(chan struct{})
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

	fmt.Println(color.GreenString("pankat-server"), "starting!")
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
