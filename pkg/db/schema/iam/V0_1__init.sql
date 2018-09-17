CREATE TABLE IF NOT EXISTS user (
	id varchar(50) PRIMARY KEY NOT NULL,
	name varchar(255) NOT NULL,
	password varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	role varchar(32) NOT NULL,
	description text
);

CREATE TABLE IF NOT EXISTS user_client (
	user_id varchar(50) PRIMARY KEY NOT NULL,
	client_id varchar(50) PRIMARY KEY NOT NULL,
	client_secret varchar(255) NOT NULL,
	description text
);

CREATE TABLE IF NOT EXISTS group (
	id varchar(50) PRIMARY KEY NOT NULL,
	name varchar(255) NOT NULL,
	description text
);

CREATE TABLE IF NOT EXISTS group_member (
	group_id varchar(50) NOT NULL,
	user_id varchar(50) NOT NULL,
	UNIQUE INDEX group_member_uniq_index (group_id, user_id)
);
