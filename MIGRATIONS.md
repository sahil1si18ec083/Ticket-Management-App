# Migration Guide

## Understanding Migrations in This Project

This project uses **SQL-based migrations** with `golang-migrate`. Every schema change follows a structured process.

---

## Process for Adding a Constraint (Example: Name Length)

### Step 1: Create Migration Files
Name migrations with a **version number** and **descriptive name**:
- **Up migration** (apply change): `000002_add_name_length_constraint.up.sql`
- **Down migration** (rollback): `000002_add_name_length_constraint.down.sql`

Location: `migrations/sql/`

**Example `000002_add_name_length_constraint.up.sql`:**
```sql
-- Add CHECK constraint to ensure name is at least 6 characters
ALTER TABLE users
ADD CONSTRAINT check_name_length CHECK (LENGTH(name) >= 6);
```

**Example `000002_add_name_length_constraint.down.sql`:**
```sql
-- Remove the CHECK constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS check_name_length;
```

**Key points:**
- **Up migration**: Changes applied when migrating forward
- **Down migration**: Changes rolled back when migrating backward
- Use `IF EXISTS` / `IF NOT EXISTS` for safety (idempotency)
- Write descriptive comments explaining what each part does

---

### Step 2: Update the Go Model
Sync the Go struct with the database schema. Add the constraint tag to match the SQL constraint.

**Example `models/user.go`:**
```go
type User struct {
    gorm.Model
    Name     string `gorm:"type:varchar(255);not null;check:LENGTH(name) >= 6" json:"name"`
    Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
    Password string `gorm:"type:varchar(100);not null" json:"password"`
}
```

**Tags explained:**
- `gorm:"type:varchar(255)"` — database column type
- `not null` — disallow NULL values
- `check:LENGTH(name) >= 6` — constraint (GORM documents this for reference)
- `json:"name"` — JSON serialization key

---

### Step 3: Test the Migration
Run the migration to verify it applies without errors:

```powershell
# Set database URL
$env:DATABASE_URL = 'postgres://user:pass@localhost:5432/dbname?sslmode=disable'

# Apply the migration
go run ./cmd/migrate -up

# Rollback to test the down migration
go run ./cmd/migrate -down

# Re-apply
go run ./cmd/migrate -up
```

---

## File Naming Convention

Migration files must follow the pattern: `XXXXXX_description.{up,down}.sql`

- **XXXXXX**: Version number (6 digits, zero-padded)
  - `000001` — initial schema
  - `000002` — first alteration
  - `000003` — second alteration
  - etc.
- **description**: Snake-case description of the change
- **{up,down}**: Direction of migration

**Examples:**
- `000001_create_users_and_tickets.up.sql`
- `000002_add_name_length_constraint.up.sql`
- `000003_add_ticket_priority_column.up.sql`

---

## Common Patterns

### Adding a Column
**Up:**
```sql
ALTER TABLE users
ADD COLUMN phone VARCHAR(20);
```

**Down:**
```sql
ALTER TABLE users
DROP COLUMN phone;
```

### Creating an Index
**Up:**
```sql
CREATE INDEX idx_users_email ON users(email);
```

**Down:**
```sql
DROP INDEX IF EXISTS idx_users_email;
```

### Renaming a Column
**Up:**
```sql
ALTER TABLE users RENAME COLUMN old_name TO new_name;
```

**Down:**
```sql
ALTER TABLE users RENAME COLUMN new_name TO old_name;
```

### Adding a Foreign Key
**Up:**
```sql
ALTER TABLE tickets
ADD CONSTRAINT fk_ticket_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE;
```

**Down:**
```sql
ALTER TABLE tickets
DROP CONSTRAINT fk_ticket_user;
```

---

## Running Migrations

### Via CLI (recommended for deployments):
```powershell
# Apply all pending migrations
go run ./cmd/migrate -up

# Rollback one migration step
go run ./cmd/migrate -down

# Apply/rollback N steps
go run ./cmd/migrate -steps 2
```

### Automatic (on server startup):
The server runs `Migrate(db)` from `bootstrap/migrations.go` when it starts. Migrations apply automatically.

```powershell
go run ./cmd/main.go
```

---

## Troubleshooting

**Migration fails with "table already exists":**
- The `CREATE TABLE IF NOT EXISTS` clause prevents errors on retry.
- If you manually ran SQL, use `go run ./cmd/migrate` to track state in the database.

**Constraint already exists:**
- Use `DROP CONSTRAINT IF EXISTS` in down migrations.
- PostgreSQL tracks constraint names; verify with `\d users` in `psql`.

**Forgot to update the Go model:**
- GORM will not enforce the constraint on the Go side (app won't validate).
- Always keep models in sync with migrations for consistency.

**Need to see migration status:**
- Run `go run ./cmd/migrate -status` (if implemented) or check `schema_migrations` table:
  ```sql
  SELECT * FROM schema_migrations ORDER BY version DESC;
  ```

---

## Summary Checklist

For **each new migration**:
- [ ] Create `.up.sql` with the change
- [ ] Create `.down.sql` with the rollback
- [ ] Name files with 6-digit version and description
- [ ] Update Go models to match schema
- [ ] Run `go run ./cmd/migrate -up` to test
- [ ] Run `go run ./cmd/migrate -down` to test rollback
- [ ] Commit both SQL files and updated models together
