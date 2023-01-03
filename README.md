# what
pankat is a static blog/wiki generator inspired by joey hess's ikiwiki.

read on in this documentation

* [pankat](src/github.com/nixcloud/cmd/pankat/README.md)
* [pankat-server](src/github.com/nixcloud/cmd/pankat-server/README.md)

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

# history

the primary motivation for rewriting ikiwiki was:
 - use pandoc as backend
 - mobile first by using bootstrap
 - more usability in navigation / posts overview
 - no interest in ikiwiki's perl

# where

a pankat generated blog can be found here:

* <https://lastlog.de/blog>

but you can easily use this software for your own blog as well!

# license
pankat is licensed AGPL v3, see LICENSE for details.

# todo

## programming

* split pankat into pankat-core, pankat and pankat-server
* rework md5 and article re-creation; add channel to communicate updated documents
* each article must know the source file it was generated from to listen for updates

* REST

  * /tags
  * /series

* ArticlesCache: add error handling

  
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

# advanced pankat editor

* gocraft/web ansprechen
* git backend ansprechen
* leaps backend ansprechen
* websockets preview mit long-polling
* lokales speichern von artikeln, wenn ./pankat -daemon -i documents -o output/ verwendet wird

* https://www.overleaf.com/4344023pmjpgq#/12921720/




# who

pankat is written and maintained by joachim schiele <js@lastlog.de>
