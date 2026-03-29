#!/usr/bin/env sh

while [ ! -f cli.pid ]; do sleep 0.1; done
PID=$(cat cli.pid)

> mem.log

while kill -0 "$PID" 2>/dev/null; do
  RSS=$(ps -o rss= -p "$PID" 2>/dev/null)
  [ -n "$RSS" ] && echo "$(date '+%Y-%m-%d %H:%M:%S') $RSS" >> mem.log
  sleep 0.2
done
