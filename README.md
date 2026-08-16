# My Project — Local MySQL via Docker

Start a local MySQL server for development (no migrations will be applied):

```sh
docker start myproject-mysql || docker run -d --name myproject-mysql \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  -e MYSQL_DATABASE=myapp \
  -p 3306:3306 mysql:8.0
```

You can also use the Makefile target:

```sh
make start-mysql
```

Run the API server after MySQL is up:

```sh
go run ./cmd
```

The application reads DB settings from environment variables; defaults are configured for a local Docker MySQL (user `root`, empty password, host `127.0.0.1:3306`, database `myapp`).
