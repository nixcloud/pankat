package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var articlesCache ArticlesCache

type Article struct {
	Title            string
	Article          []byte
	ModificationDate time.Time
	Summary          string
	Tags             []string
	Series           string
	Draft            bool
	SrcFileName      string
	DstFileName      string
	BaseFileName     string
	SrcDirectoryName string
	Anchorjs         bool
	Tocify           bool
	Timeline         bool
}

func (a Article) Render() string {
	// i would love to get rid of this initialization here and implement this 'constructor' like instead
	if articlesCache.Store == nil {
		//fmt.Println("Initializing hash map")
		articlesCache.Store = make(map[md5hash]string)
		articlesCache.load()
	}
	var text string = ""
	if articlesCache.Get(a) == "" {
		pandocProcess := exec.Command("pandoc", "-f", "markdown", "-t", "html5", "--highlight-style", "kate")
		stdin, err := pandocProcess.StdinPipe()
		if err != nil {
			fmt.Println(err)
			return "error rendering article"
		}
		buff := bytes.NewBufferString("")
		pandocProcess.Stdout = buff
		pandocProcess.Stderr = os.Stderr
		err1 := pandocProcess.Start()
		if err1 != nil {
			fmt.Println("An error occured: ", err1)
			return "error rendering article"
		}
		_, err2 := io.WriteString(stdin, string(a.Article))
		if err2 != nil {
			fmt.Println("An error occured: ", err2)
			return "error rendering article"
		}
		err3 := stdin.Close()
		if err3 != nil {
			fmt.Println("An error occured: ", err3)
			return "error rendering article"
		}
		err4 := pandocProcess.Wait()
		if err4 != nil {
			fmt.Println("An error occured during pandocProess wait: ", err4)
			return "error rendering article"
		}
		text = string(buff.Bytes())
		articlesCache.Set(a, text)
	} else {
		text = articlesCache.Get(a)
	}
	return text
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

type Articles []*Article

// https://gobyexample.com/sorting-by-functions
func (s Articles) Len() int {
	return len(s)
}
func (s Articles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Articles) Less(i, j int) bool {
	return s[i].ModificationDate.After(s[j].ModificationDate)
}

func (s Articles) NextArticle(a *Article) *Article {
	for i, elem := range s {
		if elem.Title == a.Title { // HACK
			if i-1 >= 0 {
				return s[i-1]
			}
		}
	}
	return nil
}

func (s Articles) PrevArticle(a *Article) *Article {
	for i, elem := range s {
		if elem.Title == a.Title { // HACK
			if i+1 < len(s) {
				return s[i+1]
			}
		}
	}
	return nil
}

func (s Articles) PrevArticleInSeries(a *Article) *Article {
	q := s.FilterBySeries(a.Series)
	if len(q) == 0 {
		return nil
	}
	z := q.PrevArticle(a)
	return z
}

func (s Articles) NextArticleInSeries(a *Article) *Article {
	q := s.FilterBySeries(a.Series)
	if len(q) == 0 {
		return nil
	}
	z := q.NextArticle(a)
	return z
}

func (s Articles) MakeRelativeLink(a *Article, b *Article) string {
	relativeSrcRootPath, _ := filepath.Rel(a.SrcDirectoryName, b.SrcDirectoryName)
	return relativeSrcRootPath
}

func (s Articles) TopLevel() Articles {
	var _filtered Articles
	for _, e := range s {
		if e.SrcDirectoryName == "" || e.SrcDirectoryName == "." {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) Posts() Articles {
	var _filtered Articles
	for _, e := range s {
		if e.SrcDirectoryName != "" && e.SrcDirectoryName != "." {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) FilterBySeries(t string) Articles {
	var _filtered Articles
	for _, e := range s {
		if e.Series == t {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) FilterByTag(t string) Articles {
	var _filtered Articles
	for _, e := range s {
		if contains(e.Tags, t) {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) FilterOutDrafts() Articles {
	var _filtered Articles
	for _, e := range s {
		if e.Draft == false {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

type MetaData struct {
	ArticleCount int
	Tags         map[string][]int
	Series       map[string][]int
	Years        map[int][]int
}

func (s Articles) CreateJSMetadata() MetaData {
	tagsMap := make(map[string][]int)
	seriesMap := make(map[string][]int)
	yearsMap := make(map[int][]int)
	for i, e := range s {
		m := e.ModificationDate
		year, err := strconv.Atoi(m.Format("2006"))
		if err == nil {
			if yearsMap[year] == nil {
				yearsMap[year] = []int{i}
			} else {
				yearsMap[year] = append(yearsMap[year], i)
			}
		}

		for _, t := range e.Tags {
			if tagsMap[t] == nil {
				tagsMap[t] = []int{i}
			} else {
				tagsMap[t] = append(tagsMap[t], i)
			}
		}
		z := s[i].Series
		if z != "" {
			if seriesMap[z] == nil {
				seriesMap[z] = []int{i}
			} else {
				seriesMap[z] = append(seriesMap[z], i)
			}
		}
	}
	return MetaData{len(s), tagsMap, seriesMap, yearsMap}
}
