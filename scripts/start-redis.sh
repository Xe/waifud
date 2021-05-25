#!/bin/sh

docker run \
  --name waifud-redis \
  -p 6379:6379 \
  --restart always \
  -d \
  redis:6.3.2
