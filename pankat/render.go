package pankat

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"pankat/db"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func RenderPostsBySrcFileNames(articles []string) {
	//fmt.Println("affectedArticles: ", articles)
	for _, srcFileName := range articles {
		fmt.Println("article: ", srcFileName)
		article, err := db.Instance().QueryRawBySrcFileName(srcFileName)
		if err == nil {
			fmt.Println(color.YellowString("Rendering post: "), srcFileName)
			RenderPost(article)
		} else {
			fmt.Println(color.RedString("Error rendering post: "), srcFileName, err)
		}
	}
}

func RenderPosts(articlesAll []db.Article) {
	defer timeElapsed("RenderPosts")()
	fmt.Println(color.YellowString("Rendering posts"))
	for _, article := range articlesAll {
		RenderPost(&article)
	}
}

func RenderPost(article *db.Article) {
	if Config().Verbose > 0 {
		fmt.Println("Rendering article '" + article.Title + "'")
	}
	body := Render(*article)
	navTitleArticleHTML := GenerateNavTitleArticleSource(*article, body)
	standalonePageContent := GenerateStandalonePage(*article, navTitleArticleHTML)

	sendLiveUpdateViaWS(article.DstFileName, navTitleArticleHTML)

	if (*article).Draft == true {
		fmt.Println("Article is a draft, not writing to disk: '" + article.DstFileName + "'")
		outName := filepath.Join(Config().DocumentsPath, article.DstFileName)
		if _, err := os.Stat(outName); err == nil {
			err := os.Remove(outName)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
		return
	}
	errMkdir := os.MkdirAll(Config().DocumentsPath, 0755)
	if errMkdir != nil {
		fmt.Println(errMkdir)
		panic(errMkdir)
	}
	// write to disk
	outName := filepath.Join(Config().DocumentsPath, article.DstFileName)
	err5 := os.WriteFile(outName, standalonePageContent, 0644)
	if err5 != nil {
		fmt.Println(err5)
		panic(article)
	}
	fmt.Println("Article on disk updated: '" + article.DstFileName + "'")
}

func GenerateStandalonePage(article db.Article, navTitleArticleSource string) []byte {
	buff := bytes.NewBufferString("")
	t, err := template.New("standalonePage.tmpl").
		ParseFiles("templates/standalonePage.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	noItems := struct {
		Title                 string
		SiteBrandTitle        string
		Anchorjs              bool
		Tocify                bool
		Timeline              bool
		NavTitleArticleSource string
		ArticleSourceCodeURL  string // file location from the web
		ArticleSourceCodeFS   string // file location on disk (win/linux/...) where pankat-server runs
		ArticleDstFileName    string // roadmap.html
		ShowSourceLink        bool
		LiveUpdates           bool
		SpecialPage           bool
	}{
		Title:                 article.Title,
		SiteBrandTitle:        Config().SiteTitle,
		Anchorjs:              article.Anchorjs,
		Tocify:                article.Tocify,
		Timeline:              article.Timeline,
		NavTitleArticleSource: navTitleArticleSource,
		ArticleSourceCodeFS:   article.SrcFileName,
		ArticleSourceCodeURL:  filepath.ToSlash(article.SrcFileName),
		ArticleDstFileName:    article.DstFileName,
		ShowSourceLink:        article.ShowSourceLink,
		LiveUpdates:           article.LiveUpdates,
		SpecialPage:           article.SpecialPage,
	}
	err = t.Execute(buff, noItems)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buff.Bytes()
}

func GenerateNavTitleArticleSource(article db.Article, body string) string {
	t, err := template.New("navTitleArticleSource.tmpl").
		ParseFiles("templates/navTitleArticleSource.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var timeString string
	var time time.Time

	if article.ModificationDate != time {
		date := article.ModificationDate.Format("2 Jan 2006")
		timeString += `<div id="date"><p><span id="lastupdated">` + strings.ToLower(date) + `</span></p></div>`
	}

	if len(article.Tags) > 0 {
		timeString += `<div id="tags"><p>` + tagToLinkListInTimeline(&article) + `</p></div>`
	}

	noItems := struct {
		Title       string
		TitleNAV    string
		SeriesNAV   string
		TimeString  string
		Body        string
		SpecialPage bool
	}{
		Title:       article.Title,
		TitleNAV:    GenerateArticleNavigation(&article),
		SeriesNAV:   GenerateArticleSeriesNavigation(&article),
		TimeString:  timeString,
		Body:        body,
		SpecialPage: article.SpecialPage,
	}
	generatedHTMLbuff := bytes.NewBufferString("")
	err = t.Execute(generatedHTMLbuff, noItems)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return generatedHTMLbuff.String()
}

func Render(a db.Article) string {
	ac, err := db.Instance().GetCache(a)
	if err != nil || Config().Force == 1 {
		if Config().Verbose > 1 {
			fmt.Println(color.YellowString("pandoc run for article"), a.DstFileName)
		}
		generatedHTML, err := PandocMarkdown2HTML(a.ArticleMDWNSource)
		if err != nil {
			fmt.Println("An error occurred during pandoc pipeline run: ", err)
			panic(err)
		}
		errSet := db.Instance().SetCache(a, generatedHTML)
		if errSet != nil {
			fmt.Println(errSet)
		}
		return generatedHTML
	} else {
		fmt.Println(color.YellowString("cache hit, no pandoc run for article"), a.DstFileName)
		return ac
	}
}

func GenerateArticleNavigation(article *db.Article) string {
	if article.SpecialPage == true {
		return ""
	}
	titleNAV := ""
	p, err := db.Instance().PrevArticle(*article)
	if err == nil {
		titleNAV +=
			`<span id="articleNavLeft"> <a href="` + p.DstFileName + `"> 
      <span class="glyphiconLink glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article"> </span> prev. article
    </a> </span>`
	}
	n, err := db.Instance().NextArticle(*article)
	if err == nil {
		titleNAV +=
			`<span id="articleNavRight"><a href="` + n.DstFileName + `"> 
        next article <span class="glyphiconLink glyphicon glyphicon-chevron-right" aria-hidden="true" title="next article"></span>
    </a> </span>`
	}
	return titleNAV
}

func GenerateArticleSeriesNavigation(article *db.Article) string {
	if article.SpecialPage == true {
		return ""
	}
	seriesNAV := ""

	if article.Series != "" {
		seriesNAV =
			`
      <div id="seriesContainer">
      <a href="timeline.html?filter=series::` + article.Series + `" title="article series ` + article.Series + `" class="seriesbtn btn btn-primary">` +
				article.Series + `</a>
        <header class="seriesHeader">
          <div id="seriesLeft">`
		sp, sperr := db.Instance().PrevArticleInSeries(*article)
		if sperr == nil {
			seriesNAV += `<a href="` + sp.DstFileName + `">` +
				`<span class="glyphiconLinkSeries glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article in series"></span>
            </a> `
		}
		seriesNAV += `  </div>
          <div id="seriesRight">`
		sn, snerr := db.Instance().NextArticleInSeries(*article)
		if snerr == nil {
			seriesNAV += `   <a href="` + sn.DstFileName + `">
              <span class="glyphiconLinkSeries glyphicon glyphicon-chevron-right" aria-hidden="true" title="next article in series"></span>
            </a>`
		}
		seriesNAV += `</div>
        </header>
      </div>`
	}
	return seriesNAV
}
