ALTER TABLE app_version_review
	ADD COLUMN reviewer VARCHAR(50) NOT NULL;
CREATE INDEX app_version_review_reviewer_idx
	ON app_version_review (reviewer);

ALTER TABLE app_version_review
	ADD COLUMN status_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
CREATE INDEX app_version_review_status_time_idx
	ON app_version_review (status_time);
