package pankat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	htemplate "html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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

func RenderTimeline(articles Articles) {
	defer timeElapsed("RenderTimeline")()
	fmt.Println(color.YellowString("Rendering timeline in posts.html"))

	var pageContent string
	var article Article

	article.Title = "all posts"
	article.Timeline = true

	article.SrcDirectoryName = ""

	t, err := json.Marshal(articles.CreateJSMetadata())
	if err != nil {
		fmt.Println("json.Marshal error:", err)
	}
	pageContent += `<script type="application/json" id="MetaData">` + string(t) + `</script>`

	tagsMap := make(map[string]int)
	seriesMap := make(map[string]int)

	for _, article := range articles {
		if article.Series != "" {
			seriesMap[article.Series]++
		}
		for _, t := range article.Tags {
			tagsMap[t]++
		}
	}

	// sort the tags
	tagsSlice := rankByWordCount(tagsMap)
	if GetConfig().Verbose > 0 {
		fmt.Println(color.GreenString("tagsSlice"), tagsSlice)
	}

	pageContent += `<p id="tagCloud">`
	for _, e := range tagsSlice {
		pageContent += `<a class="tagbtn btn btn-primary" onClick="setFilter('tag::` + e.Key + `', 1)">` + e.Key + `</a>`
	}
	pageContent += `</p>`

	seriesSlice := rankByWordCount(seriesMap)
	if GetConfig().Verbose > 0 {
		fmt.Println(color.GreenString("seriesSlice"), seriesSlice)
	}

	pageContent += `<p id="seriesCloud">`
	for _, e := range seriesSlice {
		pageContent += `<a class="seriesbtn btn btn-primary" onClick="setFilter('series::` + e.Key + `', 1)">` + e.Key + `</a>`
	}
	pageContent += `</p>`

	pageContent += `

    <a class="btn btn-primary" onClick="setFilter('', 1)">show all (clear filters)</a>

    <p class="lead">filter the posts (click tag/series) above:</p>
    
    <div id="timeline" class="timeline-container">
    <br class="clear">
`

	var year string

	for i, article := range articles {
		if i == 0 {
			v := article.ModificationDate.Add(1000 * 1000 * 1000 * 60 * 60 * 24 * 365) // add one year
			year = v.Format("2006")
			pageContent += `
	          <div class="timeline-wrapper pankat_year pankat_year_` + year + `">
	            <dl class="timeline-series">
                 <h2 class="timeline-time"><span>` + year + `</span></h2>`
			year = article.ModificationDate.Format("2006")
		}

		//     fmt.Println("----")
		//     fmt.Println("  ", e.Title)
		//     fmt.Println("  ", e.SrcDirectoryName)
		//     fmt.Println("  ", inputPath)

		if year != article.ModificationDate.Format("2006") {
			pageContent += `
         </dl><!-- /.timeline-series -->
       </div><!-- /.timeline-wrapper -->
       <div class="timeline-wrapper pankat_year pankat_year_` + year + `">
    	<h2 class="timeline-time"><span>` + year + `</span></h2>
			<dl class="timeline-series">`
			year = article.ModificationDate.Format("2006")
		}

		// a hacky but straight-forward way to make tagToLinkList(...) work by
		// fooling a different base article
		var v Article
		v = *article
		v.SrcDirectoryName = ""

		//     <h3>` + e.ModificationDate.Format("2 Jan 2006") + `</h3>
		//     <span class="glyphicon glyphicon-chevron-link" aria-hidden="true" title="article"></span>
		pageContent += `
          <dt class="timeline-event posting_` + strconv.Itoa(i) + `">` + article.Title + `</dt>
          <dd class="timeline-event-content posting_` + strconv.Itoa(i) + `">
            <div class="postingsEntry">
              <p class="summary">` + article.Summary + ` <a href="` + filepath.Clean(article.SrcDirectoryName+"/"+article.DstFileName) + `">open complete article</a></p>
              <p class="tag">` + tagToLinkList(&v) + `</p>
            </div>
            <br class="clear">
          </dd><!-- /.timeline-event-content -->`
		if i == len(articles)-1 {

			v := article.ModificationDate.Add(-1000 * 1000 * 1000 * 60 * 60 * 24 * 365) // add one year
			year = v.Format("2006")

			pageContent += `
            <div class="timeline-wrapper pankat_year pankat_year_` + year + `">
			<dl class="timeline-series">

    	    <h2 class="timeline-time"><span>` + year + `</span></h2>
			</dl><!-- /.timeline-series -->
			</div><!-- /.timeline-wrapper -->
			`

			pageContent += `
		    </dl><!-- /.timeline-series -->
		  </div><!-- /.timeline-wrapper -->
			`
		}
	}

	pageContent += `
	<script>
      var MetaData
      function getURLParameter(name) {
        return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)').exec(location.search)||[,""])[1].replace(/\+/g, '%20'))||null
      }

      var setFilter = function(filter, addHistory) {
        var selection = []
        try {
        var type = filter.split("::")[0]
        var identifier = filter.split("::")[1]
        } catch(e) {
          console.log("removing filter selection because of split error handling")
          for (i=0; i < MetaData.ArticleCount; i++) {
            var n = ".posting_" + i;
            $(n).css('display', 'block');
          }
          if (addHistory === 1)
		    window.pageContent.pushState('', '',  window.location.pathname);
          return
        }
        // hide all years
        $(".pankat_year").css('display', 'none');
        if (type == "tag" && typeof(MetaData.Tags[identifier]) !== "undefined") {
          selection = MetaData.Tags[identifier]
        } else if (type == "series" && typeof(MetaData.Series[identifier]) !== "undefined") {
          selection = MetaData.Series[identifier]
        } else {
          console.log("removing filter selection")
          for (i=0; i < MetaData.ArticleCount; i++) {
            var n = ".posting_" + i;
            $(n).css('display', 'block');
          }
          $(".pankat_year").css('display', 'block');
    	  if (addHistory === 1)
		    window.pageContent.pushState('', '',  window.location.pathname);
          return
        }

        //console.log(selection, type, identifier)

        for (i=0; i < MetaData.ArticleCount; i++) {
          var n = ".posting_" + i;
          if (selection.includes(i)) {
            //console.log("show", n)
            $(n).css('display', 'block');
			// show respective year
			Object.keys(MetaData.Years).forEach(function(key, index) {
			  MetaData.Years[key].forEach(function(article) {
                if (article === i) {
                  year = parseInt(key);
                  var n = ".pankat_year_" + year;
			      console.log("Showing year with jquery class: ",  n)

			      $(n).css('display', 'block');
                  var n = ".pankat_year_" + (year + 1);
			      $(n).css('display', 'block');
			      console.log("Showing year with jquery class: ",  n)

                }
			  })
			});
          } else {
            //console.log("hide", n)
            $(n).css('display', 'none');
          }
        }

        if (addHistory === 1)
          window.pageContent.pushState('', '',  window.location.pathname + '?filter=' + filter);
      }
 
      $(document).ready(function() {
        MetaData = JSON.parse(document.getElementById('MetaData', 0).innerHTML)
        var filter = getURLParameter("filter");
        //console.log("document.ready(), filter: " + filter)
		$.timeliner({
		  oneOpen: false,
		  startState: 'open'
		});
        setFilter(filter, 0)
      });

      // browser pageContent button was used, so we need to update the posts, but not the browser pageContent
      window.addEventListener("popstate", function() {
        var filter = getURLParameter("filter");
        setFilter(filter, 0);
      });
      </script>
      </div>
`
	posts := generateStandalonePage(articles, article, pageContent)

	outD := GetConfig().OutputPath + "/"
	err = os.MkdirAll(outD, 0755)
	if err != nil {
		panic(err)
	}
	outName := outD + "posts.html"
	err1 := os.WriteFile(outName, posts, 0644)
	if err1 != nil {
		panic(err1)
	}
}

func RenderFeed(articles Articles) {
	defer timeElapsed("RenderFeed")()
	fmt.Println(color.YellowString("Rendering all xml feeds"))

	var history string
	var article Article

	article.Title = "feed"

	article.SrcDirectoryName = ""

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
		<a href="feed.xml">feed.xml</a>
    </p>
    `
	for k := range tagsMap {
		generateFeedXML(articles.FilterByTag(k), "tag_"+k)
	}

	for k := range seriesMap {
		generateFeedXML(articles.FilterBySeries(k), "series_"+k)
	}

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
	feed := generateStandalonePage(articles, article, history)

	outD := GetConfig().OutputPath + "/"
	err := os.MkdirAll(outD, 0755)
	if err != nil {
		panic(err)
	}
	outName := outD + "feed.html"
	err1 := os.WriteFile(outName, feed, 0644)
	if err1 != nil {
		panic(err1)
	}
}

func generateFeedXML(articles Articles, fileName string) {
	if GetConfig().Verbose > 0 {
		fmt.Println("Generating feed: " + fileName)
	}
	feedUrl := GetConfig().SiteURL + "/" + "feed.xml"
	z := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
  <id>` + GetConfig().SiteURL + "/" + "index.html" + `</id>
  <link type="text/html" rel="alternate" href="` + feedUrl + `"/>
  <link type="application/atom+xml" rel="self" href="` + feedUrl + `"/>
  <title>` + GetConfig().SiteTitle + `</title>
  <updated>` + time.Now().Format("2006-01-02T15:04:05-07:00") + `</updated>`

	for _, e := range articles {
		url := GetConfig().SiteURL + filepath.Clean("/"+e.SrcDirectoryName+"/"+e.DstFileName)
		z += `
  <entry>
    <id>` + url + `</id>
    <link type="text/html" rel="alternate" href="` + url + `"/>
    <title>
        ` + e.Title + `
    </title>
    <updated>` + e.ModificationDate.Format("2006-01-02T15:04:05-07:00") + `</updated>`

		for _, t := range e.Tags {
			z += `<category scheme="` + GetConfig().SiteURL + `" term="` + t + `"/>`
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
	errMkdir := os.MkdirAll(GetConfig().OutputPath+"/feed", 0755)
	if errMkdir != nil {
		fmt.Println(errMkdir)
		panic(errMkdir)
	}
	feedName := filepath.Clean(GetConfig().OutputPath + "/feed/" + fileName + ".xml")
	err2 := os.WriteFile(feedName, []byte(z), 0644)
	if err2 != nil {
		fmt.Println(err2)
		panic(err2)
	}
}

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
			if GetConfig().Verbose > 0 {
				fmt.Println("Rendering article '" + article.Title + "'")
			}
			standalonePageContent := generateStandalonePage(articles, *article, article.Render())
			outD := filepath.Clean(GetConfig().OutputPath + "/" + article.SrcDirectoryName + "/")
			sendLiveUpdateViaWS(article.SrcFileName, article.Render())

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
	}
}

func generateStandalonePage(articles Articles, article Article, body string) []byte {
	buff := bytes.NewBufferString("")
	t, err := template.New("standalonePage.tmpl").
		ParseFiles("templates/standalonePage.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	relativeSrcRootPath, _ := filepath.Rel(article.SrcDirectoryName, "")
	relativeSrcRootPath = filepath.Clean(relativeSrcRootPath)
	//   fmt.Println(relativeSrcRootPath)

	//   fmt.Println("---------------")
	titleNAV := ""
	var prev string
	var next string
	//   fmt.Println(article.Title)
	p := articles.PrevArticle(&article)
	if p != nil {
		prev = filepath.Clean(relativeSrcRootPath + "/" + p.SrcDirectoryName + "/" + p.DstFileName)
		// link is active
		titleNAV +=
			`<span id="articleNavLeft"> <a href="` + prev + `"> 
      <span class="glyphiconLink glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article"> </span> prev. article
    </a> </span>`
	}
	n := articles.NextArticle(&article)
	if n != nil {
		// link is active
		next = filepath.Clean(relativeSrcRootPath + "/" + n.SrcDirectoryName + "/" + n.DstFileName)
		titleNAV +=
			`<span id="articleNavRight"><a href="` + next + `"> 
       next article <span class="glyphiconLink glyphicon glyphicon-chevron-right" aria-hidden="true" title="next article"></span>
    </a></span>`
	}
	seriesNAV := ""
	var sPrev string
	var sNext string

	if article.Series != "" {
		sp := articles.PrevArticleInSeries(&article)
		if sp != nil {
			sPrev = filepath.Clean(relativeSrcRootPath + "/" + sp.SrcDirectoryName + "/" + sp.DstFileName)
		}

		sn := articles.NextArticleInSeries(&article)
		if sn != nil {
			sNext = filepath.Clean(relativeSrcRootPath + "/" + sn.SrcDirectoryName + "/" + sn.DstFileName)
		}
		seriesNAV =
			`
      <div id="seriesContainer">
      <a href="` + relativeSrcRootPath + `/posts.html?filter=series::` + article.Series + `" title="article series ` + article.Series + `" class="seriesbtn btn btn-primary">` +
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

	var meta string
	var timeT time.Time

	if article.ModificationDate != timeT {
		meta += `<div id="date"><p><span id="lastupdated">` + article.ModificationDate.Format("2 Jan 2006") + `</span></p></div>`
	}

	if len(article.Tags) > 0 {
		meta += `<div id="tags"><p>` + tagToLinkList(&article) + `</p></div>`
	}

	noItems := struct {
		Title               string
		RelativeSrcRootPath string
		SiteURL             string
		SiteBrandTitle      string
		TitleNAV            string
		SeriesNAV           string
		Meta                string
		Anchorjs            bool
		Tocify              bool
		Timeline            bool
		Body                string
	}{
		Title:               article.Title,
		RelativeSrcRootPath: relativeSrcRootPath,
		SiteURL:             GetConfig().SiteURL,
		SiteBrandTitle:      GetConfig().SiteTitle,
		TitleNAV:            titleNAV,
		SeriesNAV:           seriesNAV,
		Meta:                meta,
		Anchorjs:            article.Anchorjs,
		Tocify:              article.Tocify,
		Timeline:            article.Timeline,
		Body:                body,
	}
	err = t.Execute(buff, noItems)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return buff.Bytes()
}

func Init() {
	pflag.String("input", "documents", "input directory  ('documents'') in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files")
	pflag.String("output", "output", "output directory ('output') all generated files will be stored there and all directories like css/ js/ images and fonts/ will be rsynced there")
	pflag.String("siteURL", "https://lastlog.de/blog", "The URL of the blog, for example: 'https://example.com/blog'")
	pflag.String("siteTitle", "lastlog.de/blog", "Title which is inserted top left, for example: 'lastlog.de/blog'")
	pflag.Int("verbose", 0, "verbosity level")
	pflag.Int("force", 0, "forced complete rebuild, not using cache")
	pflag.String("ListenAndServe", ":8000", "ip:port where pankat-server listens, for example: 'localhost:8000'")

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

	GetConfig().SiteURL = siteURL_
	GetConfig().InputPath = inputPath_
	GetConfig().OutputPath = outputPath_
	GetConfig().SiteURL = siteURL_
	GetConfig().SiteTitle = siteTitle_
	GetConfig().MyMd5HashMapJson = myMd5HashMapJson_
	GetConfig().Verbose = viper.GetInt("verbose")
	GetConfig().Force = viper.GetInt("force")
	GetConfig().ListenAndServe = viper.GetString("ListenAndServe")
}

func Rsync() {
	defer timeElapsed("rsync")()
	fmt.Println(color.GreenString("rsync"), "files to output directory")

	s, _ := os.Getwd()
	//fmt.Println("Getwd path is: ", s)

	v, _ := filepath.Rel(s, GetConfig().OutputPath)
	v = filepath.ToSlash(v)
	//fmt.Println("relative path is: ", v) // FIXME windows specific hack
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
	outIndexName := filepath.Clean(GetConfig().OutputPath + "/" + "index.html")
	errn := os.WriteFile(outIndexName, []byte(indexContent), 0644)
	if errn != nil {
		panic(errn)
	}
}

func UpdateBlog() {
	defer timeElapsed("UpdateBlog")()

	fmt.Println(color.GreenString("pankat-static"), "starting!")
	fmt.Println(color.YellowString("input Directory: "), GetConfig().InputPath)
	fmt.Println(color.YellowString("output Directory: "), GetConfig().OutputPath)

	// find all .mdwn files
	f := make([]string, 0)
	f = append(f, "")

	articles := GetTargets(".", f).FilterOutDrafts()

	// sort them by date
	sort.Sort(articles)

	fmt.Println(color.YellowString("GetTargets: found"), articles.Posts().Len(), color.YellowString("articles"))

	// override default values for posts
	for _, e := range articles.Posts() {
		e.Anchorjs = true
		e.Tocify = true
	}

	RenderPosts(articles)
	RenderTimeline(articles.Posts())
	RenderFeed(articles.Posts())
	SetMostRecentArticle(articles.Posts())
	Rsync()
}
