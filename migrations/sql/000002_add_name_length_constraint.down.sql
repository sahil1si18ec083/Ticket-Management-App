-- Remove the CHECK constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS check_name_length;
