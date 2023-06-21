package db

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func compareTagNames(a []Tag, b []Tag) error {
	if len(a) != len(b) {
		return errors.New("length of tags is not equal")
	}
	for i, v := range a {
		if v.Name != b[i].Name {
			s := fmt.Sprintf("tag names are not equal: %s != %s", v.Name, b[i].Name)
			return errors.New(s)
		}
	}
	return nil
}

func TestArticlesDatabase(t *testing.T) {
	articlesDb := NewArticlesDb()

	const longForm = "2006-01-02 15:04"
	time1, _ := time.Parse(longForm, "2019-01-01 00:00")
	article1 := Article{Title: "foo", ModificationDate: time1, Summary: "foo summary", Tags: []Tag{{Name: "Linux"}, {Name: "Go"}},
		SrcFileName: "/home/user/documents/foo.mdwn", DstFileName: "/home/user/documents/foo.html"}
	time2, _ := time.Parse(longForm, "2022-01-01 00:00")
	article2 := Article{Title: "bar", ModificationDate: time2, Series: "Linuxseries", Summary: "bar summary", Tags: []Tag{{Name: "SteamDeck"}, {Name: "Gorilla"}},
		SrcFileName: "/home/user/documents/bar.mdwn", DstFileName: "/home/user/documents/bar.html"}
	time3, _ := time.Parse(longForm, "2010-01-01 00:00")
	article3 := Article{Title: "batz", ModificationDate: time3, Series: "Linuxseries", Summary: "batz summary", Tags: []Tag{{Name: "Linux"}, {Name: "Go"}, {Name: "UniqueTag"}},
		SrcFileName: "/home/user/documents/batz.mdwn", DstFileName: "/home/user/documents/batz.html"}
	time4, _ := time.Parse(longForm, "2024-01-01 00:00")
	article4 := Article{Draft: true, Title: "draft", ModificationDate: time4, Summary: "draft summary", Tags: []Tag{{Name: "Go"}, {Name: "Linux"}},
		SrcFileName: "/home/user/documents/mydraft.mdwn", DstFileName: "/home/user/documents/mydraft.html"}
	time5, _ := time.Parse(longForm, "2024-01-01 00:00")
	article5 := Article{SpecialPage: true, Title: "draft", ModificationDate: time5,
		SrcFileName: "/home/user/documents/about.mdwn", DstFileName: "/home/user/documents/about.html"}

	// Insert the article into the database
	err := articlesDb.Add(&article1)
	if err != nil {
		panic(err)
	}
	err = articlesDb.Add(&article2)
	if err != nil {
		panic(err)
	}
	err = articlesDb.Add(&article3)
	if err != nil {
		panic(err)
	}
	err = articlesDb.Add(&article4)
	if err != nil {
		panic(err)
	}
	err = articlesDb.Add(&article5)
	if err != nil {
		panic(err)
	}

	// update item ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	err = articlesDb.Add(&article5)
	if err != nil {
		panic(err)
	}

	queryAll, err := articlesDb.QueryAll()
	assert.True(t, err == nil)
	assert.Equal(t, len(queryAll), 5)

	allArticles, err := articlesDb.Articles()
	assert.True(t, err == nil)
	assert.Equal(t, len(allArticles), 3)

	allarticles, err := articlesDb.QueryAll()
	assert.True(t, err == nil)
	assert.Equal(t, len(allarticles), 5)

	drafts, err := articlesDb.Drafts()
	assert.True(t, err == nil)
	assert.Equal(t, len(drafts), 1)

	specialpages, err := articlesDb.SpecialPages()
	assert.True(t, err == nil)
	assert.Equal(t, len(specialpages), 1)

	queryBySrcFileName, err := articlesDb.QueryRawBySrcFileName("/home/user/documents/bar.mdwn")
	assert.True(t, err == nil)
	assert.Equal(t, queryBySrcFileName.Title, "bar")
	err = compareTagNames(queryBySrcFileName.Tags, []Tag{{Name: "SteamDeck"}, {Name: "Gorilla"}})
	assert.NoError(t, err)

	mostRecentArticle, err := articlesDb.MostRecentArticle()
	assert.True(t, err == nil)
	assert.Equal(t, mostRecentArticle.SrcFileName, "/home/user/documents/bar.mdwn")

	////// Find next/previous article ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	assert.Equal(t, mostRecentArticle.ID, uint(2))
	nextArticle, err := articlesDb.NextArticle(mostRecentArticle)
	assert.Nil(t, nextArticle)
	assert.Error(t, err, "no next article")

	prevArticle, err := articlesDb.PrevArticle(mostRecentArticle)
	assert.Nil(t, err)
	assert.Equal(t, prevArticle.SrcFileName, "/home/user/documents/foo.mdwn")

	// Query articles by tag ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	tagName := "SteamDeck"
	taggedArticles, err := articlesDb.ArticlesByTag(tagName)
	assert.Nil(t, err)
	assert.Equal(t, len(taggedArticles), 1)
	assert.Equal(t, len(taggedArticles[0].Tags), 2)
	err = compareTagNames(taggedArticles[0].Tags, []Tag{{Name: "SteamDeck"}, {Name: "Gorilla"}})
	assert.NoError(t, err)

	tagName = "UniqueTag"
	taggedArticles, err = articlesDb.ArticlesByTag(tagName)
	assert.Nil(t, err)
	assert.Equal(t, len(taggedArticles), 1)

	tagName = "Linux"
	taggedArticles, err = articlesDb.ArticlesByTag(tagName)
	assert.Nil(t, err)
	assert.Equal(t, len(taggedArticles), 2)

	// Query all tags ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	tags, err := articlesDb.Tags()
	assert.Nil(t, err)
	assert.Equal(t, len(tags), 5)
	assert.Equal(t, tags, []string{"Linux", "Go", "SteamDeck", "Gorilla", "UniqueTag"})

	// Query articles by series /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	articlesBySeries, err := articlesDb.ArticlesBySeries("Linuxseries")
	assert.Nil(t, err)
	assert.Equal(t, len(articlesBySeries), 2)
	assert.Equal(t, articlesBySeries[0].SrcFileName, "/home/user/documents/bar.mdwn")
	assert.Equal(t, articlesBySeries[1].SrcFileName, "/home/user/documents/batz.mdwn")

	////// Find next/previous article in series //////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	nextArticleInSeries, err := articlesDb.NextArticleInSeries(articlesBySeries[1])
	assert.Equal(t, articlesBySeries[1].SrcFileName, "/home/user/documents/batz.mdwn")
	assert.Nil(t, err)
	assert.Equal(t, nextArticleInSeries.SrcFileName, "/home/user/documents/bar.mdwn")

	nextArticleInSeries, err = articlesDb.NextArticleInSeries(articlesBySeries[0])
	assert.Equal(t, articlesBySeries[0].SrcFileName, "/home/user/documents/bar.mdwn")
	assert.Error(t, err, "no next article in series found")

	prevArticleInSeries, err := articlesDb.PrevArticleInSeries(articlesBySeries[0])
	assert.Nil(t, err)
	assert.Equal(t, prevArticleInSeries.SrcFileName, "/home/user/documents/batz.mdwn")

	prevArticleInSeries, err = articlesDb.PrevArticleInSeries(articlesBySeries[1])
	assert.Error(t, err, "no prev article in series found")

	// Query all series /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	series, err := articlesDb.Series()
	assert.Nil(t, err)
	assert.Equal(t, len(series), 1)
	assert.Equal(t, series, []string{"Linuxseries"})

	// delete item //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	err = articlesDb.Del("/home/user/documents/mydraft.mdwn")
	if err != nil {
		panic(err)
	}
	all2, err := articlesDb.QueryAll()
	assert.True(t, err == nil)
	assert.Equal(t, len(all2), 4)
}
