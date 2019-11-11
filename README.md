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
* rework hashing: current hasing forgets to update former last article to have a 'next article' link if new article is added
* generate rss/atom feed
* create feed per tag
*  fix history writing
   example: 1. go to article https://lastlog.de/blog/posts/tour_of_nix.html
            2. click on an article tag https://lastlog.de/blog/posts.html?tag=emscripten
            3. then try 'back' button, which fails!
      maybe use backbone.js for that?

* FIXME donation button
* FIXME next/last hover shadow

* FIXME use h1 only for title, see http://pandoc.org/scripting.html filter



* SECURITY secure pandoc from passing < script>alert('hi')</script> and other evil <html tags>         find a filter system for evil html tags like script

*/ https://www.overleaf.com/4344023pmjpgq#/12921720/

* gocraft/web ansprechen
* git backend ansprechen
* leaps backend ansprechen
* websockets preview mit long-polling
* lokales speichern von artikeln, wenn ./pankat -daemon -i documents -o output/ verwendet wird



# content

Also the content of blog.lastlog.de needs some rework. Of course not related to the software but still listed here: 

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



# who
pankat is written and maintained by joachim schiele <js@lastlog.de>
