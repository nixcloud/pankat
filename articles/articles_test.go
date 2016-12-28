package pankat

import (
  "testing"
  "fmt"
  "sort"
  "time"
)

func init() {
  var articles meerk.Articles
  
  //////// article 1 ////////////////////
  article1 := new( meerk.Article)
  article1.Title = "my fancy article"
  article1.ModificationDate = time.Now()
  article1.SrcDirectoryName = "posts/my_fancy_article/"
  article1.Series = "myPony"
  article1.Tags = []string { "linux", "NixOS" }
  
  articles = append(articles, article1)
  
  //////// article 2 ////////////////////
  article2 := new( meerk.Article)
  article2.Title = "a new article deluxe"
  const longForm = "2006-01-02 15:04"
  t, _ := time.Parse(longForm, "2012-01-02 15:04")
  article2.ModificationDate = t
  article2.SrcDirectoryName = "posts/"
  article2.Tags = []string { "linux" }
  
  articles = append(articles, article2)
  
  //////// article 3 ////////////////////
  
  article3 := new( meerk.Article)
  article3.Title = "recent article"
  t1, _ := time.Parse(longForm, "2022-05-12 15:04")
  article3.ModificationDate = t1
  article3.Draft = true
  article3.Tags = []string { "linux" }
  
  article3.SrcDirectoryName = "posts/myseries/"
  article3.Series = "myPony"
  articles = append(articles, article3)
  
  ///////////////////////////////////////
}

func TestSortFunc() {
  
  // 2. sort them by date
  sort.Sort(meerk.Articles(articles))
  
}
  
  fmt.Println("- MakeRelativeLink test----------------------------")  
  fmt.Println("article1.SrcDirectoryName: ", article1.SrcDirectoryName)
  fmt.Println("article2.SrcDirectoryName: ", article2.SrcDirectoryName)
  fmt.Println(articles.MakeRelativeLink(article1, article2))
  fmt.Println(articles.MakeRelativeLink(article2, article1))
  fmt.Println("- /MakeRelativeLink test----------------------------")  
  
  z := articles.NextArticle(article1)
  fmt.Println("-next test----------------------------")
  fmt.Println(article1.Title)
  if z != nil { fmt.Println(z.Title) }
  fmt.Println("-/next test----------------------------")
  
  z1 := articles.PrevArticle(article1)
  fmt.Println("-previous test----------------------------")
  fmt.Println(article1.Title)
  if z1 != nil {fmt.Println(z1.Title)}
  fmt.Println("-/previous test----------------------------")

  z2 := articles.PrevArticleInSeries(article1)
  fmt.Println("-previousSeries test----------------------------")
  fmt.Println(article1.Title)
  if z2 != nil {fmt.Println(z2.Title)}
  fmt.Println("-/previousSeries test----------------------------")
  
  z3 := articles.NextArticleInSeries(article1)
  fmt.Println("-nextSeries test----------------------------")
  fmt.Println(article1.Title)
  if z3 != nil {fmt.Println(z3.Title)}
  fmt.Println("-/nextSeries test----------------------------")
  
  fmt.Println("=====1==1=1=1=1= articles without filter ====1=1====1==1=1===")
for i, e := range articles {
//     _ = e
    fmt.Println(i, ": ", e.Title)
    fmt.Println("     ",e.ModificationDate)
    fmt.Println("     ",e.SrcDirectoryName)
    fmt.Println("     ",e.Draft)
    fmt.Println("     ",e.Series)
    fmt.Println("     ",e.Tags)
  }

fmt.Println("~~~~~~~~~~~~~~~~~~~ FilterByDraft() ~~~~~~~~~~~~~~~~")
for i, e := range articles.FilterByDraft() {
//     _ = e
    fmt.Println(i, ": ", e.Title)
    fmt.Println("     ",e.ModificationDate)
    fmt.Println("     ",e.SrcDirectoryName)
    fmt.Println("     ",e.Draft)
    fmt.Println("     ",e.Series)    
    fmt.Println("     ",e.Tags)
  }
  
fmt.Println("~~~~~~~~~~~~~~~~~~~ .FilterByTag(\"linux\") ~~~~~~~~~~~~~~~~")
for i, e := range articles.FilterByTag("NixOS") {
//     _ = e
    fmt.Println(i, ": ", e.Title)
    fmt.Println("     ",e.ModificationDate)
    fmt.Println("     ",e.SrcDirectoryName)
    fmt.Println("     ",e.Draft)
    fmt.Println("     ",e.Series)    
    fmt.Println("     ",e.Tags)
  }

fmt.Println("~~~~~~~~~~~~~~~~~~~ .FilterBySeries(\"myPony\") ~~~~~~~~~~~~~~~~")
for i, e := range articles.FilterBySeries("myPony") {
//     _ = e
    fmt.Println(i, ": ", e.Title)
    fmt.Println("     ",e.ModificationDate)
    fmt.Println("     ",e.SrcDirectoryName)
    fmt.Println("     ",e.Draft)
    fmt.Println("     ",e.Series)    
    fmt.Println("     ",e.Tags)
  }
}






