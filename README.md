# obsolescence notice

this **go based implementation of pankat is obsoleted by the rust rewarite** at https://github.com/nixcloud/pankat-rs

# what
**pankat** is a **static blog generator** inspired by [joey hess's ikiwiki](https://ikiwiki.info/users/joey/).

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

it is used at <https://lastlog.de/blog>.

if you want to play with it, you can use my content in documents/blog.lastlog.de:

     git submodule update --init --recursive

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

# pankat external dependencies

you need to install https://pandoc.org/installing.html

## windows

* https://pandoc.org/installing.html - pandoc-3.x
* https://jmeubank.github.io/tdm-gcc/
* go 1.19 (using goland)

# license
pankat is licensed AGPL v3, see [LICENSE](LICENSE) for details.

# who

pankat is written and maintained by joachim schiele [js@lastlog.de](mailto:js@lastlog.de)



