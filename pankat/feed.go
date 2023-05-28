package pankat

import (
	"fmt"
	"github.com/fatih/color"
	htemplate "html/template"
	"os"
	"path/filepath"
	"time"
)

func RenderFeed(articles Articles) {
	defer timeElapsed("RenderFeed")()
	fmt.Println(color.YellowString("Rendering all xml feeds"))

	var history string
	var article Article

	article.Title = "feed"

	tagsMap := make(map[string]int)
	seriesMap := make(map[string]int)

	for _, a := range articles {
		if a.Series != "" {
			seriesMap[a.Series]++
		}
		for _, t := range a.Tags {
			tagsMap[t]++
		}
	}

	history += `
	<p>
		<div>want to follow the blog by feed, this one contains all articles:</div>
		<a href="feed/feed.xml">feed.xml</a>
    </p>
    `
	for k := range tagsMap {
		generateFeedXML(articles.FilterByTag(k), "tag_"+k)
	}

	for k := range seriesMap {
		generateFeedXML(articles.FilterBySeries(k), "series_"+k)
	}

	generateFeedXML(articles, "feed")

	history += `
    <h2>feed by tag/series</h2>
    `
	// sort the tags
	tagsSlice := rankByWordCount(tagsMap)

	history += `<p id="tagCloud">`
	for _, e := range tagsSlice {
		history += `<a class="tagbtn btn btn-primary" onClick="showXML('tag_` + e.Key + `')">` + e.Key + `</a>`
	}

	seriesSlice := rankByWordCount(seriesMap)
	//fmt.Println(seriesSlice)
	history += `<p id="seriesCloud">`
	for _, e := range seriesSlice {
		history += `<a class="seriesbtn btn btn-primary" onClick="showXML('series_` + e.Key + `')">` + e.Key + `</a>`
	}
	history += `</p>
	<p>
     Are you interested in a single tag or series? Then just make a selection above and copy the url from below afterwards:<br>

     <div class="feedURL" id="feedURL"> select tag or series</div>
     </p>

     <script>
     var showXML = function(name) {
        $('#feedURL')[0].innerHTML="> <a href=\"feed/" + name +".xml\">" + name +".xml</a>"
     }
     </script>
    `

	navTitleArticleSource := GenerateNavTitleArticleSource(articles, article, history)
	standalonePageContent := GenerateStandalonePage(articles, article, navTitleArticleSource)

	outD := Config().DocumentsPath + "/"
	err := os.MkdirAll(outD, 0755)
	if err != nil {
		panic(err)
	}
	outName := outD + "feed.html"
	err1 := os.WriteFile(outName, standalonePageContent, 0644)
	if err1 != nil {
		panic(err1)
	}
}

func generateFeedXML(articles Articles, fileName string) {
	if Config().Verbose > 0 {
		fmt.Println("Generating feed: " + fileName)
	}
	feedUrl := Config().SiteURL + "/feed/" + "feed.xml"
	z := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
  <id>` + Config().SiteURL + "/" + "index.html" + `</id>
  <link type="text/html" rel="alternate" href="` + feedUrl + `"/>
  <link type="application/atom+xml" rel="self" href="` + feedUrl + `"/>
  <title>` + Config().SiteTitle + `</title>
  <updated>` + time.Now().Format("2006-01-02T15:04:05-07:00") + `</updated>`

	for _, e := range articles {
		url := filepath.Join(Config().SiteURL, filepath.ToSlash(filepath.Dir(e.SrcFileName)), e.DstFileName)
		z += `
  <entry>
    <id>` + url + `</id>
    <link type="text/html" rel="alternate" href="` + url + `"/>
    <title>
        ` + e.Title + `
    </title>
    <updated>` + e.ModificationDate.Format("2006-01-02T15:04:05-07:00") + `</updated>`

		for _, t := range e.Tags {
			z += `<category scheme="` + Config().SiteURL + `" term="` + t + `"/>`
		}
		//BUG: feed needs ./posts/media/ URLs instead of ./media/ URLs
		z += `<author>
      <name>qknight</name>
      <uri>https://github.com/qknight</uri>
    </author>
    <content type="html">` + htemplate.HTMLEscaper(e.Render()) + `</content>
  </entry>`
	}

	z += `</feed>`
	errMkdir := os.MkdirAll(Config().DocumentsPath+"/feed", 0755)
	if errMkdir != nil {
		fmt.Println(errMkdir)
		panic(errMkdir)
	}
	feedName := filepath.Join(Config().DocumentsPath, "feed", fileName+".xml")
	err2 := os.WriteFile(feedName, []byte(z), 0644)
	if err2 != nil {
		fmt.Println(err2)
		panic(err2)
	}
}
