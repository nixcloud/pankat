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
	"sort"
	"strings"
)

// scan the directory for .mdwn files recursively
func findArticlesOnDisk(path string) {
	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		newDir := filepath.Join(path, entry.Name())
		//     fmt.Println("reading newDir: ", newDir)
		if entry.IsDir() {
			if entry.Name() == ".git" {
				continue
			}
			//       fmt.Println(newDir)
			findArticlesOnDisk(newDir)
		} else {
			if strings.HasSuffix(entry.Name(), ".mdwn") {
				v := strings.TrimSuffix(entry.Name(), ".mdwn") // remove .mdwn
				DstFileName := v + ".html"
				p, _ := filepath.Rel(Config().DocumentsPath, filepath.Join(path, entry.Name()))
				SrcFileName := p
				newArticle, _ := ReadRAWMDWNAndProcessPlugins(SrcFileName, DstFileName)
				db.Instance().Add(newArticle)
			}
		}
	}
}

func ReadRAWMDWNAndProcessPlugins(SrcFileName string, DstFileName string) (*db.Article, error) {
	fh, errOpen := os.Open(SrcFileName)
	if errOpen != nil {
		fmt.Println(errOpen)
		return nil, errors.New("ReadRAWMDWNAndProcessPlugins: " + errOpen.Error())
	}
	f := bufio.NewReader(fh)
	rawMDWNSourceArticle, errRead := io.ReadAll(f)
	if errRead != nil {
		fmt.Println(errRead)
		return nil, errors.New("ReadRAWMDWNAndProcessPlugins: " + errRead.Error())
	}
	errClose := fh.Close()
	if errClose != nil {
		fmt.Println(errClose)
		return nil, errors.New("ReadRAWMDWNAndProcessPlugins: " + errClose.Error())
	}
	var newArticle db.Article
	newArticle.SrcFileName = SrcFileName
	newArticle.DstFileName = DstFileName
	filename := filepath.Base(SrcFileName)
	v := strings.TrimSuffix(filename, ".mdwn")          // remove .mdwn
	newArticle.Title = strings.Replace(v, "_", " ", -1) // add whitespaces
	newArticle.ArticleMDWNSource = ProcessPlugins(rawMDWNSourceArticle, &newArticle)
	newArticle.LiveUpdates = true
	newArticle.ShowSourceLink = true
	//fmt.Println("Article after ProcessPlugins. article.LiveUpdates:", newArticle.LiveUpdates)
	//fmt.Println("Article after ProcessPlugins. article.Draft:", newArticle.Draft)
	//fmt.Println("Article after ProcessPlugins. article.SpecialPage:", newArticle.SpecialPage)
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

func emtpyFunc(string, string) {

}

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
