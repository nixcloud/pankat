# what
pankat is a static blog/wiki generator inspired by joey hess's ikiwiki.

# how to compile

    git clone https://github.com/nixcloud/pankat

afterwards use go to compile the binary

    pankat\src\github.com\nixcloud\pankat
    go build pankat

# how to use

once pankat has been compiled, use blog.lastlog.de to see how to use the software:

    pankat --input documents/blog.lastlog.de/ --output output/blog.lastlog.de

run a local webserver in output/blog.lastlog.de to review the static page generator output

    python -m http.server

then visit localhost:8000 in any webbrowser.