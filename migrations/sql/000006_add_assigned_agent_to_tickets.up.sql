ALTER TABLE tickets
ADD COLUMN IF NOT EXISTS assigned_agent_id BIGINT;

ALTER TABLE tickets
ADD CONSTRAINT fk_tickets_assigned_agent
FOREIGN KEY (assigned_agent_id)
REFERENCES users(id)
ON UPDATE CASCADE
ON DELETE SET NULL;
