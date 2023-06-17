package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	result := a.db.Where("src_file_name = ?", article.SrcFileName).FirstOrCreate(&article)
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

func (a *ArticlesDb) All() ([]Article, error) {
	var articles []Article
	result := a.db.Order("modification_date DESC").Find(&articles)
	if result.Error != nil {
		return []Article{}, nil
	}
	return articles, nil
}

func (a *ArticlesDb) Articles() ([]Article, error) {
	var articles []Article
	result := a.db.Order("modification_date DESC").Where("draft = ?", false).Where("special_page = ?", false).Find(&articles)
	if result.Error != nil {
		return []Article{}, nil
	}
	return articles, nil
}

func (a *ArticlesDb) MostRecentArticle() (Article, error) {
	var article Article
	result := a.db.Order("modification_date DESC").Where("draft = ?", false).Where("special_page = ?", false).First(&article)
	if result.Error != nil {
		return Article{}, nil
	}
	return article, nil
}

func (a *ArticlesDb) NextArticle() (Article, error) {
	//FIXME implement and test these
	var article Article
	result := a.db.Order("modification_date DESC").Where("draft = ?", false).Where("special_page = ?", false).First(&article)
	if result.Error != nil {
		return Article{}, result.Error
	}
	return article, nil
}

func (a *ArticlesDb) PrevArticle() (Article, error) {
	//FIXME implement and test these
	var article Article
	result := a.db.Order("modification_date DESC").Where("draft = ?", false).Where("special_page = ?", false).First(&article)
	if result.Error != nil {
		return Article{}, nil
	}
	return article, nil
}

//FIXME implement and test these
//func (a *ArticlesDb) NextArticleInSeries(series string) (Article, error) {}
//func (a *ArticlesDb) PrevArticleInSeries(series string) (Article, error) {}

func (a *ArticlesDb) Drafts() ([]Article, error) {
	var articles []Article
	result := a.db.Order("modification_date DESC").Where("draft = ?", true).Find(&articles)
	if result.Error != nil {
		panic(result.Error)
	}
	return articles, nil
}

func (a *ArticlesDb) SpecialPages() ([]Article, error) {
	var articles []Article
	result := a.db.Order("modification_date DESC").Where("special_page = ?", true).Find(&articles)
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
	article2 := Article{Title: "bar", ModificationDate: time2, Summary: "bar summary", Tags: []Tag{{Name: "SteamDeck"}, {Name: "Gorilla"}},
		SrcFileName: "/home/user/documents/bar.mdwn", DstFileName: "/home/user/documents/bar.html"}
	time3, _ := time.Parse(longForm, "2010-01-01 00:00")
	article3 := Article{Title: "batz", ModificationDate: time3, Summary: "batz summary", Tags: []Tag{{Name: "Linux"}, {Name: "Go"}},
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
	// try double insert
	err = articlesDb.Add(&article5)
	if err != nil {
		panic(err)
	}

	all, err := articlesDb.All()
	assert.True(t, err == nil)
	assert.Equal(t, len(all), 5)

	allarticles, err := articlesDb.All()
	assert.True(t, err == nil)
	assert.Equal(t, len(allarticles), 5)

	drafts, err := articlesDb.Drafts()
	assert.True(t, err == nil)
	assert.Equal(t, len(drafts), 1)

	specialpages, err := articlesDb.SpecialPages()
	assert.True(t, err == nil)
	assert.Equal(t, len(specialpages), 1)

	mostRecentArticle, err := articlesDb.MostRecentArticle()
	assert.True(t, err == nil)
	assert.Equal(t, mostRecentArticle.SrcFileName, "/home/user/documents/bar.mdwn")

	////// Find next/previous article ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// FIXME todo
	_, err = articlesDb.NextArticle()
	assert.Error(t, err)

	prevArticle, err := articlesDb.PrevArticle()
	assert.NotNil(t, err)
	assert.Equal(t, prevArticle.SrcFileName, "/home/user/documents/ba33r.mdwn")

	//nextArticle, err := articlesDb.NextArticle()
	//assert.Error(t, err)
	//assert.Equal(t, nextArticle.SrcFileName, "/home/user/documents/bar.mdwn")

	//// FIXME broken, []Tag is empty

	//// Query articles that are not drafts and sort by ModificationDate /////////////////////////////////////////////////////////////////////////////////////////
	//result = db.Where("draft = ?", true).Order("modification_date").Find(&articles)
	//if result.Error != nil {
	//	panic(result.Error)
	//}
	//
	//fmt.Println("draft Articles:")
	//for _, a := range articles {
	//	fmt.Printf("Date: %s, SrcFileName: %s, DstFileName: %s, Title: %s\n",
	//		a.ModificationDate, a.SrcFileName, a.DstFileName, a.Title)
	//}

	// update item ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//article4update := Article{Title: "done", ModificationDate: time4, Summary: "done summary", Tags: []Tag{{Name: "Linux"}, {Name: "Go"}},
	//	SrcFileName: "/home/user/documents/mydraft.mdwn", DstFileName: "/home/user/documents/mydraft.html"}
	//createOrUpdateArticle(db, article4update)

	// Query articles with the tag "go" //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//tagName := "SteamDeck"
	//db.Joins("INNER JOIN tags a ON a.tag_id = articles.id").Where("a.name = ? COLLATE NOCASE", tagName).Find(&articles)
	//if result.Error != nil {
	//	panic(result.Error)
	//}
	//
	//fmt.Printf("Articles with Tag '%s':\n", tagName)
	//for _, article := range articles {
	//	fmt.Printf("SrcFileName: %s, DstFileName: %s, Title: %s, Tags: %v\n",
	//		article.SrcFileName, article.DstFileName, article.Title, article.Tags)
	//}

	// delete item ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	err = articlesDb.Del("/home/user/documents/mydraft.mdwn")
	if err != nil {
		panic(err)
	}
	all2, err := articlesDb.All()
	assert.True(t, err == nil)
	assert.Equal(t, len(all2), 4)
}
