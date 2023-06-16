package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const longForm = "2006-01-02 15:04"

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

// https://gorm.io/docs/create.html
func (a *ArticlesDb) Add(article *Article) error {
	result := a.db.Create(&article)
	if result.Error != nil {
		return result.Error
	}
	//FIXME needs checking
	//article.ID          // returns inserted data's primary key
	//result.Error        // returns error
	//result.RowsAffected // returns inserted records count

	//result := db.Where("SrcFileName = ?", article.SrcFileName).Updates(&article)
	//
	//if result.Error != nil {
	//	return result.Error
	//}
	return nil
}

func (a *ArticlesDb) Del(SrcFileName string) error {
	//FIXME needs checking

	result := a.db.Where("srcfilename = ?", SrcFileName).Delete(&Article{})
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected != 1 {
		panic("no rows affected")
	}
	return nil
}

func TestArticlesDatabase(t *testing.T) {

	articlesDb := NewArticlesDb()

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

	//// Query the table to verify insertion
	var articles []Article
	result := articlesDb.db.Order("modification_date DESC").Find(&articles)
	if result.Error != nil {
		panic(result.Error)
	}

	assert.True(t, result.Error == nil)

	////// Query all articles in the db //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//fmt.Println("All Articles in DB:")
	//for _, a := range articles {
	//	fmt.Printf("Date: %s, SrcFileName: %s, DstFileName: %s, Title: %s, Tags %v\n",
	//		a.ModificationDate, a.SrcFileName, a.DstFileName, a.Title, a.Tags)
	//}
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

	////// Find next/previous article ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// FIXME todo

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
}
