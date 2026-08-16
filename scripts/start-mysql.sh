#!/usr/bin/env sh
# Start or run the Docker MySQL container for local development
docker start myproject-mysql || docker run -d --name myproject-mysql \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  -e MYSQL_DATABASE=myapp \
  -p 3306:3306 mysql:8.0
