# pankat-static

pankat-static is a static blog/wiki generator inspired by joey hess's ikiwiki.

# how to compile

    git clone https://github.com/nixcloud/pankat

compile the pankat-static binary:

    cd cmd/pankat-static
    go build pankat-static

# how to use

once `pankat-static` has been compiled, use blog.lastlog.de to see how to use the software:

    pankat-static --documents documents/blog.lastlog.de/

run a local webserver in documents/blog.lastlog.de to review the static page generator output

    cd documents/blog.lastlog.de
    python -m http.server

then visit <http://localhost:8000> in any browser.

### license

see [LICENSE](../../LICENSE)