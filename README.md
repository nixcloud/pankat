# what
**pankat** is a **blog generator** inspired by [joey hess's ikiwiki](https://ikiwiki.info/users/joey/).

pankat-static and pankat-server build a static blog.

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

## example

a pankat generated blog can be found here:

* <https://lastlog.de/blog>

but you can easily use pankat for your own blog!

# pankat-static

pankat-static can be used to render the articles.

* [pankat-static](src/github.com/nixcloud/cmd/pankat-static/README.md)

# pankat-server

pankat-server can be executed on the client to provides a live preview of the generated static blog, later the result can be uploaded.

* [pankat-server](src/github.com/nixcloud/cmd/pankat-server/README.md)

it is a helpful tool to modify drafts before publishing them.

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
- no interest in ikiwiki's perl language choice
- use pandoc as backend for rendering articles
- mobile first by using bootstrap
- more usability in navigation / posts overview

# who

pankat is written and maintained by joachim schiele <js@lastlog.de>



