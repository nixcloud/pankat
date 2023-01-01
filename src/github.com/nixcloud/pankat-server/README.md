# prototype code for using websockets with pandoc

the server should do this:

* generate a complete pandoc document into html
* send it to the client
* client uses html to display it on the right pane
* each time an update comes in, which is basically a complete document each time, the diff is built between original and update and all differences are background-colored to yellow
* ... until someone hits 'save' which creates a git commit and publishes the article

an exemplary implementation is in ../../pandoc-online-editor/ at qknights computer

consider vue with server side rendering:

* https://vuejs.org/v2/examples/hackernews.html
* https://medium.com/@adeshg7/vuejs-golang-a-rare-combination-53538b6fb918

## websocket-pandoc

websocket-pandoc is a experimental implementation of a websocket / virtual dom based updater which basically supports server side document generating with changes highlighting.

### try it

simply run:

    nix-shell

then:

    go run main.go

and afterwards do:

    chromium localhost:8080

and enjoy the editor!


### screenshots

see screenshots/ folder

### license

* user interface was inspired by: https://notehub.org
* websocket example code based on: https://github.com/golang-samples/websocket
