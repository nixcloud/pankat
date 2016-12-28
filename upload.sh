#!/bin/sh
rsync -avz --delete output/ --exclude=".*" root@nixcloud:/www/lastlog.de/blog
