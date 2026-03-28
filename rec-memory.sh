#!/usr/bin/env sh

while [ ! -f server.pid ]; do sleep 0.1; done
PID=$(cat server.pid)

> mem.log

while true; do
  echo "$(date '+%Y-%m-%d %H:%M:%S') $(ps -o rss= -p $PID)" >> mem.log
  sleep 1
done
