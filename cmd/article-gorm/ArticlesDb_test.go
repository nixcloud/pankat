package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Article struct {
	gorm.Model
	SrcFileName      string `gorm:"uniqueIndex"`
	DstFileName      string
	Title            string
	ModificationDate time.Time
	Summary          string
	Tags             []Tag `gorm:"ForeignKey:TagId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Series           string
	SpecialPage      bool
	Draft            bool
	Anchorjs         bool
	Tocify           bool
	Timeline         bool
	ShowSourceLink   bool
	LiveUpdates      bool
	Evaluated        bool
}

type Tag struct {
	gorm.Model
	TagId uint
	Name  string
}

type ArticlesDb struct {
	db *gorm.DB
}

func NewArticlesDb() *ArticlesDb {
	// Open a new SQLite database connection or create one if it doesn't exist
	db, err := gorm.Open(sqlite.Open("pankat-sqlite3.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// remove old entries
	db.Migrator().DropTable(&Article{}, &Tag{})

	// Auto-migrate the table
	err = db.AutoMigrate(&Article{}, &Tag{})
	if err != nil {
		panic(err)
	}
	return &ArticlesDb{db: db}
}

func (a *ArticlesDb) Add(article *Article) error {
	result := a.db.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).FirstOrCreate(&article)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *ArticlesDb) Del(SrcFileName string) error {
	result := a.db.Where("src_file_name = ?", SrcFileName).Delete(&Article{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *ArticlesDb) QueryAll() ([]Article, error) {
	var articles []Article
	result := a.db.Preload("Tags").Order("modification_date DESC").Find(&articles)
	if result.Error != nil {
		return []Article{}, nil
	}
	return articles, nil
}

func (a *ArticlesDb) QueryRawBySrcFileName(SrcFileName string) (*Article, error) {
	var res Article
	result := a.db.Preload("Tags").Where("src_file_name = ?", SrcFileName).First(&res)
	if result.Error != nil {
		return nil, errors.New("article not found")
	}
	return &res, nil
}

func (a *ArticlesDb) Articles() ([]Article, error) {
	var res []Article
	result := a.db.Preload("Tags").Order("modification_date DESC").Where("draft = ?", false).Where("special_page = ?", false).Find(&res)
	if result.Error != nil {
		return []Article{}, nil
	}
	return res, nil
}

func (a *ArticlesDb) MostRecentArticle() (Article, error) {
	var res Article
	result := a.db.Preload("Tags").Order("modification_date DESC").Where("draft = ?", false).Where("special_page = ?", false).First(&res)
	if result.Error != nil {
		return Article{}, nil
	}
	return res, nil
}

func (a *ArticlesDb) NextArticle(article Article) (*Article, error) {
	var res Article
	result := a.db.Preload("Tags").Where("draft = ? AND special_page = ? AND modification_date >= ?", false, false, article.ModificationDate).
		Where("id != ?", article.ID).
		Order("modification_date DESC").
		Limit(1).
		Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// return no error
		return nil, errors.New("no next article")
	}
	return &res, nil
}

func (a *ArticlesDb) PrevArticle(article Article) (*Article, error) {
	var res Article
	result := a.db.Preload("Tags").Where("draft = ? AND special_page = ? AND modification_date <= ?", false, false, article.ModificationDate).
		Where("id != ?", article.ID).
		Order("modification_date DESC").
		Limit(1).
		Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// return no error
		return nil, errors.New("no prev article")
	}
	return &res, nil
}

func (a *ArticlesDb) Tags() ([]string, error) {
	var tags []string
	result := a.db.Model(&Article{}).Where("draft = ? AND special_page = ?", false, false).Joins("INNER JOIN tags a ON a.tag_id = articles.id").Pluck("DISTINCT a.name", &tags)
	if result.Error != nil {
		panic(result.Error)
	}
	return tags, nil
}

func (a *ArticlesDb) ArticlesByTag(tagName string) ([]Article, error) {
	var articles []Article
	result := a.db.Preload("Tags").Joins("INNER JOIN tags a ON a.tag_id = articles.id").Where("draft = ? AND special_page = ?", false, false).Where("a.name = ? COLLATE NOCASE", tagName).
		Order("modification_date DESC").Find(&articles)
	if result.Error != nil {
		return []Article{}, errors.New("no articles found")
	}
	return articles, nil
}

func (a *ArticlesDb) Series() ([]string, error) {
	var seriesList []string
	result := a.db.Model(&Article{}).Where("draft = ? AND special_page = ? AND series IS NOT NULL AND series != ''", false, false).Pluck("DISTINCT series", &seriesList)
	if result.Error != nil {
		panic(result.Error)
	}
	return seriesList, nil
}

func (a *ArticlesDb) ArticlesBySeries(series string) ([]Article, error) {
	var articles []Article
	result := a.db.Preload("Tags").Where("draft = ? AND special_page = ?", false, false).Where("series = ?  COLLATE NOCASE", series).
		Order("modification_date DESC").Find(&articles)
	if result.Error != nil {
		return []Article{}, errors.New("no articles found")
	}
	return articles, nil
}

func (a *ArticlesDb) NextArticleInSeries(article Article) (Article, error) {
	var res Article
	result := a.db.Preload("Tags").Where("draft = ? AND special_page = ? AND series = ? AND modification_date >= ?", false, false, article.Series, article.ModificationDate).
		Where("id != ?", article.ID).
		Order("modification_date DESC").
		Limit(1).
		Find(&res)
	if result.Error != nil {
		return Article{}, result.Error
	}
	if result.RowsAffected == 0 {
		return Article{}, errors.New("no next article in series found")
	}
	return res, nil
}

func (a *ArticlesDb) PrevArticleInSeries(article Article) (Article, error) {
	var res Article
	result := a.db.Preload("Tags").Where("draft = ? AND special_page = ? AND series = ? AND modification_date <= ?", false, false, article.Series, article.ModificationDate).
		Where("id != ?", article.ID).
		Order("modification_date DESC").
		Limit(1).
		Find(&res)
	if result.Error != nil {
		return Article{}, result.Error
	}
	if result.RowsAffected == 0 {
		return Article{}, errors.New("no prev article in series found")
	}
	return res, nil
}

func (a *ArticlesDb) Drafts() ([]Article, error) {
	var articles []Article
	result := a.db.Preload("Tags").Order("modification_date DESC").Where("draft = ?", true).Find(&articles)
	if result.Error != nil {
		panic(result.Error)
	}
	return articles, nil
}

func (a *ArticlesDb) SpecialPages() ([]Article, error) {
	var articles []Article
	result := a.db.Preload("Tags").Order("modification_date DESC").Where("special_page = ?", true).Find(&articles)
	if result.Error != nil {
		panic(result.Error)
	}
	return articles, nil
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
	assert.Equal(t, queryBySrcFileName.Tags[0].Name, "SteamDeck")
	assert.Equal(t, queryBySrcFileName.Tags[1].Name, "Gorilla")

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
	assert.Equal(t, taggedArticles[0].Tags[0].Name, "SteamDeck")

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
