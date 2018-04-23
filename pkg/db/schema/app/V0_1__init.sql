CREATE TABLE app (
	app_id      VARCHAR(50)  NOT NULL,
	name        VARCHAR(255) NOT NULL,
	description TEXT         NOT NULL,
	icon        TEXT         NOT NULL,
	home        TEXT         NOT NULL,
	readme      TEXT         NOT NULL,
	repo_id     VARCHAR(50)  NOT NULL,
	chart_name  TEXT         NOT NULL,
	screenshots TEXT         NOT NULL,
	maintainers TEXT         NOT NULL,
	sources     TEXT         NOT NULL,
	owner       VARCHAR(50)  NOT NULL,
	status      VARCHAR(50)  NOT NULL,
	create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time TIMESTAMP    NULL,
	PRIMARY KEY (app_id)
);

CREATE INDEX app_name_idx
	ON app (name);
CREATE INDEX app_chart_name_idx
	ON app (chart_name(767));
CREATE INDEX app_repo_id_idx
	ON app (repo_id);
CREATE INDEX app_owner_idx
	ON app (owner);
CREATE INDEX app_status_idx
	ON app (status);
CREATE INDEX app_create_time_idx
	ON app (create_time);
CREATE INDEX app_status_time_idx
	ON app (status_time);
CREATE INDEX app_update_time_idx
	ON app (update_time);

CREATE TABLE app_version (
	version_id   VARCHAR(50)  NOT NULL,
	app_id       VARCHAR(50)  NOT NULL,
	name         VARCHAR(255) NOT NULL,
	description  TEXT         NOT NULL,
	package_name TEXT         NOT NULL,
	owner        VARCHAR(50)  NOT NULL,
	status       VARCHAR(50)  NOT NULL,
	create_time  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time  TIMESTAMP    NULL,
	PRIMARY KEY (version_id)
);

CREATE INDEX app_version_app_id_idx
	ON app_version (app_id);
CREATE INDEX app_version_name_idx
	ON app_version (name);
CREATE INDEX app_version_package_name_idx
	ON app_version (package_name(767));
CREATE INDEX app_version_owner_idx
	ON app_version (owner);
CREATE INDEX app_version_status_idx
	ON app_version (status);
CREATE INDEX app_version_create_time_idx
	ON app_version (create_time);
CREATE INDEX app_version_status_time_idx
	ON app_version (status_time);
CREATE INDEX app_version_update_time_idx
	ON app_version (update_time);
