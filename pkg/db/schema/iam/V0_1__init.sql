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
