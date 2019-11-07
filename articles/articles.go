package pankat

import (
	"crypto/md5"
	"path/filepath"
	"time"
)

type Article struct {
	Title            string
	SrcFileName      string
	DstFileName      string
	BaseFileName     string
	SrcDirectoryName string

	Tags             []string
	Series           string
	ModificationDate time.Time
	Hash             [md5.Size]byte
	Article          []byte
	RenderedArticle  string
	Draft            bool
	Summary          string
	Anchorjs         bool
	Tocify           bool
	Timeline         bool
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
	//   Series map[string][]int
}

func (s Articles) TagUsage() MetaData {
	tagsMap := make(map[string][]int)
	for i, e := range s {
		for _, t := range e.Tags {
			if tagsMap[t] == nil {
				tagsMap[t] = []int{i}
			} else {
				tagsMap[t] = append(tagsMap[t], i)
			}
		}
	}
	return MetaData{len(s), tagsMap}
}
