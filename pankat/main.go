package pankat

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

type Pankat struct {
}

func tagToLinkListInTimeline(a *Article) string {
	var tags []string
	tags = a.Tags
	var output string
	for _, e := range tags {
		output += `<a href="timeline.html?filter=tag::` + e + `" class="tagbtn btn btn-primary">` + e + `</a>`
	}
	return output
}

func tagToLinkList(a *Article) string {
	var tags []string
	tags = a.Tags

	var output string
	for _, e := range tags {
		output += `<a class="tagbtn btn btn-primary" onClick="setFilter('tag::` + e + `', 1)">` + e + `</a>`
	}
	return output
}

func seriesToLinkList(a *Article) string {
	var output string
	output += `<a class="seriesbtn btn btn-primary" onClick="setFilter('series::` + a.Series + `', 1)">` + a.Series + `</a>`
	return output
}

func GetTargets(path string) Articles {
	defer timeElapsed("GetTargets")()
	fmt.Println(color.YellowString("GetTargets: searching and parsing articles with *.mdwn"))
	targets := findArticlesOnDisk(path)
	return targets
}

// scan the directory for .mdwn files recursively
func findArticlesOnDisk(path string) Articles {
	var articles Articles
	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		buf := filepath.Join(path, entry.Name())
		//     fmt.Println("reading buf: ", buf)
		if entry.IsDir() {
			if entry.Name() == ".git" {
				continue
			}
			//       fmt.Println(buf)
			n := findArticlesOnDisk(buf)
			articles = append(articles, n...)
		} else {
			if strings.HasSuffix(entry.Name(), ".mdwn") {
				var a Article
				v := strings.TrimSuffix(entry.Name(), ".mdwn") // remove .mdwn
				a.Title = strings.Replace(v, "_", " ", -1)     // add whitespaces
				a.DstFileName = v + ".html"
				a.SrcFileName = filepath.Join(path, entry.Name())
				fh, errOpen := os.Open(a.SrcFileName)
				if errOpen != nil {
					fmt.Println(errOpen)
					continue
				}
				f := bufio.NewReader(fh)
				_article, errRead := io.ReadAll(f)
				if errRead != nil {
					fmt.Println(errRead)
					continue
				}
				_article = ProcessPlugins(_article, &a)
				a.ArticleMDWNSource = _article
				articles = append(articles, &a)
				errClose := fh.Close()
				if errClose != nil {
					fmt.Println(errClose)
				}
			}
		}
	}
	sort.Sort(articles)
	return articles
}

func rankByWordCount(wordFrequencies map[string]int) TagsSlice {
	pl := make(TagsSlice, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type TagsSlice []Pair

func (p TagsSlice) Len() int           { return len(p) }
func (p TagsSlice) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p TagsSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var sendLiveUpdateViaWS func(string, string) = emtpyFunc

func emtpyFunc(string, string) {

}

func OnArticleChange(f func(string, string)) {
	sendLiveUpdateViaWS = f
}

func RenderPosts(articlesAll Articles) {
	defer timeElapsed("RenderPosts")()
	fmt.Println(color.YellowString("Rendering posts"))
	for _, article := range articlesAll {
		RenderPost(articlesAll, article)
	}
}

func RenderPost(articles Articles, article *Article) {
	if Config().Verbose > 0 {
		fmt.Println("Rendering article '" + article.Title + "'")
	}
	navTitleArticleHTML := GenerateNavTitleArticleSource(articles, *article, article.Render())
	standalonePageContent := GenerateStandalonePage(articles, *article, navTitleArticleHTML)
	sendLiveUpdateViaWS(filepath.ToSlash(article.SrcFileName), navTitleArticleHTML)

	outD := Config().DocumentsPath
	errMkdir := os.MkdirAll(outD, 0755)
	if errMkdir != nil {
		fmt.Println(errMkdir)
		panic(errMkdir)
	}
	// write to disk
	outName := filepath.Join(outD, article.DstFileName)
	err5 := os.WriteFile(outName, standalonePageContent, 0644)
	if err5 != nil {
		fmt.Println(err5)
		panic(article)
	}
}

func GenerateStandalonePage(articles Articles, article Article, navTitleArticleSource string) []byte {
	buff := bytes.NewBufferString("")
	t, err := template.New("standalonePage.tmpl").
		ParseFiles("templates/standalonePage.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	noItems := struct {
		Title                 string
		SiteURL               string
		SiteBrandTitle        string
		Anchorjs              bool
		Tocify                bool
		Timeline              bool
		NavTitleArticleSource string
		ArticleSourceCodeURL  string // URL of the mdwn source code seen from the web
		ArticleSourceCodeFS   string // FS of the mdwn on disk (win/linux/...)
		SourceReference       bool
		WebsocketSupport      bool
		SpecialPage           bool
	}{
		Title:                 article.Title,
		SiteURL:               Config().SiteURL,
		SiteBrandTitle:        Config().SiteTitle,
		Anchorjs:              article.Anchorjs,
		Tocify:                article.Tocify,
		Timeline:              article.Timeline,
		NavTitleArticleSource: navTitleArticleSource,
		ArticleSourceCodeFS:   article.SrcFileName,
		ArticleSourceCodeURL:  filepath.ToSlash(article.SrcFileName),
		SourceReference:       article.SourceReference,
		WebsocketSupport:      article.WebsocketSupport,
		SpecialPage:           article.SpecialPage,
	}
	err = t.Execute(buff, noItems)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buff.Bytes()
}

func GenerateNavTitleArticleSource(articles Articles, article Article, body string) string {
	t, err := template.New("navTitleArticleSource.tmpl").
		ParseFiles("templates/navTitleArticleSource.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var meta string
	var timeT time.Time

	if article.ModificationDate != timeT {
		meta += `<div id="date"><p><span id="lastupdated">` + article.ModificationDate.Format("2 Jan 2006") + `</span></p></div>`
	}

	if len(article.Tags) > 0 {
		meta += `<div id="tags"><p>` + tagToLinkListInTimeline(&article) + `</p></div>`
	}

	noItems := struct {
		Title       string
		TitleNAV    string
		SeriesNAV   string
		Meta        string
		Body        string
		SpecialPage bool
	}{
		Title:       article.Title,
		TitleNAV:    articles.GetTitleNAV(&article),
		SeriesNAV:   articles.GetSeriesNAV(&article),
		Meta:        meta,
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

func Init() {
	pflag.String("documents", "myblog/", "input directory ('documents') in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files")
	pflag.String("siteURL", "https://lastlog.de/blog", "The URL of the blog, for example: 'https://example.com/blog'")
	pflag.String("siteTitle", "lastlog.de/blog", "Title which is inserted top left, for example: 'lastlog.de/blog'")
	pflag.Int("verbose", 0, "verbosity level")
	pflag.Int("force", 0, "forced complete rebuild, not using cache")
	pflag.String("ListenAndServe", ":8000", "ip:port where pankat-server listens, for example: 'localhost:8000'")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	input := viper.GetString("documents")
	siteURL_ := viper.GetString("siteURL")
	siteTitle_ := viper.GetString("siteTitle")

	i1, err := filepath.Abs(input)
	documentsPath := i1

	myMd5HashMapJson_ := filepath.Join(documentsPath, ".ArticlesCache.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err1 := os.Chdir(documentsPath)
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	Config().SiteURL = siteURL_
	Config().DocumentsPath = documentsPath
	Config().SiteURL = siteURL_
	Config().SiteTitle = siteTitle_
	Config().MyMd5HashMapJson = myMd5HashMapJson_
	Config().Verbose = viper.GetInt("verbose")
	Config().Force = viper.GetInt("force")
	Config().ListenAndServe = viper.GetString("ListenAndServe")

}

func SetMostRecentArticle(articlesPosts Articles) {
	mostRecentArticle := ""
	if len(articlesPosts) > 0 {
		mostRecentArticle = filepath.ToSlash(articlesPosts[0].DstFileName)
	} else {
		mostRecentArticle = "timeline.html"
	}
	indexContent :=
		`
<html>
<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
<meta http-equiv="Pragma" content="no-cache" />
<meta http-equiv="Expires" content="0" />
<meta http-equiv="refresh" content="0; url=` + mostRecentArticle + `" />
</html>
`
	outIndexName := filepath.Join(Config().DocumentsPath, "index.html")
	errn := os.WriteFile(outIndexName, []byte(indexContent), 0644)
	if errn != nil {
		panic(errn)
	}
}

func UpdateBlog() {
	defer timeElapsed("UpdateBlog")()
	fmt.Println(color.GreenString("pankat-static"), "starting!")
	fmt.Println(color.YellowString("Documents path: "), Config().DocumentsPath)

	articles := GetTargets(".").FilterOutDrafts()
	fmt.Println(color.YellowString("GetTargets: found"), articles.Targets().Len(), color.YellowString("articles"))

	RenderPosts(articles)

	articles = articles.FilterOutSpecialPages().Targets()
	RenderTimeline(articles)
	RenderFeed(articles)
	SetMostRecentArticle(articles)
}
