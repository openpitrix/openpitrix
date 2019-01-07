
ALTER TABLE job
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX job_owner_path_idx
	ON job (owner_path);
UPDATE job
SET owner_path = CONCAT(':', owner);
