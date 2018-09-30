CREATE TABLE IF NOT EXISTS user (
	user_id     VARCHAR(50)  NOT NULL,
	username    VARCHAR(255) NOT NULL,
	password    VARCHAR(255) NOT NULL,
	email       VARCHAR(255) NOT NULL,
	role        VARCHAR(50)  NOT NULL DEFAULT 'user',
	status      VARCHAR(50)  NOT NULL,
	description TEXT,
	create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (user_id)
);

CREATE INDEX user_email_idx
	ON user (email);
CREATE INDEX user_status_idx
	ON user (status);
CREATE INDEX user_create_time_idx
	ON user (create_time);

CREATE TABLE IF NOT EXISTS user_client (
	client_id     VARCHAR(50)  NOT NULL,
	user_id       VARCHAR(50)  NOT NULL,
	client_secret VARCHAR(255) NOT NULL,
	status        VARCHAR(50)  NOT NULL,
	description   TEXT,
	create_time   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (client_id)
);

CREATE INDEX user_client_user_id_idx
	ON user_client (user_id);
CREATE INDEX user_client_status_idx
	ON user_client (status);
CREATE INDEX user_client_create_time_idx
	ON user_client (create_time);

CREATE TABLE IF NOT EXISTS `group` (
	group_id    VARCHAR(50)  NOT NULL,
	name        VARCHAR(255) NOT NULL DEFAULT '',
	status      VARCHAR(50)  NOT NULL,
	description TEXT,
	create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (group_id)
);

CREATE INDEX group_status_idx
	ON `group` (status);
CREATE INDEX group_create_time_idx
	ON `group` (create_time);

CREATE TABLE IF NOT EXISTS group_member (
	group_id    VARCHAR(50) NOT NULL,
	user_id     VARCHAR(50) NOT NULL,
	create_time TIMESTAMP   NOT NULL  DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (group_id, user_id)
);

CREATE INDEX group_member_create_time_idx
	ON group_member (create_time);

CREATE TABLE IF NOT EXISTS user_password_reset (
	reset_id    VARCHAR(50) NOT NULL,
	user_id     VARCHAR(50) NOT NULL,
	status      VARCHAR(50) NOT NULL,
	create_time TIMESTAMP   NOT NULL  DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (reset_id)
);

CREATE INDEX user_password_reset_status_idx
	ON user_password_reset (status);
CREATE INDEX user_password_reset_create_time_idx
	ON user_password_reset (create_time);
