# what
pankat is a static blog/wiki generator inspired by ikiwiki from joey hess.

the primary motivation for rewriting ikiwiki was:
 - use pandoc as backend
 - use bootstrap (mobile first) in the frontend
 - more usability in navigation / posts overview

![A screenshot featuring pankat](https://raw.githubusercontent.com/nixcloud/pankat/master/screenshots/pankat.jpg)

# where

a pankat generated blog can be found here:

* <https://lastlog.de/blog>

# license
pankat is licensed AGPL v3, see LICENSE for details.

# todo

* feed: //BUG: feed needs ./posts/media/ URLs instead of ./media/ URLs
* feed: fix BUG <link rel="alternate" type="application/atom+xml" title="lastlog.de/blog Atom feed" href="{{.RelativeSrcRootPath}}/feed.xml" />
* feed: check feed quality, for instance pictures
* ArticlesCache: add error handling
  * http://blog.j7mbo.com/bypassing-golangs-lack-of-constructors/
* move templates and general stuff into base dir

* FIXME next/last hover shadow
* FIXME use h1 only for title, see http://pandoc.org/scripting.html filter
* FIXME add donation button
* SECURITY secure pandoc from passing < script>alert('hi')</script> and other evil <html tags>         find a filter system for evil html tags like script

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
