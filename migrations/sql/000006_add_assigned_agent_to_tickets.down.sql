ALTER TABLE tickets
DROP CONSTRAINT IF EXISTS fk_tickets_assigned_agent;

ALTER TABLE tickets
DROP COLUMN IF EXISTS assigned_agent_id;
