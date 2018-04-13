CREATE TABLE repo (
	repo_id     VARCHAR(50) PRIMARY KEY             NOT NULL,
	name        VARCHAR(255)                        NOT NULL,
	description TEXT                                NOT NULL,
	url         TEXT                                NOT NULL,
	credential  JSON                                NOT NULL,
	type        VARCHAR(50)                         NOT NULL,
	visibility  VARCHAR(50)                         NOT NULL,
	owner       VARCHAR(50)                         NOT NULL,
	status      VARCHAR(50)                         NOT NULL,
	create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX repo_name_idx
	ON repo (name);
CREATE INDEX repo_type_idx
	ON repo (type);
CREATE INDEX repo_visibility_idx
	ON repo (visibility);
CREATE INDEX repo_owner_idx
	ON repo (owner);
CREATE INDEX repo_status_idx
	ON repo (status);
CREATE INDEX repo_create_time_idx
	ON repo (create_time);
CREATE INDEX repo_status_time_idx
	ON repo (status_time);

CREATE TABLE repo_provider (
	repo_id  VARCHAR(50) NOT NULL,
	provider VARCHAR(50) NOT NULL,
	PRIMARY KEY (repo_id, provider)
);

CREATE TABLE repo_label (
	repo_label_id VARCHAR(50) PRIMARY KEY                   NOT NULL,
	repo_id       VARCHAR(50)                               NOT NULL,
	label_key     VARCHAR(50)                               NOT NULL,
	label_value   VARCHAR(255)                              NOT NULL,
	create_time   TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);

CREATE INDEX repo_label_repo_id_idx
	ON repo_label (repo_id);
CREATE INDEX repo_label_label_key_idx
	ON repo_label (label_key);
CREATE INDEX repo_label_label_value_idx
	ON repo_label (label_value);
CREATE INDEX repo_label_create_time_idx
	ON repo_label (create_time);

CREATE TABLE repo_selector (
	repo_selector_id VARCHAR(50) PRIMARY KEY                   NOT NULL,
	repo_id          VARCHAR(50)                               NOT NULL,
	selector_key     VARCHAR(50)                               NOT NULL,
	selector_value   VARCHAR(255)                              NOT NULL,
	create_time      TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);

CREATE INDEX repo_selector_repo_id_idx
	ON repo_selector (repo_id);
CREATE INDEX repo_selector_selector_key_idx
	ON repo_selector (selector_key);
CREATE INDEX repo_selector_selector_value_idx
	ON repo_selector (selector_value);
CREATE INDEX repo_selector_create_time_idx
	ON repo_selector (create_time);

CREATE TABLE repo_event (
	repo_event_id VARCHAR(50) PRIMARY KEY             NOT NULL,
	repo_id       VARCHAR(50)                         NOT NULL,
	owner         VARCHAR(50)                         NOT NULL,
	status        VARCHAR(50)                         NOT NULL,
	result        TEXT                                NOT NULL,
	create_time   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX repo_event_repo_id_idx
	ON repo_event (repo_id);
CREATE INDEX repo_event_owner_idx
	ON repo_event (owner);
CREATE INDEX repo_event_status_idx
	ON repo_event (status);
CREATE INDEX repo_event_create_time_idx
	ON repo_event (create_time);
CREATE INDEX repo_event_status_time_idx
	ON repo_event (status_time);
