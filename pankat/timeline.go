package pankat

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strconv"
)

func RenderTimeline(articles Articles) {
	defer timeElapsed("RenderTimeline")()
	fmt.Println(color.YellowString("Rendering timeline in posts.html"))

	var pageContent string
	var article Article

	article.Title = "all posts timeline"
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
              <p class="summary">` + article.Summary + ` <a href="` + filepath.Clean(article.DstFileName) + `">open complete article</a></p>
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
      </div>
`

	navTitleArticleSource := GenerateNavTitleArticleSource(articles, article, pageContent)
	standalonePageContent := GenerateStandalonePage(articles, article, navTitleArticleSource)

	outD := GetConfig().DocumentsPath + "/"
	err = os.MkdirAll(outD, 0755)
	if err != nil {
		panic(err)
	}
	outName := outD + "posts.html"
	err1 := os.WriteFile(outName, standalonePageContent, 0644)
	if err1 != nil {
		panic(err1)
	}
}
