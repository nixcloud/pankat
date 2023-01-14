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

# pankat-docker

Build programs:

    docker build -t pankat-docker:latest .

Use program `pankat-static`:

    docker run -it --rm -v ${PWD}/documents:/documents -v ${PWD}/output:/output pankat-docker:latest pankat-static --documents /documents/blog.lastlog.de/

Use program `pankat-server`:

    docker run -it --rm -p 8000:8000 -v ${PWD}/documents:/documents pankat-docker:latest pankat-server --documents /documents/blog.lastlog.de/

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

* consider using nix-instantiate instead own md5 hash system
* merge images and posts/media folder

* BUG: timeline needs beginning year and end year, some don't, so they are not rendered

* tidy generated html code
  * FIX gohtml.FormatBytes() is not working properly, css needs fixes for posts.html, did not work for normal pages, duno why - but won't do this ATM

* live preview
  * lacks tags, creation date and title
  * live updates of TOC not working: $("#toc").tocify(); called twice does nothing.
  * update cmd/pankat-server/ws/server.go to use pub/sub system for websocket where clients can register a certain page; if registered page is changed on the source side live updates are sent
  * use v8 https://github.com/rogchap/v8go with serverside diffDOM.js to mainly send diffs to the client
  * integrate leaps editor

* make evaluation lazy, rework md5 and article re-creation, rework xml, rework timeline
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

* scroll up: some articles render it inside the paper, some outside
 
* mobile html is the horror:
  * font size bogus, 
  * background paper while it would look better fullscreen
  * < and > for series is not expanding to vertical size
  * main menu is hard to read
  * source code boxes are tiny compared to the rest
* add drafts subpage which lists all articles in draft state 

* ArticlesCache: add error handling

* BUG https://github.com/nixcloud/pankat/issues/3
  a href="https://lastlog.de/blog/media/tuex.png
  should be
  a href="https://lastlog.de/blog/posts/media/tuex.png
* FIX bug with regexp where grep_and_vim_idea.mdwn contains plugin calls which should not be rendered
   func processPlugins(_article []byte, article *Article) []byte {
  var _articlePostprocessed []byte

  re := regexp.MustCompile("\\[\\[!(.*?)\\]\\]")

* create hello world example so someone else can use this software
* implement comment system FIXME
  see example: https://www.reddit.com/r/golang/comments/1xbxzk/default_value_in_structs/
* FIXME create a [[!pandocFormat mdwn]] plugin which makes more pandoc dialects available
* write a links plugin
* remove media/* and HTML documents where there is no reference to anymore (GC)

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
