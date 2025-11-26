# Database Migrations - Name Length Constraint

## Overview

This document covers the migration we added to enforce a minimum name length of 6 characters in the `users` table.

---

## Migration Details

### Files Created/Modified

**Migration Files:**
- `migrations/sql/000002_add_name_length_constraint.up.sql` — Adds the constraint
- `migrations/sql/000002_add_name_length_constraint.down.sql` — Removes the constraint

**Model File:**
- `models/user.go` — Updated with `check:LENGTH(name) >= 6` tag

**API Request Validation:**
- `models/auth_requests.go` — Updated `SignupRequest.Name` with `min=6` binding

---

## What Was Changed

### 1. Database Migration (Up)
**File:** `migrations/sql/000002_add_name_length_constraint.up.sql`

```sql
-- Add CHECK constraint to ensure name is at least 6 characters
ALTER TABLE users
ADD CONSTRAINT check_name_length CHECK (LENGTH(name) >= 6);
```

**What it does:**
- Adds a PostgreSQL `CHECK` constraint to the `users.name` column
- Enforces that any `name` value must have a length of at least 6 characters
- The constraint applies to all new inserts and updates

### 2. Database Migration (Down/Rollback)
**File:** `migrations/sql/000002_add_name_length_constraint.down.sql`

```sql
-- Remove the CHECK constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS check_name_length;
```

**What it does:**
- Removes the constraint if it exists
- Allows rolling back to the previous schema if needed
- Safe to run even if the constraint was already removed (`IF EXISTS`)

### 3. Go Model Update
**File:** `models/user.go`

```go
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null;check:LENGTH(name) >= 6" json:"name"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
}
```

**What changed:**
- Added `check:LENGTH(name) >= 6` tag to the `Name` field
- This documents the constraint for GORM and serves as reference for developers

### 4. API Request Validation
**File:** `models/auth_requests.go`

```go
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
```

**What changed:**
- Added `min=6` validation to the `Name` field
- Gin framework validates this on the request before it reaches your handler
- Provides immediate feedback to API clients with validation errors

---

## Why Three Layers of Validation?

| Layer | Purpose | When It Works |
|-------|---------|---------------|
| **API Validation** (`Gin` binding) | Fast feedback to client | Request received |
| **Go Model** (GORM tags) | Reference for developers | Documentation |
| **Database Constraint** (SQL CHECK) | Last-line defense | Insert/update to DB |

This defense-in-depth approach ensures:
- ✅ Invalid data doesn't reach the database
- ✅ Existing data in DB is always valid
- ✅ Direct database inserts also respect the constraint

---

## How to Apply This Migration

### Prerequisites
- PostgreSQL database running and accessible
- `.env` file with `DATABASE_URL` configured
- Project built successfully (`go build ./...`)

### Commands to Run

**1. Set your database URL (if not in `.env`):**
```powershell
$env:DATABASE_URL = 'postgresql://user:password@host:port/dbname?sslmode=require'
```

**2. Apply all pending migrations (including this one):**
```powershell
cd 'c:\Users\smart\Downloads\Files (2)\ticketing-app-gin-golang'
go run ./cmd/migrate -up
```

**3. Verify the migration applied:**
```powershell
# Connect to PostgreSQL and check
psql -U neondb_owner -h ep-winter-frost-adsed2jl-pooler.c-2.us-east-1.aws.neon.tech -d neondb

# In psql terminal, check the constraint:
\d users
```

You should see:
```
Check constraints:
    "check_name_length" CHECK (length(name::text) >= 6)
```

---

## How to Rollback (If Needed)

**To undo this migration:**
```powershell
cd 'c:\Users\smart\Downloads\Files (2)\ticketing-app-gin-golang'
go run ./cmd/migrate -down
```

This will:
1. Run the `.down.sql` file
2. Remove the `check_name_length` constraint from the `users` table
3. Revert to the previous schema

---

## Testing the Constraint

### Test 1: API Validation (Should Reject)
```bash
# Try to signup with name < 6 characters
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob",
    "email": "bob@example.com",
    "password": "password123"
  }'

# Response: 400 Bad Request (Gin validates min=6)
# Error: "Key: 'SignupRequest.Name' Error:Field validation for 'Name' failed on the 'min' tag"
```

### Test 2: API Validation (Should Accept)
```bash
# Try to signup with name >= 6 characters
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Robert",
    "email": "robert@example.com",
    "password": "password123"
  }'

# Response: 200 OK (or appropriate response)
```

### Test 3: Database Constraint (Direct Insert)
```sql
-- This will fail (if someone bypasses the API):
INSERT INTO users (name, email, password) 
VALUES ('Bob', 'bob@test.com', 'hashed_password');
-- Error: new row for relation "users" violates check constraint "check_name_length"

-- This will succeed:
INSERT INTO users (name, email, password) 
VALUES ('Robert', 'robert@test.com', 'hashed_password');
-- Success
```

---

## Migration File Structure

All migrations follow the naming convention:
```
XXXXXX_description.{up,down}.sql
```

**Current migrations:**
- `000001_create_users_and_tickets.{up,down}.sql` — Initial schema
- `000002_add_name_length_constraint.{up,down}.sql` — Name constraint (this one)

**For future migrations:**
- `000003_your_next_change.{up,down}.sql`
- `000004_another_change.{up,down}.sql`
- etc.

---

## Troubleshooting

### Error: "constraint already exists"
**Cause:** Migration was already applied.

**Solution:** Check migration status:
```sql
SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 5;
```

If `000002` is listed, the migration already ran.

### Error: "DATABASE_URL not set"
**Cause:** Environment variable not configured.

**Solution:** Set it before running:
```powershell
$env:DATABASE_URL = 'postgresql://...'
go run ./cmd/migrate -up
```

### Error: "column does not exist"
**Cause:** Previous migration didn't run.

**Solution:** Apply all migrations:
```powershell
go run ./cmd/migrate -up
```

---

## Summary

**What we did:**
1. ✅ Created SQL migration to add name length constraint
2. ✅ Added rollback migration for safety
3. ✅ Updated Go model to document the constraint
4. ✅ Added API-level validation in `SignupRequest`
5. ✅ Documented the entire process

**How to use:**
- Run: `go run ./cmd/migrate -up`
- Rollback: `go run ./cmd/migrate -down`
- Auto-run on startup: Server runs migrations on boot

**Defense in depth:**
- API validates early (Gin binding)
- Database enforces constraint (SQL CHECK)
- Go model documents intent

---

## Next Steps

To add more migrations:
1. Create `000003_description.up.sql` and `.down.sql` files in `migrations/sql/`
2. Update relevant Go models
3. Run `go run ./cmd/migrate -up`

See `MIGRATIONS.md` in the project root for detailed migration guidelines.
