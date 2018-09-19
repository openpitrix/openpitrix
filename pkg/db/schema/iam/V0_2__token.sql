CREATE TABLE IF NOT EXISTS token (
	token_id     VARCHAR(255)  NOT NULL,
	client_id     VARCHAR(50)  NOT NULL,
	refresh_token    VARCHAR(255) NOT NULL,
	scope    VARCHAR(255) NOT NULL,
	user_id    VARCHAR(50) NOT NULL,
	status      VARCHAR(50)  NOT NULL,
	create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

	UNIQUE KEY (refresh_token),
	PRIMARY KEY (token_id)
);

CREATE INDEX token_user_id_idx
	ON token (user_id);
CREATE INDEX token_status_idx
	ON token (status);
CREATE INDEX token_create_time_idx
	ON token (create_time);
CREATE INDEX token_refresh_token_idx
	ON token (refresh_token);
