# Ticket Management App — Migrations

This project uses SQL-based migrations stored under `migrations/sql/` and the `golang-migrate` library to apply them.

Running migrations (Powershell):

1. Install the migration dependency and tidy modules:
```powershell
cd 'c:\Users\smart\Downloads\Files (2)\ticketing-app-gin-golang'
go get github.com/golang-migrate/migrate/v4
go mod tidy
```

2. Set your `DATABASE_URL` environment variable (example PostgreSQL):
```powershell
$env:DATABASE_URL = 'postgres://user:pass@localhost:5432/dbname?sslmode=disable'
```

3. Run migrations via the provided CLI (applies up migrations):
```powershell
go run ./cmd/migrate -up
```

4. Rollback (apply down):
```powershell
go run ./cmd/migrate -down
```

Notes
- Migration files are in `migrations/sql/` and follow the `000001_name.up.sql` / `.down.sql` convention.
- The server also runs migrations automatically on startup via `bootstrap.ConnectDatabase()`.
