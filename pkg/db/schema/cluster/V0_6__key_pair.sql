CREATE TABLE IF NOT EXISTS key_pair (
	key_pair_id VARCHAR(50)   NOT NULL,
	name        VARCHAR(50)   NULL,
	description VARCHAR(1000) NULL,
	pub_key         TEXT          NOT NULL,
	create_time TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	owner       VARCHAR(255)  NOT NULL,
	INDEX key_pair_name_index (name ASC),
	INDEX key_pair_create_time_index (create_time ASC),
	INDEX key_pair_key_index (pub_key(767)),
	INDEX key_pair_owner_index (owner ASC),
	PRIMARY KEY (key_pair_id)
);

CREATE TABLE IF NOT EXISTS node_key_pair (
  node_id     varchar(50) NOT NULL,
  key_pair_id varchar(50) NOT NULL,
  INDEX node_key_pair_node_id_index (node_id ASC),
  INDEX node_key_pair_key_pair_id_index (key_pair_id ASC),
  PRIMARY KEY(node_id, key_pair_id)
);
