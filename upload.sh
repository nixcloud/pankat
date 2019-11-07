#!/bin/sh
rsync -avz --delete output/ --exclude=".*" root@nixcloud-root:/www/lastlog.de/blog
