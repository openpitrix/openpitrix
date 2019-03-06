/**
promotion
**/
CREATE TABLE IF NOT EXISTS combination_price (
	combination_price_id VARCHAR(50) NOT NULL UNIQUE,
	combination_binding_id   VARCHAR(50) NOT NULL,
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
	probation_id       VARCHAR(50) NOT NULL UNIQUE,
	sku_id  VARCHAR(50) NOT NULL,
	attribute_id          VARCHAR(50) NOT NULL
	COMMENT 'the value in attribute of probation',
	status                 VARCHAR(16)          DEFAULT 'active'
	COMMENT 'active, deleted',
	start_time            TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
	end_time            TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
	create_time            TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
	status_time            TIMESTAMP            DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (probation_id)
);




/** the records of probation resource used by user **/
CREATE TABLE IF NOT EXISTS probation_record (
	probation_sku_id VARCHAR(50) NOT NULL,
	user_id          VARCHAR(50) NOT NULL,
	num              INT         NOT NULL DEFAULT 1,
	create_time      TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
	probation_times  JSON COMMENT '[probation_time1, ...]',
	PRIMARY KEY (probation_sku_id, user_id)
);

CREATE TABLE IF NOT EXISTS dicount (
	discount_id      VARCHAR(50)  NOT NULL,
	name             VARCHAR(255) NOT NULL,
	limits           JSON COMMENT '{resource:.., sku:.., price:.., user:.., regoin:..}',
	discount_value   DECIMAL(8, 2) COMMENT 'the price value to cut down',
	discount_percent DECIMAL(2, 2) COMMENT 'the price percent to cut down, there is only one of discount_value and discount_percent',
	start_time       TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	end_time         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	create_time      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status           VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted, overtime',
	mark             TEXT,
	PRIMARY KEY (discount_id)
);

CREATE TABLE IF NOT EXISTS coupon (
	coupon_id   VARCHAR(50)   NOT NULL,
	name        VARCHAR(50)   NOT NULL,
	limits      JSON COMMENT '{resource:.., sku:.., price:.., user:.., regoin:...}',
	balance     DECIMAL(8, 2) NOT NULL,
	count       INT           NOT NULL,
	limit_num   INT           NOT NULL DEFAULT 1,
	start_time  TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	end_time    TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	status      VARCHAR(16)            DEFAULT 'active'
	COMMENT 'active, deleted, overtime',
	mark        TEXT,
	PRIMARY KEY (coupon_id)
);

CREATE TABLE IF NOT EXISTS coupon_received (
	coupon_received_id VARCHAR(50)   NOT NULL,
	coupon_id          VARCHAR(50)   NOT NULL,
	user_id            VARCHAR(50)   NOT NULL,
	balance            DECIMAL(8, 2) NOT NULL,
	status             VARCHAR(16) DEFAULT 'received'
	COMMENT 'received, active, overtime',
	create_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (coupon_received_id)
);
