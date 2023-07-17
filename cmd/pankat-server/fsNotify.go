package main

import (
	"fmt"
	"github.com/qknight/watcher"
	"log"
	"os"
	"pankat"
	"pankat/db"
	"path/filepath"
	"strings"
	"time"
)

func fsRemoveMDWN(eventRelFileName string) {
	fmt.Println("File removed:", eventRelFileName)
	affectedArticles, err := db.Instance().Del(eventRelFileName)
	if err == nil {
		pankat.RenderPostsBySrcFileNames(affectedArticles)
		pankat.RenderTimeline()
		pankat.SetMostRecentArticle()
	} else {
		fmt.Println("File removal error: ", err)
	}
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
						// note: remove/rename is not trigged by a drag'n'drop of a file in the explorer -> no clue
						if event.Op == watcher.Remove {
							fmt.Println("File remove detected in: ", eventRelFileName)
							fsRemoveMDWN(eventRelFileName)
						}
						if event.Op == watcher.Rename {
							eventOldRelFileName, _ := filepath.Rel(documentsPath, event.OldPath)
							fmt.Println("File rename detected in: ", eventOldRelFileName)
							fsRemoveMDWN(eventOldRelFileName)
						}
						if event.Op == watcher.Write || event.Op == watcher.Create || event.Op == watcher.Rename {
							fmt.Println("File write|create|rename detected in: ", eventRelFileName)
							newArticle, err := pankat.CreateArticleFromFilesystemMarkdown(eventRelFileName)
							if err != nil {
								fmt.Println(err)
								break
							}
							dbArticle, affectedArticles, errNew := db.Instance().Set(newArticle)
							if errNew != nil {
								fmt.Println(errNew)
							} else {
								if dbArticle == nil {
									fmt.Println("error: dbArticle is nil")
									break
								}
								pankat.RenderPost(dbArticle)
								pankat.RenderPostsBySrcFileNames(affectedArticles)
								pankat.RenderTimeline()
								pankat.SetMostRecentArticle()
							}
							break
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
		log.Fatalln(err)
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
