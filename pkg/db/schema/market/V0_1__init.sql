CREATE TABLE IF NOT EXISTS market (
	market_id   VARCHAR(50)  NOT NULL,
	name        VARCHAR(255) NOT NULL,
	visibility  VARCHAR(50)  NOT NULL,
	status      VARCHAR(50)  NOT NULL,
	owner       VARCHAR(50)  NOT NULL,
	description TEXT         NOT NULL,
	create_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (market_id)
);

CREATE INDEX market_id_idx
	ON market (market_id);
CREATE INDEX market_name_idx
	ON market (name);
CREATE INDEX market_visibility_idx
	ON market (visibility);
CREATE INDEX market_status_idx
	ON market (status);
CREATE INDEX market_owner_idx
	ON market (owner);
CREATE INDEX market_create_time_idx
	ON market (create_time);
CREATE INDEX market_status_time_idx
	ON market (status_time);

CREATE TABLE IF NOT EXISTS market_user (
	market_id   VARCHAR(50) NOT NULL,
	user_id     VARCHAR(50) NOT NULL,
	owner       VARCHAR(50) NOT NULL,
	create_time TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (market_id,user_id)
);

CREATE INDEX market_user_market_id_idx
	ON market_user (market_id);
CREATE INDEX market_user_user_id_idx
	ON market_user (user_id);
CREATE INDEX market_user_owner_idx
	ON market_user (owner);
CREATE INDEX market_user_create_time_idx
	ON market_user (create_time);