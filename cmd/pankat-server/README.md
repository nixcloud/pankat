# pankat-server

pankat-server is a static blog/wiki generator inspired by joey hess's ikiwiki, 
but it features live preview during article writing. you can use any editor you like and once you hit save the browser will update the article without a browser page reload.

# how to compile

    git clone https://github.com/nixcloud/pankat

compile the pankat-server binary:

    cd  cmd/pankat-server
    go build pankat-server

# how to use

the **pankat-server** will generate all **static pages**. once this is done it will open a
socket on <http://localhost:8000> which needs to be visited with a browser. there one can experience
the blog as if it were deployed on a webserver. 

the main difference being:

* it offers a live preview of pages edits using js/websocket
* a list of draft articles

start it like this:

1. start the pankat-server binary

       pankat-server --documents documents/blog.lastlog.de/ 

2. afterward open this in firefox/chromium 

       localhost:8000

3. using any editor like notepad/vim/... open the mdwn documents with any editor and one a change is saved it will update the article in the browser without reload

## draft editing

a draft is a normal mdwn document which contains **[[!draft]]** in the source code.

see them on <localhost:8000/draft>

these rules apply for a draft:

* it will not be listed in the timeline 
* no html document will be generated, only live preview
* the source code (mdwn) probably gets published, so be aware of that

## special pages

SpecialPages like about, roadmap, websocket should:

* should contain left, right navigational keyboard commands
* should not contain a navigation bar
* not appear in the timeline
* however, unlike drafts, they should render as normal page

Use [[!SpecialPage]] to mark a page as special page.

## file/filename changes

if a file, take about.mdwn as an example, is

* moved to a different folder, you need to restart pankat-server as watcher can't detect this on windows
* renamed, you need to change to the new url in the browser manually 

## markdown syntax

pankat uses pandoc markdown syntax but uses various plugins to extend it:

  * [[!specialpage]]
  * [[!draft]]
  * [[!meta]]
  * [[!series]]
  * [[!tag]]
  * [[!img]]
  * [[!summary]]
  * [[!title]]
  * ... more, see implementation & lastlog.de/blog source code for examples

you can also use inline html/css/javascript in your markdown files.

### license

see [LICENSE](../../LICENSE)
