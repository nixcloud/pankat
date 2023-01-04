# what
**pankat** is a **static blog generator** inspired by [joey hess's ikiwiki](https://ikiwiki.info/users/joey/).

both pankat-static and pankat-server build a static blog.

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

## example

a pankat generated blog can be found here:

* <https://lastlog.de/blog>

but you can easily use this software for your own blog as well!

# pankat-static

pankat-static can be used to render the articles.

* [pankat-static](src/github.com/nixcloud/cmd/pankat-static/README.md)

# pankat-server

pankat-server can be executed on the client to provides a live preview of the generated static blog, later the result can be uploaded.

* [pankat-server](src/github.com/nixcloud/cmd/pankat-server/README.md)

# license
pankat is licensed AGPL v3, see [LICENSE](LICENSE) for details.

# history

the primary motivation for rewriting ikiwiki was:
- use pandoc as backend
- mobile first by using bootstrap
- more usability in navigation / posts overview
- no interest in ikiwiki's perl

# todo

## programming

* make everything lazy, rework md5 and article re-creation, rework rsync, rework fsnotify checks on all files 
* add channel to inject updated documents notifications
  * asdf 
* each article must know the source file it was generated from to listen for updates

* REST

  * /tags
  * /series

* ArticlesCache: add error handling
* feed.xml is not generated anymore, might have worked in f5e3232f1df691f3a3b21ca54b77c2b13a9db564
* re-arrange directories: 
* move templates and general stuff into base dir
* fix rsync (windows and linux)
* integrate leaps editor
* BUG: pandoc integration with parser '-s' of html head/body and migration to the go template

later

* create hello world example so someone else can use this software
* implement comment system FIXME
   see example: https://www.reddit.com/r/golang/comments/1xbxzk/default_value_in_structs/

## content

the content of blog.lastlog.de needs rework 

* write a new article from time to time!? ;-)
* rework warning/info/danger/error ...
* write summary for each article; write missing summary plugin (chatGPT?)
* rewrite title names
* fix images, add class="noFancy"
* check h1,h2,...
* use <div class="warn">...</div>
* check [[!series ogre]] for other series like qt
* libnoise_viewer.html fix video width

* commit history using git and add revert link like ikiwiki does FIXME

// FIXME create a [[!pandocFormat mdwn]] plugin which makes more pandoc dialects available

# who

pankat is written and maintained by joachim schiele <js@lastlog.de>
