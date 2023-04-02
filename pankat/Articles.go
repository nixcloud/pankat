package pankat

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var articlesCache ArticlesCache

type Article struct {
	Title             string
	ArticleMDWNSource []byte
	ModificationDate  time.Time
	Summary           string
	Tags              []string
	Series            string
	SrcFileName       string // foo.mdwn
	SrcDirectoryName  string // /home/user/documents (lacks foo.mdwn)
	DstFileName       string // /home/user/documents/foo.html
	SpecialPage       bool   // used for timeline.html, about.html (not added to timeline if true, not added in list of articles)
	Draft             bool
	Anchorjs          bool
	Tocify            bool
	Timeline          bool // generating timeline.html uses this flag in RenderTimeline(..)
	SourceReference   bool // switch for showing the document source mdwn at bottom
	WebsocketSupport  bool // live update support via WS on/off
}

func PandocMarkdown2HTML(articleMarkdown []byte) (string, error) {
	pandocProcess := exec.Command("pandoc", "-f", "markdown", "-t", "html5", "--highlight-style", "kate")
	stdin, err := pandocProcess.StdinPipe()
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return "", err
	}
	buff := bytes.NewBufferString("")
	pandocProcess.Stdout = buff
	pandocProcess.Stderr = os.Stderr
	err1 := pandocProcess.Start()
	if err1 != nil {
		fmt.Println("An error occurred: ", err1)
		return "", err1
	}
	_, err2 := io.WriteString(stdin, string(articleMarkdown))
	if err2 != nil {
		fmt.Println("An error occurred: ", err2)
		return "", err2
	}
	err3 := stdin.Close()
	if err3 != nil {
		fmt.Println("An error occurred: ", err3)
		return "", err3
	}
	err4 := pandocProcess.Wait()
	if err4 != nil {
		fmt.Println("An error occurred during pandocProess wait: ", err4)
		fmt.Println("An error occurred: ", err4)
	}
	return string(buff.Bytes()), nil
}

func (a Article) Render() string {
	// FIXME i would love to get rid of this initialization here and implement this 'constructor' like instead
	if articlesCache.Store == nil {
		//fmt.Println("Initializing hash map")
		articlesCache.Store = make(map[md5hash]string)
		articlesCache.load()
	}
	var text string = ""
	if articlesCache.Get(a) == "" {
		if GetConfig().Verbose > 1 {
			fmt.Println(color.YellowString("pandoc run for article"), a.DstFileName)
		}
		text, err := PandocMarkdown2HTML(a.ArticleMDWNSource)
		if err != nil {
			fmt.Println("An error occurred during pandoc pipeline run: ", err)
			panic(err)
		}
		articlesCache.Set(a, text)
	} else {
		fmt.Println(color.YellowString("cache hit, no pandoc run for article"), a.DstFileName)
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

func (s Articles) GetTitleNAV(article *Article) string {
	articles := s
	//   fmt.Println("---------------")
	titleNAV := ""
	//   fmt.Println(article.Title)
	p := articles.PrevArticle(article)
	if p != nil {
		// link is active
		titleNAV +=
			`<span id="articleNavLeft"> <a href="` + p.DstFileName + `"> 
      <span class="glyphiconLink glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article"> </span> prev. article
    </a> </span>`
	}
	n := articles.NextArticle(article)
	if n != nil {
		// link is active
		titleNAV +=
			`<span id="articleNavRight"><a href="` + n.DstFileName + `"> 
        next article <span class="glyphiconLink glyphicon glyphicon-chevron-right" aria-hidden="true" title="next article"></span>
    </a> </span>`
	}

	return titleNAV
}

func (s Articles) GetSeriesNAV(article *Article) string {
	articles := s
	seriesNAV := ""
	var sPrev string
	var sNext string

	if article.Series != "" {
		sp := articles.PrevArticleInSeries(article)
		if sp != nil {
			sPrev = sp.DstFileName
		}

		sn := articles.NextArticleInSeries(article)
		if sn != nil {
			sNext = sn.DstFileName
		}
		seriesNAV =
			`
      <div id="seriesContainer">
      <a href="timeline.html?filter=series::` + article.Series + `" title="article series ` + article.Series + `" class="seriesbtn btn btn-primary">` +
				article.Series + `</a>
        <header class="seriesHeader">
          <div id="seriesLeft">`
		if sp != nil {
			seriesNAV += `<a href="` + sPrev + `">` +
				`<span class="glyphiconLinkSeries glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article in series"></span>
            </a> `
		}
		seriesNAV += `  </div>
          <div id="seriesRight">`
		if sn != nil {
			seriesNAV += `   <a href="` + sNext + `">
              <span class="glyphiconLinkSeries glyphicon glyphicon-chevron-right" aria-hidden="true" title="next article in series"></span>
            </a>`
		}
		seriesNAV += `</div>
        </header>
      </div>`
	}
	return seriesNAV
}

func (s Articles) Targets() Articles {
	var _filtered Articles
	for _, e := range s {
		e.SourceReference = true
		e.WebsocketSupport = true
		if e.SpecialPage == true {
			e.Anchorjs = false
			e.Tocify = false
		} else {
			e.Anchorjs = true
			e.Tocify = true
		}
		_filtered = append(_filtered, e)
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

func (s Articles) FilterOutSpecialPages() Articles {
	var _filtered Articles
	for _, e := range s {
		if e.SpecialPage == false {
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
