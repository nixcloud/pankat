# what
**pankat** is a **static blog generator** inspired by [joey hess's ikiwiki](https://ikiwiki.info/users/joey/).

both pankat-static and pankat-server build a static blog.

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

## example

a pankat generated blog can be found here:

* <https://lastlog.de/blog>

but you can easily use this software for your own blog as well!

# pankat-static

pankat-static can be used to render the articles.

* [pankat-static](src/github.com/nixcloud/cmd/pankat-static/README.md)

# pankat-server

pankat-server can be executed on the client to provides a live preview of the generated static blog, later the result can be uploaded.

* [pankat-server](src/github.com/nixcloud/cmd/pankat-server/README.md)

# license
pankat is licensed AGPL v3, see [LICENSE](LICENSE) for details.

# history

the primary motivation for rewriting ikiwiki was:
- use pandoc as backend
- mobile first by using bootstrap
- more usability in navigation / posts overview
- no interest in ikiwiki's perl

# todo

## programming

* BUG: pandoc integration with parser '-s' of html head/body and migration to the go template

* extend MetaData for articles: must know the source file it was generated from to listen for updates, the md5 hash, ...

* make evaluation lazy, rework md5 and article re-creation, rework rsync, rework xml, rework timeline
* // when to rerender article? Articles.go
  // change in
  // - title
  // - Article
  // - ModificationDate
  // - tags
  // - series
  // - SrcFileName
  // - DstFileName
  // - BaseFileName
  // - SrcDirectoryName
  // - Anchorjs         
  // - Tocify           
  //
  // but also when
  // previous/next article have these changes
  // - DstFileName
  * ArticlesCache: add error handling

* create list of articles in draft state 

* use Article.timeline field, don't call RenderPosts twice

     articlesTopLevel := articlesAll.TopLevel().FilterOutDrafts()
     articlesPosts := articlesAll.Posts().FilterOutDrafts()
     RenderPosts(articlesTopLevel)
     RenderPosts(articlesPosts)

* rework fsnotify to know exactly which file was changed
  * https://medium.com/@skdomino/watch-this-file-watching-in-go-5b5a247cf71f
 
* add drafts subpage which lists all articles in draft state 
* add pub/sub system for websocket where clients can register a certain page with a hash

* add interface to inject updated documents notifications
  * `func RenderPosts(articles Articles) { ... fmt.Println("NOTIFICATION: ", e.DstFileName)` 

* https://github.com/nixcloud/pankat/issues/3
  a href="https://lastlog.de/blog/media/tuex.png
  should be
  a href="https://lastlog.de/blog/posts/media/tuex.png
* feed.xml is not generated anymore, might have worked in f5e3232f1df691f3a3b21ca54b77c2b13a9db564
* re-arrange directories: 
  * move templates and general stuff into base dir
* generalize rsync (windows and linux)
* integrate leaps editor
* BUG: timeline needs beginning year and end year, some have not both
later

* create hello world example so someone else can use this software
* implement comment system FIXME
   see example: https://www.reddit.com/r/golang/comments/1xbxzk/default_value_in_structs/
* FIXME create a [[!pandocFormat mdwn]] plugin which makes more pandoc dialects available
* FIX bug with regexp where grep_and_vim_idea.mdwn contains plugin calls which should not be rendered
   func processPlugins(_article []byte, article *Article) []byte {
  var _articlePostprocessed []byte

  re := regexp.MustCompile("\\[\\[!(.*?)\\]\\]")

## content

the content of blog.lastlog.de needs rework 

* write a new article from time to time!? ;-)
* rework warning/info/danger/error ...
* write summary for each article; write missing summary plugin (chatGPT?)
* rewrite title names
* fix images, add class="noFancy"
* check h1,h2,...
* use <div class="warn">...</div>
* libnoise_viewer.html fix video width

# who

pankat is written and maintained by joachim schiele <js@lastlog.de>
