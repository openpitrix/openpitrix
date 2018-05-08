ALTER TABLE app_version
	ADD COLUMN sequence int DEFAULT 0 NOT NULL;

CREATE INDEX app_version_sequence_idx
	ON app_version (sequence);
