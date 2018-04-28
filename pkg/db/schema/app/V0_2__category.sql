CREATE TABLE category (
	category_id VARCHAR(50)  NOT NULL,
	name        VARCHAR(255) NOT NULL,
	locale      JSON         NUll,
	owner       VARCHAR(50)  NOT NULL,
	create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time TIMESTAMP    NULL,
	PRIMARY KEY (category_id)
);

CREATE INDEX category_name_idx
	ON category (name);
CREATE INDEX category_owner_idx
	ON category (owner);
CREATE INDEX category_create_time_idx
	ON category (create_time);
CREATE INDEX category_update_time_idx
	ON category (update_time);

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
