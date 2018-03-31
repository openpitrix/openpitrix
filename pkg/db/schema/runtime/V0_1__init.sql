CREATE TABLE runtime (
	runtime_id            VARCHAR(50) PRIMARY KEY             NOT NULL,
	name                  VARCHAR(50)                         NOT NULL,
	description           TEXT                                NOT NULL,
	provider              VARCHAR(50)                         NOT NULL,
	runtime_url           TEXT                                NOT NULL,
	zone                  VARCHAR(50)                         NOT NULL,
	owner                 VARCHAR(50)                         NOT NULL,
	status                VARCHAR(50)                         NOT NULL,
	runtime_credential_id VARCHAR(50)                         NOT NULL,
	create_time           TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time           TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX runtime_name_idx
		ON runtime (name);

CREATE INDEX runtime_provider_idx
		ON runtime (provider);

CREATE INDEX runtime_zone_idx
		ON runtime (zone);

CREATE INDEX runtime_owner_idx
		ON runtime (owner);

CREATE INDEX runtime_status_idx
		ON runtime (status);

CREATE INDEX  runtime_create_time_idx
	ON runtime(create_time);

CREATE INDEX  runtime_status_time_idx
	ON runtime(status_time);

CREATE TABLE runtime_label (
	runtime_label_id VARCHAR(50) PRIMARY KEY             NOT NULL,
	runtime_id       VARCHAR(50)                         NOT NULL,
	label_key        VARCHAR(50)                         NOT NULL,
	label_value      VARCHAR(255)                        NOT NULL,
	create_time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX runtime_label_runtime_id_idx
		ON runtime_label (runtime_id);
CREATE INDEX runtime_label_label_key_idx
		ON runtime_label (label_key);
CREATE INDEX runtime_label_label_value_idx
		ON runtime_label (label_value);
CREATE INDEX runtime_label_create_time_idx
		ON runtime_label (create_time);

CREATE TABLE runtime_credential (
	runtime_credential_id VARCHAR(50) PRIMARY KEY             NOT NULL,
	content               JSON                                NOT NULL,
	create_time           TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX runtime_credential_create_time_idx
	ON runtime_credential (create_time);

