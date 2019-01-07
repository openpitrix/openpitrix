
ALTER TABLE runtime
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX runtime_owner_path_idx
	ON runtime (owner_path);
UPDATE runtime
SET owner_path = CONCAT(':', owner);

ALTER TABLE runtime_credential
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX runtime_credential_owner_path_idx
	ON runtime_credential (owner_path);
UPDATE runtime_credential
SET owner_path = CONCAT(':', owner);
