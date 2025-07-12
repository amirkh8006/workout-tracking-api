# Workout Tracking API
This is a code repository for a workout tracking API project.

## Setup

The API project is built from scratch. Before watching the course, you should install:
- [Go](https://go.dev/doc/install) (version 1.24.2 or higher)
- [Postgres](https://www.postgresql.org/download/) and any DB tool like psql or Sequel Ace to run SQL queries.
- [Docker and Docker Compose](https://www.docker.com/)
- [Goose](https://github.com/pressly/goose) for migration tool


## Setup Tips
- The Docker container exposes Postgres on the default port of `5432`. If you already have Postgres or something else running on that port and you get a connection error, you can use an alternate port but updating the `docker-compose.yml` to be something like `"5433:5432"`.
- If you get a "command not found" error when running `goose -version`, it's because the `$HOME/go/bin` directory is not added to your `PATH`. You can fix this temporarily by running `export PATH=$HOME/go/bin:$PATH`, but this will not persist if you close your terminal. A permanent fix would require adding `export PATH=$HOME/go/bin:$PATH` to your `.zshrc` or `.bashrc`.

## Tests

```bash
$ cd internal/store/
$ go test .
```