/**
promotion
**/
CREATE TABLE IF NOT EXISTS combination_price (
	combination_price_id VARCHAR(50) NOT NULL UNIQUE,
	sku_id               VARCHAR(50) NOT NULL,
	attribute_id         VARCHAR(50) NOT NULL COMMENT 'metering attribute id',
	prices               JSON COMMENT '{upto: price1, ...}',
	currency             VARCHAR(50) NOT NULL  DEFAULT 'cny',
	status               VARCHAR(16)           DEFAULT 'active'
	COMMENT 'active, deleted',
	create_time          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	status_time          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (combination_price_id)
);


/** probation **/
CREATE TABLE IF NOT EXISTS probation (
	probation_id VARCHAR(50) NOT NULL UNIQUE,
	sku_id       VARCHAR(50) NOT NULL,
	value        JSON        NOT NULL
	COMMENT 'the value of probation: {attribute_id_1: value_1, attribute_id_2: value_2, ...}',
	status       VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted',
	start_time   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	end_time     TIMESTAMP,
	create_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (probation_id)
);


/** the records of probation resource used by user **/
CREATE TABLE IF NOT EXISTS probation_record (
	probation_id VARCHAR(50) NOT NULL,
	user_id      VARCHAR(50) NOT NULL,
	status       VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, used',
	create_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (probation_id, user_id)
);


CREATE TABLE IF NOT EXISTS dicount (
	discount_id      VARCHAR(50)  NOT NULL,
	name             VARCHAR(255) NOT NULL,
	owner            VARCHAR(50)  NOT NULL,
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
	coupon_id   VARCHAR(50)   NOT NULL,
	name        VARCHAR(50)   NOT NULL,
	owner       VARCHAR(50)   NOT NULL,
	limit_ids   JSON COMMENT '[spu_id1, .., sku_id1, .., price_id1, ..]',
	balance     DECIMAL(8, 2) NOT NULL,
	count       INT           NOT NULL,
	remain      INT           NOT NULL,
	limit_num   INT           NOT NULL DEFAULT 1,
	status      VARCHAR(16)            DEFAULT 'active'
	COMMENT 'active, deleted, overdue',
	start_time  TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	end_time    TIMESTAMP,
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	description TEXT,
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
