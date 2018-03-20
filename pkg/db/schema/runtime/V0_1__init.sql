CREATE TABLE runtime_env (
	runtime_env_id  VARCHAR(50) PRIMARY KEY             NOT NULL,
	name            VARCHAR(50)                         NOT NULL,
	description     TEXT                                NOT NULL,
	runtime_env_url TEXT                                NOT NULL,
	owner           VARCHAR(50)                         NOT NULL,
	status          VARCHAR(50)                         NOT NULL,
	create_time     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE runtime_env_label (
	runtime_env_label_id VARCHAR(50) PRIMARY KEY NOT NULL,
	runtime_env_id       VARCHAR(50)             NOT NULL,
	label_key            VARCHAR(50)             NOT NULL,
	label_value          TEXT                    NOT NULL
);

CREATE TABLE runtime_env_credential (
	runtime_env_credential_id VARCHAR(50) PRIMARY KEY             NOT NULL,
	name                      VARCHAR(50)                         NOT NULL,
	description               TEXT                                NOT NULL,
	owner                     VARCHAR(50)                         NOT NULL,
	content                   JSON,
	status                    VARCHAR(50)                         NOT NULL,
	create_time               TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time               TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE runtime_env_attached_credential (
	runtime_env_id            VARCHAR(50) NOT NULL,
	runtime_env_credential_id VARCHAR(50) NOT NULL,
	PRIMARY KEY (runtime_env_id, runtime_env_credential_id)
);
