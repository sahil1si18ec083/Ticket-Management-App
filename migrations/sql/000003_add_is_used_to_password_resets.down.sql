-- Remove is_used column from password_resets table
ALTER TABLE password_resets DROP COLUMN is_used;
