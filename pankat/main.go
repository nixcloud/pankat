package pankat

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
	"pankat/db"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// scan the directory for .mdwn files recursively
func findArticlesOnDisk(path string) {
	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		newDir := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if entry.Name() == ".git" {
				continue
			}
			findArticlesOnDisk(newDir)
		} else {
			if strings.HasSuffix(entry.Name(), ".mdwn") {
				SrcFileName, _ := filepath.Rel(Config().DocumentsPath, filepath.Join(path, entry.Name()))
				newArticle, err := CreateArticleFromFilesystemMarkdown(SrcFileName)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// ignoring related articles updates (they all get rendered anyway)
				db.Instance().Set(newArticle)
			}
		}
	}
}

func CreateArticleFromFilesystemMarkdown(SrcFileName string) (*db.Article, error) {
	fh, errOpen := os.Open(SrcFileName)
	if errOpen != nil {
		fmt.Println(errOpen)
		return nil, errors.New("CreateArticleFromFilesystemMarkdown: " + errOpen.Error())
	}
	f := bufio.NewReader(fh)
	rawMDWNSourceArticle, errRead := io.ReadAll(f)
	if errRead != nil {
		fmt.Println(errRead)
		return nil, errors.New("CreateArticleFromFilesystemMarkdown: " + errRead.Error())
	}
	errClose := fh.Close()
	if errClose != nil {
		fmt.Println(errClose)
		return nil, errors.New("CreateArticleFromFilesystemMarkdown: " + errClose.Error())
	}
	var newArticle db.Article
	newArticle.SrcFileName = SrcFileName

	filenameMDWN := filepath.Base(SrcFileName)
	filename := strings.TrimSuffix(filenameMDWN, ".mdwn")
	newArticle.DstFileName = filename + ".html"

	newArticle.Title = strings.Replace(filename, "_", " ", -1) // add whitespaces
	newArticle.ArticleMDWNSource = ProcessPlugins(rawMDWNSourceArticle, &newArticle)
	newArticle.LiveUpdates = true
	newArticle.ShowSourceLink = true
	return &newArticle, nil
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

func emtpyFunc(string, string) {}

func OnArticleChange(f func(string, string)) {
	sendLiveUpdateViaWS = f
}

func Init() {
	pflag.String("documents", "myblog/", "input directory ('documents') in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files")
	pflag.String("siteTitle", "lastlog.de/blog", "Title which is inserted top left, for example: 'lastlog.de/blog'")
	pflag.Int("verbose", 0, "verbosity level")
	pflag.Int("force", 0, "forced complete rebuild, not using cache")
	pflag.String("ListenAndServe", ":8000", "ip:port where pankat-server listens, for example: 'localhost:8000'")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	input := viper.GetString("documents")
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

	Config().DocumentsPath = documentsPath
	Config().SiteTitle = siteTitle_
	Config().MyMd5HashMapJson = myMd5HashMapJson_
	Config().Verbose = viper.GetInt("verbose")
	Config().Force = viper.GetInt("force")
	Config().ListenAndServe = viper.GetString("ListenAndServe")

	mdwnSource := "# hello"
	text, err := PandocMarkdown2HTML([]byte(mdwnSource))
	if err != nil {
		fmt.Println("An error occurred during pandoc pipeline run: ", err)
		panic(err)
	}
	//text will be: <h1 id="hello">hello</h1>
	match, _ := regexp.MatchString(".*h1.*hello.*h1.*", text)
	if match != true {
		fmt.Println("An error occurred during pandoc pipeline run result match: ", err)
		panic(err)
	}
}

func SetMostRecentArticle() {
	mostRecentArticle := ""
	article, err := db.Instance().MostRecentArticle()
	if err != nil {
		mostRecentArticle = "timeline.html"
	} else {
		mostRecentArticle = filepath.ToSlash(article.DstFileName)
	}
	indexContent :=
		`<!DOCTYPE html>
<html  xmlns="http://www.w3.org/1999/xhtml">
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

	dbArticles, errQuerryAll := db.Instance().QueryAll()
	if errQuerryAll == nil {
		for _, article := range dbArticles {
			_, err := os.Stat(article.SrcFileName)
			if err != nil {
				fmt.Println(color.YellowString("Article: "), article.SrcFileName, color.RedString("does not exist anymore, deleting from database"))
				db.Instance().Del(article.SrcFileName)
			}
		}
	}

	findArticlesOnDisk(Config().DocumentsPath)

	articles, err := db.Instance().Articles()
	if err != nil {
		fmt.Errorf("Error: %s", err)
	} else {
		fmt.Println(color.YellowString("Articles: "), len(articles), color.YellowString("articles"))
		RenderPosts(articles)
		RenderTimeline()
	}

	specialPages, err := db.Instance().SpecialPages()
	if err != nil {
		fmt.Errorf("Error: %s", err)
	} else {
		fmt.Println(color.YellowString("SpecialPages: "), len(specialPages), color.YellowString("SpecialPages"))
		RenderPosts(specialPages)
	}
	SetMostRecentArticle()
}
