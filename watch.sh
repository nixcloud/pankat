#!/bin/sh

while true; do 
  inotifywait -r -e modify *; 
  ./update.sh
done
