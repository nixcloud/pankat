# pankat-server

the **pankat-server** will generate **static pages** and once done will wait for further modifications or new articles.

    pankat-server --input documents/blog.lastlog.de/ --output output/blog.lastlog.de

afterwards open this in firefox/chromium 

    localhost:8000

planned: articles will have an `edit button` somewhere and if clicked the article will open using [leaps](https://github.com/Jeffail/leaps) on the left and a preview render on the right.

**note: this is still WIP**

what does not work yet

* edit button
* leaps editor integration

what works:

* WS reload command on documents/blog.lastlog.de updates
* localhost:8000 shows the documents

## virtual dom ideas

* http://fiduswriter.github.io/diffDOM/

* https://vuejs.org/v2/examples/hackernews.html
* https://medium.com/@adeshg7/vuejs-golang-a-rare-combination-53538b6fb918

## websockets in go

a few examples how to use

* https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/08.2.html
* https://github.com/golang-samples/websocket/blob/master/websocket-chat/src/chat/client.go


## websocket-pandoc

websocket-pandoc is an experimental implementation of a websocket / virtual dom based updater which basically supports server side document generating with changes highlighting.

### license

see LICENSE file, AGPL v3

* user interface was inspired by: https://notehub.org
* websocket example code based on: https://github.com/golang-samples/websocket
