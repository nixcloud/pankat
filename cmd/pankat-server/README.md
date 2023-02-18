# pankat-server

the **pankat-server** will first generate **static pages**. once this is done it will open a
socket on localhost:8000 which needs to be visited with a browser. there one can experience
the blog as if it were deployed on a webserver but the main difference is, that it offers a live
preview of the pages.

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

### license

see LICENSE file, AGPL v3
