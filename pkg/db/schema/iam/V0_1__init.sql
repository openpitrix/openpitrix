CREATE TABLE IF NOT EXISTS user (
	user_id     varchar(50) PRIMARY KEY NOT NULL,
	username    varchar(255)            NOT NULL,
	password    varchar(255)            NOT NULL,
	email       varchar(255)            NOT NULL,
	role        varchar(32)             NOT NULL,
	status      varchar(50)             NOT NULL,
	description text
);

CREATE TABLE IF NOT EXISTS user_client (
	client_id     varchar(50) PRIMARY KEY NOT NULL,
	user_id       varchar(50)             NOT NULL,
	client_secret varchar(255)            NOT NULL,
	status        varchar(50)             NOT NULL,
	description   text
);

CREATE TABLE IF NOT EXISTS group (
	group_id    varchar(50) PRIMARY KEY NOT NULL,
	name        varchar(255) DEFAULT '' NOT NULL,
	status      varchar(50)             NOT NULL,
	description text
);

CREATE TABLE IF NOT EXISTS group_member (
	group_id varchar(50) PRIMARY KEY NOT NULL,
	user_id  varchar(50) PRIMARY KEY NOT NULL,
);

CREATE TABLE IF NOT EXISTS user_password_reset (
	reset_id    varchar(50) PRIMARY KEY             NOT NULL,
	user_id     varchar(255)                        NOT NULL,
	status      varchar(50)                         NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
