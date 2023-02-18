package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"pankat"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func getArticlesFilteredByDraftsExceptOne(eventRelFileName string) pankat.Articles {
	var _filtered pankat.Articles
	for _, e := range pankat.GetTargets(".") {
		if e.Draft == false || filepath.Clean(e.SrcDirectoryName+"/"+e.SrcFileName) == eventRelFileName {
			_filtered = append(_filtered, e)
		}
	}
	sort.Sort(_filtered)
	return _filtered
}

func fsNotifyWatchDocumentsDirectory(directory string) {
	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.FileInfo.IsDir() == false {
					if strings.HasSuffix(event.Name(), ".mdwn") {
						documentsPath, err := os.Getwd()
						if err != nil {
							log.Println(err)
						}
						eventRelFileName, _ := filepath.Rel(documentsPath, event.Path)
						if event.Op == watcher.Remove {
							fmt.Println("file removed:", event.Name())
						}
						if event.Op == watcher.Write || event.Op == watcher.Create {
							fmt.Println("File write|create detected in ", eventRelFileName)
							articles := getArticlesFilteredByDraftsExceptOne(eventRelFileName)
							for _, article := range articles {
								if filepath.Clean(article.SrcDirectoryName+"/"+article.SrcFileName) == eventRelFileName {
									fmt.Println("pankat.RenderPost(articles, article)")
									article.SourceReference = true
									article.WebsocketSupport = true
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
