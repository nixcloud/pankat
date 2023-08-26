# what
**pankat** is a **static blog generator** inspired by [joey hess's ikiwiki](https://ikiwiki.info/users/joey/).

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

it is used at <https://lastlog.de/blog>.

# pankat-static

pankat-static can be used to render the articles into static html files.

* [pankat-static](cmd/pankat-static/README.md)

# pankat-server

pankat-server is used on the client to have life-preview during writing.

* [pankat-server](cmd/pankat-server/README.md)

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

# who

pankat is written and maintained by joachim schiele [js@lastlog.de](mailto:js@lastlog.de)



