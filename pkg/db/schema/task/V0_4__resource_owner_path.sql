
ALTER TABLE task
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX task_owner_path_idx
	ON task (owner_path);
UPDATE task
SET owner_path = CONCAT(':', owner);
