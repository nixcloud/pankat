package pankat

var articlesCache ArticlesCache

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

type Articles []*Article

// https://gobyexample.com/sorting-by-functions
func (s Articles) Len() int {
	return len(s)
}
func (s Articles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Articles) Less(i, j int) bool {
	return s[i].ModificationDate.After(s[j].ModificationDate)
}

func (s Articles) NextArticle(a *Article) *Article {
	for i, elem := range s {
		if elem.SrcFileName == a.SrcFileName {
			if i-1 >= 0 {
				return s[i-1]
			}
		}
	}
	return nil
}

func (s Articles) PrevArticle(a *Article) *Article {
	for i, elem := range s {
		if elem.SrcFileName == a.SrcFileName {
			if i+1 < len(s) {
				return s[i+1]
			}
		}
	}
	return nil
}

func (s Articles) PrevArticleInSeries(a *Article) *Article {
	q := s.FindArticleWithSeries(a.Series)
	if len(q) == 0 {
		return nil
	}
	z := q.PrevArticle(a)
	return z
}

func (s Articles) NextArticleInSeries(a *Article) *Article {
	q := s.FindArticleWithSeries(a.Series)
	if len(q) == 0 {
		return nil
	}
	z := q.NextArticle(a)
	return z
}

func (s Articles) GenerateArticleNavigation(article *Article) string {
	articles := s
	//   fmt.Println("---------------")
	titleNAV := ""
	//   fmt.Println(article.Title)
	p := articles.PrevArticle(article)
	if p != nil {
		// link is active
		titleNAV +=
			`<span id="articleNavLeft"> <a href="` + p.DstFileName + `"> 
      <span class="glyphiconLink glyphicon glyphicon-chevron-left" aria-hidden="true" title="previous article"> </span> prev. article
    </a> </span>`
	}
	n := articles.NextArticle(article)
	if n != nil {
		// link is active
		titleNAV +=
			`<span id="articleNavRight"><a href="` + n.DstFileName + `"> 
        next article <span class="glyphiconLink glyphicon glyphicon-chevron-right" aria-hidden="true" title="next article"></span>
    </a> </span>`
	}

	return titleNAV
}

func (s Articles) GenerateArticleSeriesNavigation(article *Article) string {
	articles := s
	seriesNAV := ""
	var sPrev string
	var sNext string

	if article.Series != "" {
		sp := articles.PrevArticleInSeries(article)
		if sp != nil {
			sPrev = sp.DstFileName
		}

		sn := articles.NextArticleInSeries(article)
		if sn != nil {
			sNext = sn.DstFileName
		}
		seriesNAV =
			`
      <div id="seriesContainer">
      <a href="timeline.html?filter=series::` + article.Series + `" title="article series ` + article.Series + `" class="seriesbtn btn btn-primary">` +
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
	return seriesNAV
}

// FIXME refactor this Targets(), this should be done during Add(article)
func (s Articles) Targets() Articles {
	var _filtered Articles
	for _, e := range s {
		e.SourceReference = true
		e.WebsocketSupport = true
		if e.SpecialPage == true {
			e.Anchorjs = false
			e.Tocify = false
		} else {
			e.Anchorjs = true
			e.Tocify = true
		}
		_filtered = append(_filtered, e)
	}
	return _filtered
}

func (s Articles) FindArticleWithSeries(t string) Articles {
	var _filtered Articles
	for _, e := range s {
		if e.Series == t {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) FindArticlesWithTag(t string) Articles {
	var _filtered Articles
	for _, e := range s {
		if contains(e.Tags, t) {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) FilterOutDrafts() Articles {
	var _filtered Articles
	for _, e := range s {
		if e.Draft == false {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}

func (s Articles) FilterOutSpecialPages() Articles {
	var _filtered Articles
	for _, e := range s {
		if e.SpecialPage == false {
			_filtered = append(_filtered, e)
		}
	}
	return _filtered
}
