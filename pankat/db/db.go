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
	ID                uint   `gorm:"primarykey"`
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
	ID    uint `gorm:"primarykey"`
	TagId uint
	Name  string
}

func (a Article) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		SrcFileName       string
		DstFileName       string
		ArticleMDWNSource []byte
		Title             string
		Summary           string
		Tags              []Tag
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
		Tags:              a.Tags,
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
	type Alias Article
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func (a Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name string
	}{
		Name: a.Name,
	})
}

func (a *Tag) UnmarshalJSON(data []byte) error {
	type Alias Tag
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

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
		// Auto-migrate the table
		err = db.AutoMigrate(&Article{}, &Tag{})
		if err != nil {
			panic(err)
		}
		dbInstance = &ArticlesDb{db: db}
	}
	return dbInstance
}

func mergeMaps(map1, map2 map[string]bool) map[string]bool {
	mergedMap := make(map[string]bool)

	// Add key-value pairs from map1
	for key, value := range map1 {
		mergedMap[key] = value
	}

	// Add key-value pairs from map2
	for key, value := range map2 {
		mergedMap[key] = value
	}

	return mergedMap
}

func getKeysFromMap(mergedMap map[string]bool) []string {
	keys := make([]string, 0, len(mergedMap))

	for key := range mergedMap {
		keys = append(keys, key)
	}

	return keys
}

func (a *ArticlesDb) GetRelatedArticles(article Article) map[string]bool {
	var affectedArticles = make(map[string]bool)
	if article.Draft || article.SpecialPage {
		return affectedArticles
	}
	next, err := a.NextArticle(article)
	if err == nil {
		affectedArticles[next.SrcFileName] = true
	}
	prev, err := a.PrevArticle(article)
	if err == nil {
		affectedArticles[prev.SrcFileName] = true
	}
	nextSeries, err := a.NextArticleInSeries(article)
	if err == nil {
		affectedArticles[nextSeries.SrcFileName] = true
	}
	prevSeries, err := a.PrevArticleInSeries(article)
	if err == nil {
		affectedArticles[prevSeries.SrcFileName] = true
	}
	return affectedArticles
}

func (a *ArticlesDb) Set(article *Article) (*Article, []string, error) {
	var affectedArticles = make(map[string]bool)
	newArticle := Article{}
	// a hack to update properties, since won't work with article
	update := make(map[string]interface{})
	update["title"] = article.Title
	update["src_file_name"] = article.SrcFileName
	update["dst_file_name"] = article.DstFileName
	update["article_mdwn_source"] = article.ArticleMDWNSource
	update["modification_date"] = article.ModificationDate
	update["summary"] = article.Summary
	update["series"] = article.Series
	update["special_page"] = article.SpecialPage
	update["draft"] = article.Draft
	update["anchorjs"] = article.Anchorjs
	update["tocify"] = article.Tocify
	update["timeline"] = article.Timeline
	update["show_source_link"] = article.ShowSourceLink
	update["live_updates"] = article.LiveUpdates
	update["evaluated"] = article.Evaluated

	err := a.db.Transaction(func(tx *gorm.DB) error {
		oldArticle := Article{}
		ret := tx.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).Limit(1).Find(&oldArticle)
		if ret.Error != nil && !errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return ret.Error
		}
		if ret.RowsAffected == 1 {
			affectedArticles = mergeMaps(affectedArticles, a.GetRelatedArticles(oldArticle))
			result := tx.Preload("Tags").Model(&Article{}).Where("src_file_name = ?", article.SrcFileName).Updates(update)
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
			updatedArticle := Article{}
			ret := tx.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).First(&updatedArticle)
			if ret.Error != nil {
				tx.Rollback()
				return ret.Error
			}
			if ret.RowsAffected != 1 {
				tx.Rollback()
				return errors.New("coudn't find article to delete tags for")
			}
			affectedArticles = mergeMaps(affectedArticles, a.GetRelatedArticles(updatedArticle))
			err := tx.Model(&updatedArticle).Unscoped().Association("Tags").Unscoped().Clear()
			if err != nil {
				tx.Rollback()
				return err
			}
			err = tx.Preload("Tags").Model(&updatedArticle).Association("Tags").Replace(article.Tags)
			if err != nil {
				tx.Rollback()
				return err
			}
			ret = tx.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).First(&newArticle)
			if ret.Error != nil {
				return ret.Error
			}
			return nil
		}
		// create new article
		ret = tx.Create(article)
		if ret.Error != nil {
			tx.Rollback()
			return ret.Error
		}
		if ret.RowsAffected != 1 {
			tx.Rollback()
			return errors.New("couldn't create article")
		}
		ret = tx.Preload("Tags").Where("src_file_name = ?", article.SrcFileName).First(&newArticle)
		if ret.Error != nil {
			return ret.Error
		}
		affectedArticles = mergeMaps(affectedArticles, a.GetRelatedArticles(newArticle))
		return nil
	})
	if err != nil {
		return nil, []string{}, err
	}
	return &newArticle, getKeysFromMap(affectedArticles), nil
}

func (a *ArticlesDb) Del(SrcFileName string) ([]string, error) {
	var affectedArticles = make(map[string]bool)
	article := Article{}
	ret := a.db.Preload("Tags").Where("src_file_name = ?", SrcFileName).First(&article)
	if ret.Error != nil {
		return []string{}, ret.Error
	}
	if ret.RowsAffected == 0 || ret.RowsAffected > 1 {
		return []string{}, errors.New("Coudn't find article to delete")
	}
	affectedArticles = mergeMaps(affectedArticles, a.GetRelatedArticles(article))
	err := a.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&article).Unscoped().Association("Tags").Unscoped().Clear()
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Unscoped().Delete(&article).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return getKeysFromMap(affectedArticles), nil
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
		return nil, result.Error
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
		return nil, result.Error
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
		return nil, result.Error
	}
	return articles, nil
}

func (a *ArticlesDb) SpecialPages() ([]Article, error) {
	var articles []Article
	result := a.db.Preload("Tags").Order("modification_date DESC").Where("special_page = ?", true).Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	return articles, nil
}
