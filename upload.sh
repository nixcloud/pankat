#!/bin/sh
rsync -avz --delete output/ --exclude=".*" root@joshi-server-NG:/var/lib/nixcloud/webservices/apache-lastlog/www/blog
