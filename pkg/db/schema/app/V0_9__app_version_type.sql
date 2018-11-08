ALTER TABLE app_version
	ADD COLUMN type VARCHAR(50) NOT NULL;


CREATE INDEX app_version_type_idx
	ON app_version (type);