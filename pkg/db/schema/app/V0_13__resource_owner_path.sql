ALTER TABLE app
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX app_owner_path_idx
	ON app (owner_path);
UPDATE app
SET owner_path = CONCAT(':', owner);

ALTER TABLE app_version
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX app_version_owner_path_idx
	ON app_version (owner_path);
UPDATE app_version
SET owner_path = CONCAT(':', owner);

ALTER TABLE app_version_review
	ADD COLUMN owner VARCHAR(255) NOT NULL;
CREATE INDEX app_version_review_owner_idx
	ON app_version_review (owner);

ALTER TABLE app_version_review
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX app_version_review_owner_path_idx
	ON app_version_review (owner_path);
UPDATE app_version_review
SET owner_path = CONCAT(':', owner);


ALTER TABLE app_version_audit
	ADD COLUMN owner VARCHAR(255) NOT NULL;
CREATE INDEX app_version_audit_owner_idx
	ON app_version_audit (owner);


ALTER TABLE app_version_audit
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX app_version_audit_owner_path_idx
	ON app_version_audit (owner_path);
UPDATE app_version_audit
SET owner_path = CONCAT(':', owner);

ALTER TABLE category
	ADD COLUMN owner_path VARCHAR(50) NOT NULL;
CREATE INDEX category_owner_path_idx
	ON category (owner_path);
UPDATE category
SET owner_path = CONCAT(':', owner);
