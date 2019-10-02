#!/usr/bin/env sh
docker build -t omnitureproxy .
docker run -t --rm -p 3000:3000 omnitureproxy