
ALTER TABLE repo
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX repo_owner_path_idx
	ON repo (owner_path);
UPDATE repo
SET owner_path = CONCAT(':', owner);

ALTER TABLE repo_event
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX repo_event_owner_path_idx
	ON repo_event (owner_path);
UPDATE repo_event
SET owner_path = CONCAT(':', owner);
