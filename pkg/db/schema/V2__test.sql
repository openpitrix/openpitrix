CREATE TABLE test_selector (
	repo_selector_id  VARCHAR(50) PRIMARY KEY             NOT NULL,
	repo_id           VARCHAR(50)                         NOT NULL,
	selector_key      VARCHAR(50)                         NOT NULL,
	selector_value    VARCHAR(50)                         NOT NULL,
	status            VARCHAR(50)                         NOT NULL,
	create_time       TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	status_time       TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
