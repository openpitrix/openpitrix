CREATE TABLE attachment (
	attachment_id   VARCHAR(50) NOT NULL,
	create_time     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (attachment_id)
);

CREATE INDEX attachment_create_time_idx
	ON attachment (create_time);