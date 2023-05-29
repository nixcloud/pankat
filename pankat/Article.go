package pankat

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

type Article struct {
	Title             string
	ArticleMDWNSource []byte
	ModificationDate  time.Time
	Summary           string
	Tags              []string
	Series            string
	SrcFileName       string // /home/user/documents/foo.mdwn
	DstFileName       string // /home/user/documents/foo.html
	//NextArticle         *Article
	//PrevArticle         *Article
	//NextArticleInSeries *Article
	//PrevArticleInSeries *Article
	SpecialPage      bool // used for timeline.html, about.html (not added to timeline if true, not added in list of articles)
	Draft            bool
	Anchorjs         bool
	Tocify           bool
	Timeline         bool // generating timeline.html uses this flag in RenderTimeline(..)
	SourceReference  bool // switch for showing the document source mdwn at bottom
	WebsocketSupport bool // live update support via WS on/off
}

//func (a Article) Hash() md5hash {
//	bytes, err := json.Marshal(a)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return md5.Sum(bytes)
//}

func (a Article) Render() string {
	if articlesCache.Store == nil {
		//fmt.Println("Initializing hash map")
		articlesCache.Store = make(map[md5hash]string)
		articlesCache.load()
	}
	if articlesCache.Get(a) == "" {
		if Config().Verbose > 1 {
			fmt.Println(color.YellowString("pandoc run for article"), a.DstFileName)
		}
		text, err := PandocMarkdown2HTML(a.ArticleMDWNSource)
		if err != nil {
			fmt.Println("An error occurred during pandoc pipeline run: ", err)
			panic(err)
		}
		articlesCache.Set(a, text)
		return text
	} else {
		fmt.Println(color.YellowString("cache hit, no pandoc run for article"), a.DstFileName)
		return articlesCache.Get(a)
	}
}
