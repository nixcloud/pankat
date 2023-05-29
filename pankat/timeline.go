package pankat

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type MetaData struct {
	ArticleCount int
	Tags         map[string][]int
	Series       map[string][]int
	Years        map[int][]int
}

func (s Articles) CreateJSMetadata() MetaData {
	tagsMap := make(map[string][]int)
	seriesMap := make(map[string][]int)
	yearsMap := make(map[int][]int)
	for i, e := range s {
		m := e.ModificationDate
		year, err := strconv.Atoi(m.Format("2006"))
		if err == nil {
			if yearsMap[year] == nil {
				yearsMap[year] = []int{i}
			} else {
				yearsMap[year] = append(yearsMap[year], i)
			}
		}

		for _, t := range e.Tags {
			if tagsMap[t] == nil {
				tagsMap[t] = []int{i}
			} else {
				tagsMap[t] = append(tagsMap[t], i)
			}
		}
		z := s[i].Series
		if z != "" {
			if seriesMap[z] == nil {
				seriesMap[z] = []int{i}
			} else {
				seriesMap[z] = append(seriesMap[z], i)
			}
		}
	}
	return MetaData{len(s), tagsMap, seriesMap, yearsMap}
}

func RenderTimeline(articles Articles) {
	defer timeElapsed("RenderTimeline")()
	fmt.Println(color.YellowString("Rendering timeline into timeline.html"))

	var pageContent string
	var article Article

	article.Title = "all posts timeline"
	article.Timeline = true

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
	if Config().Verbose > 0 {
		fmt.Println(color.GreenString("tagsSlice"), tagsSlice)
	}

	pageContent += `<p id="tagCloud">`
	for _, e := range tagsSlice {
		pageContent += `<a class="tagbtn btn btn-primary" onClick="setFilter('tag::` + e.Key + `', 1)">` + e.Key + `</a>`
	}
	pageContent += `</p>`

	seriesSlice := rankByWordCount(seriesMap)
	if Config().Verbose > 0 {
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

	// create a map of years
	yearsMap := make(map[int]Articles)
	for _, article := range articles {
		yearsMap[article.ModificationDate.Year()] = append(yearsMap[article.ModificationDate.Year()], article)
	}

	keys := make([]int, 0, len(yearsMap))
	for k := range yearsMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(a, b int) bool {
		return keys[a] > keys[b]
	})
	var articleCount int

	genTimelineseries := func(articles Articles) string {
		var ret string

		for _, article := range articles {
			// hack to make tagToLinkList(...) work with relative directory ./ vs. ../
			var v Article
			v = *article

			ret += `
          <dt class="timeline-event posting_` + strconv.Itoa(articleCount) + `">` + article.Title + `</dt>
          <dd class="timeline-event-content posting_` + strconv.Itoa(articleCount) + `">
            <div class="postingsEntry">
              <p class="summary">` + article.Summary + ` <a href="` + filepath.ToSlash(article.DstFileName) + `">open complete article</a></p>
              <p class="tag">` + tagToLinkList(&v) + seriesToLinkList(&v) + `</p>
            </div>
            <br class="clear">
          </dd><!-- /.timeline-event-content -->`
			articleCount += 1
		}
		return ret
	}

	genYear := func(year int, articles Articles) string {
		var ret string
		ret += `<div class="timeline-wrapper pankat_year pankat_year_` + strconv.Itoa(year+1) + `">
		<dl class="timeline-series">
        <h2 class="timeline-time"><span>` + strconv.Itoa(year+1) + `</span></h2>`
		ret += genTimelineseries(articles)
		ret += `
		</dl><!-- /.timeline-series -->
		</div><!-- /.timeline-wrapper -->`
		return ret
	}

	for _, year := range keys {
		pageContent += genYear(year, yearsMap[year])
		if _, ok := yearsMap[year-1]; !ok {
			pageContent += genYear(year-1, Articles{})
		}
	}

	pageContent += `</div><!-- /.timeline-container -->`

	navTitleArticleSource := GenerateNavTitleArticleSource(articles, article, pageContent)
	standalonePageContent := GenerateStandalonePage(articles, article, navTitleArticleSource)

	outD := Config().DocumentsPath + "/"
	err = os.MkdirAll(outD, 0755)
	if err != nil {
		panic(err)
	}
	outName := outD + "timeline.html"
	err1 := os.WriteFile(outName, standalonePageContent, 0644)
	if err1 != nil {
		panic(err1)
	}
}
