-- Add CHECK constraint to ensure name is at least 6 characters
ALTER TABLE users
ADD CONSTRAINT check_name_length CHECK (LENGTH(name) >= 6);
