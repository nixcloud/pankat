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

func getArticlesFilteredByDraftsExceptOne(eventRelFileName string) (pankat.Articles, error) {
	var _filtered pankat.Articles
	var found bool = false
	for _, e := range pankat.GetTargets(".") {
		if found == false {
			if filepath.Join(e.SrcDirectoryName, e.SrcFileName) == filepath.FromSlash(eventRelFileName) {
				found = true
				_filtered = append(_filtered, e)
			}
		}
		if e.Draft == false {
			_filtered = append(_filtered, e)
		}
	}
	if found == false {
		return pankat.Articles{}, fmt.Errorf("File %s not found in targets", eventRelFileName)
	}

	sort.Sort(_filtered)
	return _filtered, nil
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
							articles, err := getArticlesFilteredByDraftsExceptOne(eventRelFileName)
							if err != nil {
								fmt.Println(err)
								continue
							}
							for _, article := range articles {
								if filepath.Join(article.SrcDirectoryName, article.SrcFileName) == filepath.FromSlash(eventRelFileName) {
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
