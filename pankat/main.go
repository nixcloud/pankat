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
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
)

type Pankat struct {
}

func tagToLinkList(a *Article) string {
	var tags []string
	tags = a.Tags

	var output string
	for _, e := range tags {
		//     fmt.Println("----------------")
		//     fmt.Println(outputPath)
		//     fmt.Println(a.SrcDirectoryName)

		// HACK should be moved to pankat-core
		relativeSrcRootPath, _ := filepath.Rel(a.SrcDirectoryName, "")
		relativeSrcRootPath = filepath.Clean(relativeSrcRootPath)

		output += `<a href="` + relativeSrcRootPath + `/posts.html?filter=tag::` + e + `" class="tagbtn btn btn-primary">` + e + `</a>`
	}
	return output
}

func GetTargets(path string, ret []string) Articles {
	defer timeElapsed("GetTargets")()

	fmt.Println(color.YellowString("GetTargets: searching and parsing articles with *.mdwn"))
	return getTargets_(path, ret)
}

// scan the direcotry for .mdwn files recurively
func getTargets_(path string, ret []string) Articles {
	var articles Articles
	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		buf := path + "/" + entry.Name()
		//     fmt.Println("reading buf: ", buf)
		if entry.IsDir() {
			if entry.Name() == ".git" {
				continue
			}
			ret = append(ret, buf)
			//       fmt.Println(buf)
			n := getTargets_(buf, ret)
			articles = append(articles, n...)
		} else {
			if strings.HasSuffix(entry.Name(), ".mdwn") {
				var a Article
				v := strings.TrimSuffix(entry.Name(), ".mdwn") // remove .mdwn

				a.Title = strings.Replace(v, "_", " ", -1) // add whitespaces
				a.DstFileName = v + ".html"
				a.BaseFileName = v
				a.SrcFileName = entry.Name()
				a.SrcDirectoryName = path
				fh, err := os.Open(path + "/" + entry.Name())

				f := bufio.NewReader(fh)

				if err != nil {
					fmt.Println(err)
					panic(err)
				}
				_article, err := io.ReadAll(f)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}

				_article = ProcessPlugins(_article, &a)

				a.Article = _article
				articles = append(articles, &a)
				err = fh.Close()
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
			}
		}
	}

	return articles
}

func ProcessPlugins(_article []byte, article *Article) []byte {
	var _articlePostprocessed []byte

	re := regexp.MustCompile("\\[\\[!(.*?)\\]\\]")
	z := re.FindAllIndex(_article, -1)

	prevPos := 0
	var foundPlugins []string
	for i := 0; i <= len(z); i++ {
		if i == len(z) {
			_articlePostprocessed = append(_articlePostprocessed, _article[prevPos:]...)
			break
		}
		n := z[i]

		// include normal content (not plugin processed)
		if prevPos != n[0] {
			_articlePostprocessed = append(_articlePostprocessed, _article[prevPos:n[0]]...)
		}

		// include plugin processed stuff
		t, name := callPlugin(_article[n[0]:n[1]], article)
		foundPlugins = append(foundPlugins, name)
		_articlePostprocessed = append(_articlePostprocessed, t...)
		prevPos = n[1]
	}
	if GetConfig().Verbose > 1 {
		fmt.Println(article.DstFileName, color.GreenString("plugins:"), foundPlugins)
	}
	return _articlePostprocessed
}

func callPlugin(in []byte, article *Article) ([]byte, string) {
	a := len(in) - 2
	p := string(in[3:a])
	//   fmt.Println(p)
	var output []byte

	f := strings.Fields(p)
	var name string
	if len(f) > 0 {
		name = f[0]
	} else {
		var z []byte
		return z, ""
	}

	//   fmt.Println("\n=========== ", name, " ===========")

	switch name {
	case "meta":
		re := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+")
		z := re.FindIndex(in)
		var t time.Time
		if z != nil {
			s := string(in[z[0]:z[1]])
			//           fmt.Println(s)
			const longForm = "2006-01-02 15:04"
			t, _ = time.Parse(longForm, s)
			article.ModificationDate = t
			//           fmt.Println(t)
		}
		// 	case "warning":
		//     if len(f) > 1 {
		//       o := `<div id="bar">` + strings.Join(f[1:len(f)], " ") + `</div>`
		//       output = []byte(o)
		//     }
	case "series":
		if len(f) > 1 {
			article.Series = strings.Join(f[1:], " ")
		}
	case "tag":
		if len(f) > 1 {
			article.Tags = f[1:]
		}
	case "draft":
		article.Draft = true

	case "img":
		b := strings.Join(f[1:], " ")
		//      fmt.Println("\n------------\n", article.SrcDirectoryName)
		//      fmt.Println(f[1])

		//HACK should be moved to pankat-core
		relativeSrcRootPath, _ := filepath.Rel(article.SrcDirectoryName, "./posts")
		relativeSrcRootPath = filepath.Clean(relativeSrcRootPath)
		//      fmt.Println(relativeSrcRootPath)

		o := `<a href="` + relativeSrcRootPath + "/" + f[1] + `"><img src=` + relativeSrcRootPath + "/" + b + `></a>`
		output = []byte(o)

	case "summary":
		article.Summary = strings.Join(f[1:], " ")
	default:
		fmt.Println(article.SrcFileName + ": plugin '" + name + "'" + color.RedString(" NOT supported"))
	}
	return output, name
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
	var RenderRuns []Articles
	RenderRuns = append(RenderRuns, articlesAll.TopLevel())
	RenderRuns = append(RenderRuns, articlesAll.Posts())
	for _, articles := range RenderRuns {
		for _, article := range articles {
			RenderPost(articles, article)
		}
	}
}

func RenderPost(articles Articles, article *Article) {
	if GetConfig().Verbose > 0 {
		fmt.Println("Rendering article '" + article.Title + "'")
	}
	navTitleArticleSource := GenerateNavTitleArticleSource(articles, *article, article.Render())
	standalonePageContent := generateStandalonePage(articles, *article, navTitleArticleSource)
	outD := filepath.Clean(GetConfig().DocumentsPath + "/" + article.SrcDirectoryName + "/")
	sendLiveUpdateViaWS(article.SrcFileName, navTitleArticleSource)

	errMkdir := os.MkdirAll(outD, 0755)
	if errMkdir != nil {
		fmt.Println(errMkdir)
		panic(errMkdir)
	}
	// write to disk
	outName := filepath.Clean(outD + "/" + article.DstFileName)
	err5 := os.WriteFile(outName, standalonePageContent, 0644)
	if err5 != nil {
		fmt.Println(err5)
		panic(article)
	}
}

func generateStandalonePage(articles Articles, article Article, navTitleArticleSource string) []byte {
	buff := bytes.NewBufferString("")
	t, err := template.New("standalonePage.tmpl").
		ParseFiles("templates/standalonePage.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	relativeSrcRootPath, _ := filepath.Rel(article.SrcDirectoryName, "")
	relativeSrcRootPath = filepath.Clean(relativeSrcRootPath)

	noItems := struct {
		Title                 string
		RelativeSrcRootPath   string
		SiteURL               string
		SiteBrandTitle        string
		Anchorjs              bool
		Tocify                bool
		Timeline              bool
		NavTitleArticleSource string
		SrcDirectoryName      string
		ArticleSourceCodeURL  string
		SourceReference       bool
		WebsocketSupport      bool
		NavAndSeriesElements  bool
	}{
		Title:                 article.Title,
		RelativeSrcRootPath:   relativeSrcRootPath,
		SiteURL:               GetConfig().SiteURL,
		SiteBrandTitle:        GetConfig().SiteTitle,
		Anchorjs:              article.Anchorjs,
		Tocify:                article.Tocify,
		Timeline:              article.Timeline,
		NavTitleArticleSource: navTitleArticleSource,
		SrcDirectoryName:      article.SrcDirectoryName,
		ArticleSourceCodeURL:  filepath.Clean(article.SrcDirectoryName + "/" + article.SrcFileName),
		SourceReference:       article.SourceReference,
		WebsocketSupport:      article.WebsocketSupport,
		NavAndSeriesElements:  article.NavAndSeriesElements,
	}
	err = t.Execute(buff, noItems)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buff.Bytes()
}

func GenerateNavTitleArticleSource(articles Articles, article Article, body string) string {
	buff := bytes.NewBufferString("")
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
		meta += `<div id="tags"><p>` + tagToLinkList(&article) + `</p></div>`
	}

	noItems := struct {
		Title                string
		TitleNAV             string
		SeriesNAV            string
		Meta                 string
		Body                 string
		NavAndSeriesElements bool
	}{
		Title:                article.Title,
		TitleNAV:             articles.GetTitleNAV(&article),
		SeriesNAV:            articles.GetSeriesNAV(&article),
		Meta:                 meta,
		Body:                 body,
		NavAndSeriesElements: article.NavAndSeriesElements,
	}
	err = t.Execute(buff, noItems)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buff.String()
}

func Init() {
	pflag.String("documents", "myblog/", "input directory ('documents'') in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files")
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

	myMd5HashMapJson_ := filepath.Clean(documentsPath + "/.ArticlesCache.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err1 := os.Chdir(documentsPath)
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	GetConfig().SiteURL = siteURL_
	GetConfig().DocumentsPath = documentsPath
	GetConfig().SiteURL = siteURL_
	GetConfig().SiteTitle = siteTitle_
	GetConfig().MyMd5HashMapJson = myMd5HashMapJson_
	GetConfig().Verbose = viper.GetInt("verbose")
	GetConfig().Force = viper.GetInt("force")
	GetConfig().ListenAndServe = viper.GetString("ListenAndServe")
}

func SetMostRecentArticle(articlesPosts Articles) {
	mostRecentArticle := ""
	if len(articlesPosts) > 0 {
		mostRecentArticle = filepath.Clean(articlesPosts[0].SrcDirectoryName + "/" + articlesPosts[0].DstFileName)
	} else {
		mostRecentArticle = "posts.html"
	}
	indexContent :=
		`
<html>
<meta http-equiv="refresh" content="0; url=` + mostRecentArticle + `" />
</html>
`
	outIndexName := filepath.Clean(GetConfig().DocumentsPath + "/" + "index.html")
	errn := os.WriteFile(outIndexName, []byte(indexContent), 0644)
	if errn != nil {
		panic(errn)
	}
}

func GetArticles() Articles {

	// find all .mdwn files
	f := make([]string, 0)
	f = append(f, "")

	articles := GetTargets(".", f).FilterOutDrafts()

	// sort them by date
	sort.Sort(articles)

	fmt.Println(color.YellowString("GetTargets: found"), articles.Posts().Len(), color.YellowString("articles"))

	return articles
}

func UpdateBlog() {
	defer timeElapsed("UpdateBlog")()

	fmt.Println(color.GreenString("pankat-static"), "starting!")
	fmt.Println(color.YellowString("Documents path: "), GetConfig().DocumentsPath)

	articles := GetArticles()
	RenderPosts(articles)
	RenderTimeline(articles.Posts())
	RenderFeed(articles.Posts())
	SetMostRecentArticle(articles.Posts())
}
