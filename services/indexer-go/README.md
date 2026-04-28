# indexer-go

Chain indexer service for the Titular protocol.

## Quickstart

```bash
go build ./...
go run .
```

## Build

```bash
docker build -t titular/indexer-go .
```

## Database migrations

The Postgres schema is owned by `internal/db` and applied via the `migrate`
binary in `cmd/migrate/`:

```bash
go build -o ./bin/migrate ./cmd/migrate
DATABASE_URL=postgres://titular:titular@localhost:5432/titular?sslmode=disable \
    ./bin/migrate up
./bin/migrate version
./bin/migrate down
```

Available subcommands: `up`, `down`, `steps <n>`, `version`, `force <v>`.
Run `./bin/migrate -h` for the full usage. See
[`internal/db/README.md`](internal/db/README.md) for migration authoring
conventions and the integration-test harness.
