package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"pankat"
	"path/filepath"
	"strings"
	"time"
)

func fsNotifyWatchDocumentsDirectory(directory string) {
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
							fmt.Println("Name, full path", event.Name(), event.Path)
							articles := pankat.GetArticles()
							for _, article := range articles {
								if article.SrcFileName == event.Name() {
									fmt.Println("pankat.RenderPost(articles, article)")
									pankat.RenderPost(articles, article)
								}
							}
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
