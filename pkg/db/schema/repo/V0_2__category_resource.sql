CREATE TABLE category_resource (
	category_id VARCHAR(50) NOT NULL,
	resource_id VARCHAR(50) NOT NULL,
	status      VARCHAR(50) NOT NULL,
	create_time TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (category_id, resource_id)
);

CREATE INDEX category_resource_status_idx
	ON category_resource (status);
CREATE INDEX category_resource_create_time_idx
	ON category_resource (create_time);
CREATE INDEX category_resource_status_time_idx
	ON category_resource (status_time);
