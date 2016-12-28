#!/bin/sh

while true; do 
  inotifywait -r -e modify *; 
  ./pankat -i documents -o output; 
done
