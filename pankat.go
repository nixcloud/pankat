package main

import (
	"./articles" // pankat
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
  htemplate "html/template"
	"time"
)

// maybe migrate to gopkg.in/alecthomas/kingpin.v2 as flag replacement
var inputPath string
var iArg = flag.String("i", "", "input directory")
var outputPath string
var oArg = flag.String("o", "", "output directory")

// var mode string
// var modeArg = flag.String("mode", "blog", "operation mode: 'blog' or 'wiki'")

var SiteURL = "https://lastlog.de/blog"
var SiteBrandTitle = "lastlog.de/blog"

func tagToLinkList(a *pankat.Article) string {
	var tags []string
	tags = a.Tags
	var output string
	for _, e := range tags {
		//     fmt.Println("----------------")
		//     fmt.Println(outputPath)
		//     fmt.Println(a.SrcDirectoryName)

		//HACK should be moved to pankat.Articles
		relativeSrcRootPath, _ := filepath.Rel(a.SrcDirectoryName, "")
		relativeSrcRootPath = path.Clean(relativeSrcRootPath)

		output += `<a href="` + relativeSrcRootPath + `/posts.html?tag=` + e + `" class="tagbtn btn btn-primary">` + e + `</a>`
	}
	return output
}

func main() {
	flag.Usage = func() {
		a :=
			`Usage of '` + filepath.Base(os.Args[0]) + `', the ultimate pandoc wiki/blog generator, is: 
  -i input directory  (must be given)
     in this directory it is expected to find about.mdwn and posts/ among other top level *.mdwn files
  -o output directory (must be given)
     all generated files will be stored there and all directories like css/ js/ images and fonts/ will be rsynced there

  example: ./` + filepath.Base(os.Args[0]) + ` -i documents -o output
`
		fmt.Fprintf(os.Stderr, a)
		os.Exit(1)
	}
	flag.Parse()

	if *iArg == "" || *oArg == "" {
		flag.Usage()
	}
	inputPath = path.Clean(*iArg + "/")
	fmt.Println("srcDirectory: ", inputPath)
	outputPath = *oArg + "/"
  
	tt, err := filepath.Abs(outputPath)
	outputPath = tt

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("outDirectory: ", outputPath)

	err1 := os.Chdir(inputPath)
	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}
	
	// find all .mdwn files
	f := make([]string, 0)
	f = append(f, "")
	articlesAll := getTargets(".", f)

	articlesTopLevel := articlesAll.TopLevel().FilterByDraft()
	articlesPosts := articlesAll.Posts().FilterByDraft()

	// sort them by date
	sort.Sort(pankat.Articles(articlesPosts))

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
	renderPosts(articlesTopLevel)

	// generate posts.html (timeline)
	renderPostsTimeline(articlesPosts)

	// generate articles
	renderPosts(articlesPosts)

	// copy static files as fonts/css/js to the output folder
	commands := []string{
		"rsync -av --delete --relative js " + outputPath,
		"rsync -av --delete --relative css " + outputPath,
		"rsync -av --delete --relative fonts " + outputPath,
		"rsync -av --delete --relative images " + outputPath,
		"rsync -av --delete --relative posts/media " + outputPath,
	}

	for _, el := range commands {
		parts := strings.Fields(el)
		head := parts[0]
		parts = parts[1:len(parts)]
		out, err := exec.Command(head, parts...).Output()
		if err != nil {
			fmt.Println(err)
		}
		//       fmt.Println("executing: ", el)
		//       fmt.Println(string(out))
		_ = out
	}

	//BUG  md5 optimization forgets to update former last article to have a 'next article' link if new article is added
	
	// generate rss/atom feed
	// FIXME create feed per tag
	renderFeed(articlesPosts)
  
	// BUG fix history writing
	//   example: 1. go to article https://lastlog.de/blog/posts/tour_of_nix.html
	//            2. click on an article tag https://lastlog.de/blog/posts.html?tag=emscripten
	//            3. then try 'back' button, which fails!
  //      maybe use backbone.js for that?

	// FIXME donation button
  // FIXME next/last hover shadow
  
  // FIXME use h1 only for title, see http://pandoc.org/scripting.html filter
  
  //////////////////////////////////////// main features ///////////////////////////////////////////////////
	
  // SECURITY secure pandoc from passing <script> and other evil <html tags>
  //          find a filter system for evil html tags like <script>  
  
  // https://www.overleaf.com/4344023pmjpgq#/12921720/
	// FIXME - gocraft/web ansprechen
  //       - git backend ansprechen
  //       - leaps backend ansprechen
  //       - websockets preview mit long-polling
  //       - lokales speichern von artikeln, wenn ./pankat -daemon -i documents -o output/ verwendet wird
  
	// FIXME - integrate a wiki switch
  // FIXME - pankat release 


	//////////////////////////////////////// /main features ///////////////////////////////////////////////////

	// FIXME for each article:
	//  - rework warning/info/danger/error ...
	//  - Summary for each article
	//  - rewrite title names
	//  - fix images, add class="noFancy"
	//  - check h1,h2,...
	//  - use <div class="warn">...</div>
	//  - check [[!series ogre]] for other series like qt
	//  - libnoise_viewer.html fix video width

	// - commit history using git and add revert link like ikiwiki does FIXME
	// - implement comment system FIXME
	//   see example: https://www.reddit.com/r/golang/comments/1xbxzk/default_value_in_structs/

	// BUG pandoc integration with parser '-s' of html head/body and migration to the go template

	// FIXME create a [[!pandocFormat mdwn]] plugin which makes more pandoc dialects available
  
	mostRecentArticle := ""
	if len(articlesPosts) > 0 {
		mostRecentArticle = path.Clean(articlesPosts[0].SrcDirectoryName + "/" + articlesPosts[0].DstFileName)
	} else {
		mostRecentArticle = "posts.html"
	}
	indexContent :=
		`
<html>
<meta http-equiv="refresh" content="0; url=` + mostRecentArticle + `" />
</html>
`

	outIndexName := path.Clean(outputPath + "/" + "index.html")
	errn := ioutil.WriteFile(outIndexName, []byte(indexContent), 0644)
	if errn != nil {
		panic(errn)
	}
}

// scan the direcotry for .mdwn files recurively
func getTargets(path string, ret []string) pankat.Articles {
	var A pankat.Articles
	entries, _ := ioutil.ReadDir(path)
	for _, entry := range entries {
		buf := path + "/" + entry.Name()
		//     fmt.Println("reading buf: ", buf)
		if entry.IsDir() {
			if entry.Name() == ".git" {
				continue
			}
			ret = append(ret, buf)
			//       fmt.Println(buf)
			n := getTargets(buf, ret)
			A = append(A, n...)
		} else {
			if strings.HasSuffix(entry.Name(), ".mdwn") {
				var a pankat.Article
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
					os.Exit(1)
				}
				defer fh.Close()

				_article, err := ioutil.ReadAll(f)
				if err != nil {
					fmt.Println(err)
					panic(err)
					os.Exit(1)
				}

				_article = filterDocument(_article, &a)

				a.Hash = md5.Sum(_article)
				a.Article = _article
				A = append(A, &a)
			}
		}
	}
	return A
}

func filterDocument(_article []byte, article *pankat.Article) []byte {
	var _articlePostprocessed []byte

	re := regexp.MustCompile("\\[\\[!(.*?)\\]\\]")
	z := re.FindAllIndex(_article, -1)

	prevPos := 0
	for i := 0; i <= len(z); i++ {
		if i == len(z) {
			_articlePostprocessed = append(_articlePostprocessed, _article[prevPos:len(_article)]...)
			break
		}
		n := z[i]

		// include normal content (not plugin processed)
		if prevPos != n[0] {
			_articlePostprocessed = append(_articlePostprocessed, _article[prevPos:n[0]]...)
		}

		// include plugin processed stuff
		_articlePostprocessed = append(_articlePostprocessed, callPlugin(_article[n[0]:n[1]], article)...)
		prevPos = n[1]
	}
	return _articlePostprocessed
}

func callPlugin(in []byte, article *pankat.Article) []byte {
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
		return z
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
			article.Series = strings.Join(f[1:len(f)], " ")
		}
	case "tag":
		if len(f) > 1 {
			article.Tags = f[1:len(f)]
		}
	case "draft":
		article.Draft = true

	case "img":
		b := strings.Join(f[1:len(f)], " ")
		//      fmt.Println("\n------------\n", article.SrcDirectoryName)
		//      fmt.Println(f[1])

		//HACK should be moved to pankat.Articles
		relativeSrcRootPath, _ := filepath.Rel(article.SrcDirectoryName, "./posts")
		relativeSrcRootPath = path.Clean(relativeSrcRootPath)
		//      fmt.Println(relativeSrcRootPath)

		o := `<a href="` + relativeSrcRootPath + "/" + f[1] + `"><img src=` + relativeSrcRootPath + "/" + b + `></a>`
		output = []byte(o)

	case "summary":
		article.Summary = strings.Join(f[1:len(f)], " ")
	default:
		fmt.Println(name, " plugin, called from ", article.SrcFileName, " NOT supported")
	}

	return output
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

func renderPostsTimeline(articles pankat.Articles) {
	// http://codepen.io/jplhomer/pen/lgfus
	var history string
	var article pankat.Article
	article.SrcDirectoryName = ""

	b, err := json.Marshal(articles.TagUsage())
	if err != nil {
		fmt.Println("json.Marshal error:", err)
	}
	history += `<script type="application/json" id="tagsMap">` + string(b) + `</script>`

	tagsMap := make(map[string]int)

	for _, a := range articles {
		for _, t := range a.Tags {
			tagsMap[t] = tagsMap[t] + 1
		}
	}
	// sort the tags
	tagsSlice := rankByWordCount(tagsMap)

	history += `<p id="tagCloud">`
	for _, e := range tagsSlice {
		zz := "'" + e.Key + "'"
		history += `<a class="tagbtn btn btn-primary" onClick="showTag(` + zz + `)">` + e.Key + `</a>`
	}
	history += `</p>`

	history += `
    <p class="lead">A timeline showing all blog postings.</p>

    <div id="timelineFilter" style="display: none">
      <p>Currently this filter is set:</p>
      <a id="timelineButton"  class="tagbtn btn btn-primary">bios</a> <span id="timelineFilterCancel" class="glyphicon glyphicon-remove" onclick="showTag('')"></span>
    </div>
    
    <div id="timeline" class="timeline-container">
      <a class="timeline-toggle btn btn-primary">+expand all</a>

      <br class="clear">
      <div class="timeline-wrapper">
      <dl class="timeline-series">`

	var year string
	
	for i, e := range articles {
		if i == 0 {
      v := e.ModificationDate.Add(1000*1000*1000*60*60*24*365) // add one year
			year = v.Format("2006")
			history += `<h2 class="timeline-time"><span>` + year + `</span></h2>`
      year = e.ModificationDate.Format("2006")
		}

		//     fmt.Println("----")
		//     fmt.Println("  ", e.Title)
		//     fmt.Println("  ", e.SrcDirectoryName)
		//     fmt.Println("  ", inputPath)

		if year != e.ModificationDate.Format("2006") {
			history += `
         </dl><!-- /.timeline-series -->
       </div><!-- /.timeline-wrapper -->
       <div class="timeline-wrapper">`

			history += `<h2 class="timeline-time"><span>` + year + `</span></h2>`
			history += `<dl class="timeline-series">`
			year = e.ModificationDate.Format("2006")
		}

		// a hacky but straight-forward way to make tagToLinkList(...) work by
		// fooling a different base article
		var v pankat.Article
		v = *e
		v.SrcDirectoryName = ""

		//     <h3>` + e.ModificationDate.Format("2 Jan 2006") + `</h3>
		//     <span class="glyphicon glyphicon-chevron-link" aria-hidden="true" title="article"></span>
		history += `
          <dt id="` + strconv.Itoa(i) + `" class="timeline-event"><a>` + e.Title + `</a></dt>
          <dd class="timeline-event-content" id="` + strconv.Itoa(i) + "EX" + `">
            <div class="postingsEntry">
              <p class="summary">` + e.Summary + ` <a href="` + path.Clean(e.SrcDirectoryName+"/"+e.DstFileName) + `">open complete article</a></p>
              <p class="tag">` + tagToLinkList(&v) + `</p>
            </div>
            <br class="clear">
          </dd><!-- /.timeline-event-content -->`
	}

	history += `
        </dl><!-- /.timeline-series -->
      </div><!-- /.timeline-wrapper -->
      
      <script>
      var tagsMap
      var showTag = function(tagName) {
        var count = tagsMap.ArticleCount
        if (tagName === "") {
          for (i=0; i < count; i++) { var n = "#" + i; $(n).css('display', 'block'); }
          $('#timelineFilter').fadeOut("slow");
          window.history.pushState('', '',  window.location.pathname);
        } else {
          for (i=0; i < count; i++) { var n = "#" + i; $(n).css('display', 'none'); }
          if (typeof(tagsMap.Tags[tagName]) !== "undefined") {
            for (i=0; i < tagsMap.Tags[tagName].length; i++) { 
              var n = "#" + (tagsMap.Tags[tagName])[i]; 
              $(n).css('display', 'block'); 
            }
          }
          $('#timelineFilter').css('display','block')
          $('#timelineButton')[0].innerHTML=tagName
          window.history.pushState('', '',  window.location.pathname + '?tag=' + tagName);
        }
      }
      $(document).ready(function() {
        tagsMap = JSON.parse(document.getElementById('tagsMap').innerHTML)
        var tag = getURLParameter("tag");
        if (typeof(tag) !== "undefined" && tag !== null) {
          console.log("read() : " + tag)
          showTag(tag);
        }
      });
      function getURLParameter(name) {
        return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)').exec(location.search)||[,""])[1].replace(/\+/g, '%20'))||null
      }
      window.addEventListener("popstate", function() {
        var tag = getURLParameter("tag");
        if (typeof(tag) !== "undefined" && tag !== null) {
          console.log("addEvent() : " + tag)
          showTag(tag);
        }
      });
      </script>
      </div>
`

	article.Title = "posts - timeline"
	article.Timeline = true

	page := generateStandalonePage(articles, article, history)

	outD := outputPath + "/"
	os.MkdirAll(outD, 0755)
	outName := outD + "posts.html"
	err1 := ioutil.WriteFile(outName, page, 0644)
	if err1 != nil {
		panic(err1)
	}
}

func renderFeed(articles pankat.Articles) {
  feedUrl := SiteURL + "/" + "feed.xml"
    z := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
  <id>`+SiteURL + "/" + "index.html"+`</id>
  <link type="text/html" rel="alternate" href="`+feedUrl+`"/>
  <link type="application/atom+xml" rel="self" href="`+feedUrl+`"/>
  <title>`+ SiteBrandTitle +`</title>
  <updated>`+ time.Now().Format("2006-01-02T15:04:05-07:00") +`</updated>`
  
    for _, e := range articles {
      url := SiteURL + path.Clean("/" + e.SrcDirectoryName + "/" + e.DstFileName)
      z += `
  <entry>
    <id>`+ url +`</id>
    <link type="text/html" rel="alternate" href="`+ url +`"/>
    <title>
        ` + e.Title + `
    </title>
    <updated>` + e.ModificationDate.Format("2006-01-02T15:04:05-07:00") + `</updated>`
  
  for _, t := range e.Tags {
    z+=`<category scheme="`+SiteURL+`" term="`+t+`"/>`
  }

    z += `<author>
      <name>qknight</name>
      <uri>https://github.com/qknight</uri>
    </author>
    <content type="html">` + htemplate.HTMLEscaper(string(e.RenderedArticle)) + `</content>
  </entry>`
    }
    z += `</feed>`
    
    outName := path.Clean(outputPath + "/" + "feed.xml")
    err2 := ioutil.WriteFile(outName, []byte(z), 0644)
    if err2 != nil {
      fmt.Println(err2)
      panic(err2)
    }
}

type CachedArticleData struct {
  Article string
  Hash [md5.Size]byte
}

func renderPosts(articles pankat.Articles) {
// BUG if hash found, article won't be rendered (that was the idea), but then the feed lacks <content>...</content>
//     so article has to be cached among the md5 sum
  myMd5HashMap := make(map[string]CachedArticleData)
  myMd5HashMapJson := path.Clean(outputPath + "/" + ".MyMd5HashMap.json")
  
  b, errReadFile := ioutil.ReadFile(myMd5HashMapJson)
  if errReadFile != nil {
    fmt.Println(errReadFile)
  } else {
    jBuff := bytes.NewBufferString(string(b))
    
    dec := json.NewDecoder(jBuff)
    
    if err := dec.Decode(&myMd5HashMap); err != nil {
      fmt.Println(err)
    }
  }
  
	for _, e := range articles {
    key := path.Clean(e.SrcDirectoryName + "/" + e.SrcFileName)
    if (myMd5HashMap[key].Hash == e.Hash) {
//       fmt.Println(e.Hash, " already generated, not generating again")
      e.RenderedArticle = myMd5HashMap[key].Article
      continue
    } else {
      //myMd5HashMap[key] = CachedArticleData{"",""}
    }
		fmt.Println(e.Title)

		pandocProcess := exec.Command("pandoc", "-f", "markdown", "-t", "html5", "--highlight-style", "kate")
		stdin, err := pandocProcess.StdinPipe()
		if err != nil {
			fmt.Println(err)
			continue
		}

		buff := bytes.NewBufferString("")
		pandocProcess.Stdout = buff
		pandocProcess.Stderr = os.Stderr

		err1 := pandocProcess.Start()
		if err1 != nil {
			fmt.Println("An error occured: ", err1)
			continue
		}

		io.WriteString(stdin, string(e.Article))
		stdin.Close()
		pandocProcess.Wait()

    e.RenderedArticle = string(buff.Bytes())
    
		standalonePageContent := generateStandalonePage(articles, *e, string(buff.Bytes()))

		outD := path.Clean(outputPath + "/" + e.SrcDirectoryName + "/")
		//     fmt.Println(outD + e.DstFileName)
		//     fmt.Println(string(standalonePageContent))

		os.MkdirAll(outD, 0755)

		// write to disk
		outName := path.Clean(outD + "/" + e.DstFileName)
		err2 := ioutil.WriteFile(outName, standalonePageContent, 0644)
		if err2 != nil {
			fmt.Println(err2)
			panic(e)
		}
		// FIXME add new article to cache
		myMd5HashMap[key] = CachedArticleData{e.RenderedArticle, md5.Sum([]byte(e.Article))}
	}
	
// 	fmt.Println("All article hashes:")
// 	for k,v := range myMd5HashMap {
//     fmt.Println(k, " ", v)
// 	}
	
	jsonBuff := bytes.NewBufferString("")
  enc := json.NewEncoder(jsonBuff)

	if errEnc := enc.Encode(&myMd5HashMap); errEnc != nil {
    fmt.Println(errEnc)
  }
//   fmt.Println(string(jsonBuff.Bytes()))
  
  errn := ioutil.WriteFile(myMd5HashMapJson, jsonBuff.Bytes(), 0644)
  if errn != nil {
    panic(errn)
  }
}

func generateStandalonePage(articles pankat.Articles, article pankat.Article, body string) []byte {
	buff := bytes.NewBufferString("")
	t, err := template.New("standalonePage.tmpl").
		ParseFiles("templates/standalonePage.tmpl")
	if err != nil {
		fmt.Println(err)
		panic(err)
		os.Exit(1)
	}

	//HACK should be moved to pankat.Articles
	relativeSrcRootPath, _ := filepath.Rel(article.SrcDirectoryName, "")
	relativeSrcRootPath = path.Clean(relativeSrcRootPath)
	//   fmt.Println(relativeSrcRootPath)

	//   fmt.Println("---------------")
	titleNAV := ""
	var prev string
	var next string
	//   fmt.Println(article.Title)
	p := articles.PrevArticle(&article)
	if p != nil {
		prev = path.Clean(relativeSrcRootPath + "/" + p.SrcDirectoryName + "/" + p.DstFileName)
		// link is active
		titleNAV +=
		  `<span id="articleNavLeft"> <a href="` + prev + `"> 
      <span class="glyphiconLink glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article"> </span> prev. article
    </a> </span>`
	} 
	n := articles.NextArticle(&article)
	if n != nil {
		// link is active
		next = path.Clean(relativeSrcRootPath + "/" + n.SrcDirectoryName + "/" + n.DstFileName)
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
			sPrev = path.Clean(relativeSrcRootPath + "/" + sp.SrcDirectoryName + "/" + sp.DstFileName)
		}

		sn := articles.NextArticleInSeries(&article)
		if sn != nil {
			sNext = path.Clean(relativeSrcRootPath + "/" + sn.SrcDirectoryName + "/" + sn.DstFileName)
		}
		seriesNAV =
			`
      <div id="seriesContainer">
      <a href="` + relativeSrcRootPath + `/posts.html?series=` + article.Series + `" title="article series `+article.Series+`">` +
				article.Series + `</a>
        <header class="seriesHeader">
          <div id="seriesLeft">`
     if sp != nil {
       seriesNAV += `<a href="` + sPrev + `">` +
      `<span class="glyphiconLinkSeries glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article in series"></span>
            </a> `
     }
     seriesNAV +=   `  </div>
          <div id="seriesRight">`
     if sn != nil {     
       seriesNAV +=  `   <a href="` + sNext + `">
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
		SiteURL:             SiteURL,
		SiteBrandTitle:      SiteBrandTitle,
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
