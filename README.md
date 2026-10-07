# Social API

A Go REST API backed by PostgreSQL. The API currently exposes a health-check
endpoint and includes storage and database migration scaffolding for users and
posts.

## Requirements

- Go 1.26.5 or newer
- Docker and Docker Compose
- GNU Make
- The [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI

## Getting started

1. Copy the example environment file:

   ```powershell
   Copy-Item .envrc.example .envrc
   ```

   Update `.envrc` if your local PostgreSQL credentials or connection details
   differ from the defaults.

2. Start PostgreSQL:

   ```powershell
   docker compose up -d
   ```

   PostgreSQL is exposed on `localhost:5433` and uses the database and
   credentials configured in `.envrc`.

3. Apply the database migrations:

   ```powershell
   make migrate-up
   ```

4. Start the API:

   ```powershell
   go run .\cmd\api
   ```

   The server listens on `http://localhost:8080` by default.

5. Verify that the API is running:

   ```powershell
   Invoke-WebRequest http://localhost:8080/api/v1/health
   ```

   A healthy server responds with `ok`.

## Development

To roll back the latest migration:

```powershell
make migrate-down
```

The repository includes an `.air.toml` configuration for live reloading with
[Air](https://github.com/air-verse/air):

```powershell
air
```

## Project layout

```text
cmd/
  api/                 HTTP API entry point and routes
  migrate/migrations/  SQL database migrations
internal/
  db/                  PostgreSQL connection setup
  env/                 Environment loading helpers
  store/               Data access layer
scripts/               Database initialization scripts
docker-compose.yml     Local PostgreSQL service
Makefile               Migration commands
```

## Configuration

Configuration is loaded from `.envrc` when the API starts. The available
settings are documented in `.envrc.example`, including:

- `ADDR`
- `DB_ADDR`
- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`
- `DB_MAX_IDLE_TIME`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`

