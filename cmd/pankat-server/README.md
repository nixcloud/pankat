# pankat-server

pankat-server is a static blog/wiki generator inspired by joey hess's ikiwiki but it features live-updates.

# how to compile

    git clone https://github.com/nixcloud/pankat

afterwards use go to compile the binary

    cd  cmd/pankat-server
    go build pankat-server

# how to use

the **pankat-server** will first generate **static pages**. once this is done it will open a
socket on localhost:8000 which needs to be visited with a browser. there one can experience
the blog as if it were deployed on a webserver but the main difference is, that it offers a live
preview of the pages, draft editing and a roadmap.

1. start the server

       pankat-server --documents documents/blog.lastlog.de/ 

2. afterwards open this in firefox/chromium 

       localhost:8000

3. open the mdwn documents with any editor and one a change is saved it will update the article in the browser without reload

## draft editing

drafts are not published. they are only visible in the browser when the pankat-server is running on localhost:8000/draft

that is, if a document contains a [[!draft]] reference:

* it will not be rendered into the timeline (/posts) 
* no html document will be generated thus published
* the source code (mdwn) probably gets published, so be aware of that

Use [[!draft]] to mark a page as draft.

## special pages

SpecialPages like about, roadmap, websocket should:

* should contain left, right navigational keyboard commands
* should not contain a navigation bar
* not appear in the timeline
* however, unlike drafts, they should render as normal page

Use [[!SpecialPage]] to mark a page as special page.

## file/filename changes

if a file about.mdwn is

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

see LICENSE file, AGPL v3
