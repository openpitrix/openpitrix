ALTER table app_version_audit
	ADD review_id VARCHAR(50) NOT NULL;
ALTER table app_version
	ADD review_id VARCHAR(50) NOT NULL;

CREATE TABLE app_version_review
(
	review_id  VARCHAR(50) NOT NULL,
	version_id VARCHAR(50) NOT NULL,
	app_id     VARCHAR(50) NOT NULL,
	status     VARCHAR(50) NOT NULL,
	phase      JSON        NOT NULL,
	PRIMARY KEY (review_id)
);

CREATE INDEX app_version_review_version_id_idx
	ON app_version_review (version_id);
CREATE INDEX app_version_review_app_id_idx
	ON app_version_review (app_id);
CREATE INDEX app_version_review_status_idx
	ON app_version_review (status);

ALTER TABLE app
	ADD COLUMN tos TEXT NOT NULL;
ALTER TABLE app
	ADD COLUMN abstraction TEXT NOT NULL;
