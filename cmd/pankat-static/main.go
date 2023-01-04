package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"pankat"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	pflag.String("input", "documents", "input directory  ('documents'') in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files")
	pflag.String("output", "output", "output directory ('output') all generated files will be stored there and all directories like css/ js/ images and fonts/ will be rsynced there")
	pflag.String("siteURL", "https://lastlog.de/blog", "The URL of the blog, for example: 'https://example.com/blog'")
	pflag.String("siteTitle", "lastlog.de/blog", "Title which is inserted top left, for example: 'lastlog.de/blog'")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	input := viper.GetString("input")
	output := viper.GetString("output")
	siteURL_ := viper.GetString("siteURL")
	siteTitle_ := viper.GetString("siteTitle")

	i1, err := filepath.Abs(input)
	inputPath_ := i1

	myMd5HashMapJson_ := filepath.Clean(inputPath_ + "/.ArticlesCache.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	o1, err := filepath.Abs(output)
	outputPath_ := o1

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err1 := os.Chdir(inputPath_)
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	pankat.GetConfig().SiteURL = siteURL_
	pankat.GetConfig().InputPath = inputPath_
	pankat.GetConfig().OutputPath = outputPath_
	pankat.GetConfig().SiteURL = siteURL_
	pankat.GetConfig().SiteTitle = siteTitle_
	pankat.GetConfig().MyMd5HashMapJson = myMd5HashMapJson_

	fmt.Println(color.GreenString("pankat"), "starting!")
	fmt.Println("input Directory: ", pankat.GetConfig().InputPath)
	fmt.Println("output Directory: ", pankat.GetConfig().OutputPath)

	// find all .mdwn files
	f := make([]string, 0)
	f = append(f, "")
	articlesAll := pankat.GetTargets(".", f) // FIXME should that not be in Pankat.InputPath?

	articlesTopLevel := articlesAll.TopLevel().FilterOutDrafts()
	articlesPosts := articlesAll.Posts().FilterOutDrafts()

	// sort them by date
	sort.Sort(articlesPosts)

	// override default values for posts
	for _, e := range articlesPosts {
		e.Anchorjs = true
		e.Tocify = true
	}

	for _, e := range articlesPosts {
		_ = e
		// 		    fmt.Println(e.Title)
		//     z:=e.SrcDirectoryName + "/" + e.SrcFileName
		//     fmt.Println("   ", z)
		//     fmt.Println("   ",e.SrcDirectoryName)
		//     fmt.Println("   ",e.SrcFileName)
		//     fmt.Println("   ",e.ModificationDate)
		//     fmt.Println("   ",e.Tags)
		//     fmt.Println("   ",e.Draft)
		// 		    fmt.Println("   ",e.Series)
	}

	// generate about.html and additional pages
	pankat.RenderPosts(articlesTopLevel)

	// generate articles
	pankat.RenderPosts(articlesPosts)

	// generate posts.html (timeline)
	pankat.RenderTimeline(articlesPosts)

	// render feed.html
	pankat.RenderFeed(articlesPosts)

	s, _ := os.Getwd()
	//fmt.Println("Getwd path is: ", s)

	v, _ := filepath.Rel(s, pankat.GetConfig().OutputPath)
	v = filepath.ToSlash(v)
	fmt.Println("relative path is: ", v) // FIXME windows specific hack

	// copy static files as fonts/css/js to the output folder
	commands := []string{
		"rsync -av --relative js " + v,
		"rsync -av --relative css " + v,
		"rsync -av --relative fonts " + v,
		"rsync -av --relative images " + v,
		"rsync -av --relative posts/media " + v,
	}

	for _, el := range commands {
		parts := strings.Fields(el)
		head := parts[0]
		parts = parts[1:]
		fmt.Println("executing: ", el)
		out, err := exec.Command(head, parts...).Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		_ = out
	}
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
	outIndexName := filepath.Clean(pankat.GetConfig().OutputPath + "/" + "index.html")
	errn := ioutil.WriteFile(outIndexName, []byte(indexContent), 0644)
	if errn != nil {
		panic(errn)
	}
}
