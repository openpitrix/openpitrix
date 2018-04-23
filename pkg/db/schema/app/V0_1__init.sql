CREATE TABLE app (
	app_id      VARCHAR(50) PRIMARY KEY             NOT NULL,
	name        VARCHAR(255)                        NOT NULL,
	description TEXT                                NOT NULL,
	icon        TEXT                                NOT NULL,
	home        TEXT                                NOT NULL,
	readme      TEXT                                NOT NULL,
	repo_id     VARCHAR(50)                         NOT NULL,
	chart_name  TEXT                                NOT NULL,
	screenshots TEXT                                NOT NULL,
	maintainers TEXT                                NOT NULL,
	sources     TEXT                                NOT NULL,
	owner       VARCHAR(50)                         NOT NULL,
	status      VARCHAR(50)                         NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time TIMESTAMP                           NULL
);

CREATE TABLE app_version (
	version_id   VARCHAR(50) PRIMARY KEY             NOT NULL,
	app_id       VARCHAR(50)                         NOT NULL,
	name         VARCHAR(255)                        NOT NULL,
	description  TEXT                                NOT NULL,
	package_name TEXT                                NOT NULL,
	owner        VARCHAR(50)                         NOT NULL,
	status       VARCHAR(50)                         NOT NULL,
	create_time  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time  TIMESTAMP                           NULL
);
