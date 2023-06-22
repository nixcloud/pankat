package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Article struct {
	gorm.Model
	SrcFileName       string `gorm:"uniqueIndex"`
	DstFileName       string
	ArticleMDWNSource []byte
	Title             string
	ModificationDate  time.Time
	Summary           string
	Tags              []Tag `gorm:"ForeignKey:TagId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Series            string
	SpecialPage       bool
	Draft             bool
	Anchorjs          bool
	Tocify            bool
	Timeline          bool
	ShowSourceLink    bool
	LiveUpdates       bool
	Evaluated         bool
}

type Tag struct {
	gorm.Model
	TagId uint
	Name  string
}

func toTags(tags []string) []Tag {
	var t []Tag
	for _, tag := range tags {
		t = append(t, Tag{Name: tag})
	}
	return t
}

func fromTags(tags []Tag) []string {
	var t []string
	for _, tag := range tags {
		t = append(t, tag.Name)
	}
	return t
}

func (a Article) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		SrcFileName       string
		DstFileName       string
		ArticleMDWNSource []byte
		Title             string
		Summary           string
		Tags              []string
		Series            string
		SpecialPage       bool
		Draft             bool
		Anchorjs          bool
		Tocify            bool
		Timeline          bool
		ShowSourceLink    bool
		LiveUpdates       bool
		Evaluated         bool
		ModificationDate  string
	}{
		SrcFileName:       a.SrcFileName,
		DstFileName:       a.DstFileName,
		ArticleMDWNSource: a.ArticleMDWNSource,
		Title:             a.Title,
		Summary:           a.Summary,
		Tags:              fromTags(a.Tags),
		Series:            a.Series,
		SpecialPage:       a.SpecialPage,
		Draft:             a.Draft,
		Anchorjs:          a.Anchorjs,
		Tocify:            a.Tocify,
		Timeline:          a.Timeline,
		ShowSourceLink:    a.ShowSourceLink,
		LiveUpdates:       a.LiveUpdates,
		Evaluated:         a.Evaluated,
		ModificationDate:  a.ModificationDate.Format(time.RFC3339),
	})
}

func (a *Article) UnmarshalJSON(data []byte) error {
	aux := &struct {
		SrcFileName       string
		DstFileName       string
		ArticleMDWNSource []byte
		Title             string
		Summary           string
		Tags              []string
		Series            string
		SpecialPage       bool
		Draft             bool
		Anchorjs          bool
		Tocify            bool
		Timeline          bool
		ShowSourceLink    bool
		LiveUpdates       bool
		Evaluated         bool
		ModificationDate  string
	}{
		SrcFileName:       a.SrcFileName,
		DstFileName:       a.DstFileName,
		ArticleMDWNSource: a.ArticleMDWNSource,
		Title:             a.Title,
		Summary:           a.Summary,
		Tags:              fromTags(a.Tags),
		Series:            a.Series,
		SpecialPage:       a.SpecialPage,
		Draft:             a.Draft,
		Anchorjs:          a.Anchorjs,
		Tocify:            a.Tocify,
		Timeline:          a.Timeline,
		ShowSourceLink:    a.ShowSourceLink,
		LiveUpdates:       a.LiveUpdates,
		Evaluated:         a.Evaluated,
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	a.ModificationDate, err = time.Parse(time.RFC3339, aux.ModificationDate)
	return err
}

//func (a Article) MarshalJSON() ([]byte, error) {
//	type Alias Article
//	return json.Marshal(&struct {
//		Alias
//		ModificationDate string
//	}{
//		Alias:            (Alias)(a),
//		ModificationDate: a.ModificationDate.Format(time.RFC3339),
//	})
//}
//
//func (a *Article) UnmarshalJSON(data []byte) error {
//	type Alias Article
//	aux := &struct {
//		*Alias
//		ModificationDate string
//	}{
//		Alias: (*Alias)(a),
//	}
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//	var err error
//	a.ModificationDate, err = time.Parse(time.RFC3339, aux.ModificationDate)
//	return err
//}

var lock = &sync.Mutex{}

type ArticlesDb struct {
	db *gorm.DB
}

var dbInstance *ArticlesDb

func Instance() *ArticlesDb {
	lock.Lock()
	defer lock.Unlock()
	if dbInstance == nil {
		fmt.Println("Creating DB instance now.")
		// Open a new SQLite database connection or create one if it doesn't exist
		db, err := gorm.Open(sqlite.Open("pankat-sqlite3.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		// remove old entries
		//db.Migrator().DropTable(&Article{}, &Tag{})

		// Auto-migrate the table
		err = db.AutoMigrate(&Article{}, &Tag{})
		if err != nil {
			panic(err)
		}
		dbInstance = &ArticlesDb{db: db}
	}
	return dbInstance
}

//func (a *ArticlesDb) Add(article *Article) error {
//	result := a.db.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).FirstOrCreate(&article)
//	if result.Error != nil {
//		return result.Error
//	}
//	if result.RowsAffected == 0 {
//		// return no error
//		fmt.Println("article was not updated")
//	}
//
//	return nil
//}

func (a *ArticlesDb) Add(article *Article) error {
	var existingArticle Article
	result := a.db.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).First(&existingArticle)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = a.db.Create(article)
	} else {
		a.Del(article.SrcFileName)
		result = a.db.Save(&existingArticle)
	}

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		fmt.Println("article was not updated")
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
		return Article{}, result.Error
	}
	return res, nil
}

func (a *ArticlesDb) NextArticle(article Article) (*Article, error) {
	var res Article
	result := a.db.Preload("Tags").Where("draft = ? AND special_page = ? AND modification_date >= ?", false, false, article.ModificationDate).
		Where("id != ?", article.ID).
		Order("modification_date").
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

func (a *ArticlesDb) AllTagsInDB() ([]string, error) {
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

func (a *ArticlesDb) AllSeriesInDB() ([]string, error) {
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
		Order("modification_date").
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
