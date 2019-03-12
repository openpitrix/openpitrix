/**
promotion
**/
CREATE TABLE IF NOT EXISTS combination_price (
	combination_price_id   VARCHAR(50) NOT NULL UNIQUE,
	combination_binding_id VARCHAR(50) NOT NULL,
	prices                 JSON COMMENT '{upto: price1, ...}',
	currency               VARCHAR(50) NOT NULL  DEFAULT 'cny',
	status                 VARCHAR(16)           DEFAULT 'active'
	COMMENT 'active, deleted',
	create_time            TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	status_time            TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (combination_price_id)
);


/** probation **/
CREATE TABLE IF NOT EXISTS probation (
	probation_id VARCHAR(50) NOT NULL UNIQUE,
	sku_id       VARCHAR(50) NOT NULL,
	attribute_id VARCHAR(50) NOT NULL
	COMMENT 'the value in attribute of probation',
	status       VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted',
	start_time   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	end_time     TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	create_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (probation_id)
);


/** the records of probation resource used by user **/
CREATE TABLE IF NOT EXISTS probation_record (
	probation_id VARCHAR(50)   NOT NULL,
	user_id      VARCHAR(50)   NOT NULL,
	remain       DECIMAL(8, 2) NOT NULL,
	status       VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, used',
	create_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (probation_id, user_id)
);


CREATE TABLE IF NOT EXISTS dicount (
	discount_id      VARCHAR(50)  NOT NULL,
	owner            VARCHAR(50)  NOT NULL,
	name             VARCHAR(255) NOT NULL,
	limit_ids        JSON COMMENT '[spu_id1, .., sku_id1, .., price_id1, ..]',
	discount_value   DECIMAL(8, 2) COMMENT 'the price value to cut down',
	discount_percent DECIMAL(2, 2) COMMENT 'the price percent to cut down, there is only one of discount_value and discount_percent',
	status           VARCHAR(16) DEFAULT 'ready'
	COMMENT 'ready, active, deleted, overdue',
	start_time       TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	end_time         TIMESTAMP,
	create_time      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	description      TEXT,
	PRIMARY KEY (discount_id)
);


CREATE TABLE IF NOT EXISTS coupon (
	coupon_id     VARCHAR(50)   NOT NULL,
	name          VARCHAR(50)   NOT NULL,
	owner         VARCHAR(50)   NOT NULL,
	limit_ids     JSON COMMENT '[spu_id1, .., sku_id1, .., price_id1, ..]',
	balance       DECIMAL(8, 2) NOT NULL,
	count         INT           NOT NULL,
	remain        INT           NOT NULL,
	limit_num_per INT           NOT NULL DEFAULT 1,
	status        VARCHAR(16)            DEFAULT 'active'
	COMMENT 'active, deleted, overdue',
	start_time    TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	end_time      TIMESTAMP,
	create_time   TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	status_time   TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	description   TEXT,
	PRIMARY KEY (coupon_id)
);


CREATE TABLE IF NOT EXISTS coupon_received (
	coupon_received_id VARCHAR(50)   NOT NULL,
	coupon_id          VARCHAR(50)   NOT NULL,
	user_id            VARCHAR(50)   NOT NULL,
	remain             DECIMAL(8, 2) NOT NULL,
	status             VARCHAR(16) DEFAULT 'received'
	COMMENT 'received, active, used, overdue',
	create_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (coupon_received_id)
);


/** coupon used record **/
CREATE TABLE IF NOT EXISTS coupon_used (
	coupon_used_id VARCHAR(50)   NOT NULL,
	coupon_received_id VARCHAR(50)   NOT NULL,
	contract_id          VARCHAR(50)   NOT NULL,
	balance             DECIMAL(8, 2) NOT NULL,
	currency            VARCHAR(50) NOT NULL  DEFAULT 'cny',
	status             VARCHAR(16) DEFAULT 'undetermined'
	COMMENT 'undetermined --> done / refunded',
	create_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (coupon_used_id)
);
