-- Add is_used column to password_resets table
ALTER TABLE password_resets ADD COLUMN is_used BOOLEAN DEFAULT false NOT NULL;
