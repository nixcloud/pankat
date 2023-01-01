# what
pankat is a static blog/wiki generator inspired by joey hess's ikiwiki.

the primary motivation for rewriting ikiwiki was:
 - use pandoc as backend
 - mobile first by using bootstrap
 - more usability in navigation / posts overview
 - no interest in ikiwiki's perl

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

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

# where

a pankat generated blog can be found here:

* <https://lastlog.de/blog>

# license
pankat is licensed AGPL v3, see LICENSE for details.

# todo

* ArticlesCache: add error handling
  * http://blog.j7mbo.com/bypassing-golangs-lack-of-constructors/
  
* re-arrange directories: 
  * update the go version, use more recent nixpkgs
  * make it build with nix
  * move templates and general stuff into base dir
  * create hello world example so someone else can use this software

# content

Also the content of blog.lastlog.de needs some rework. Of course not related to the software but still listed here: 

* write a new article from time to time!? ;-)
* rework warning/info/danger/error ...
* write summary for each article
* rewrite title names
* fix images, add class="noFancy"
* check h1,h2,...
* use <div class="warn">...</div>
* check [[!series ogre]] for other series like qt
* libnoise_viewer.html fix video width

* commit history using git and add revert link like ikiwiki does FIXME
* implement comment system FIXME
   see example: https://www.reddit.com/r/golang/comments/1xbxzk/default_value_in_structs/

* BUG pandoc integration with parser '-s' of html head/body and migration to the go template

// FIXME create a [[!pandocFormat mdwn]] plugin which makes more pandoc dialects available

# advanced pankat editor

* gocraft/web ansprechen
* git backend ansprechen
* leaps backend ansprechen
* websockets preview mit long-polling
* lokales speichern von artikeln, wenn ./pankat -daemon -i documents -o output/ verwendet wird

* https://www.overleaf.com/4344023pmjpgq#/12921720/




# who

pankat is written and maintained by joachim schiele <js@lastlog.de>
