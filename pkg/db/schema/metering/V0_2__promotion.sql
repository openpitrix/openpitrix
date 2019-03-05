/**
Promotion
**/

/** the combination spu **/
CREATE TABLE IF NOT EXISTS combination_spu (
	combination_spu_id VARCHAR(50) NOT NULL UNIQUE,
	spu_ids            JSON        NOT NULL
	COMMENT 'combination spu ids',
	create_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	status             VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted',
	PRIMARY KEY (combination_spu_id)
);

/** the sku of combination resources **/
CREATE TABLE IF NOT EXISTS combination_sku (
	combination_sku_id     VARCHAR(50) NOT NULL UNIQUE,
	combination_spu_id     VARCHAR(50) NOT NULL
	COMMENT 'the id of combination_spu_id',
	attribute_ids          JSON        NOT NULL
	COMMENT 'sku attributes of spu: {spuId:[attId, ..], ..}',
	metering_attribute_ids JSON        NOT NULL
	COMMENT 'sku metering attributes of spu: {spuId:[attId, ..], ..}',
	create_time            TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time            TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	status                 VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted',
	PRIMARY KEY (combination_sku_id)
);


CREATE TABLE IF NOT EXISTS combination_price (
	combination_price_id VARCHAR(50) NOT NULL UNIQUE,
	combination_sku_id   VARCHAR(50) NOT NULL,
	spu_id               VARCHAR(50) NOT NULL,
	attribute_id         VARCHAR(50) NOT NULL,
	prices               JSON COMMENT '{upto: price1, ...}',
	currency             VARCHAR(50) NOT NULL  DEFAULT 'cny',
	start_time           TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	end_time             TIMESTAMP   NULL,
	create_time          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	status_time          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	status               VARCHAR(16)           DEFAULT 'active'
	COMMENT 'active, deleted',
	INDEX price_sku_index (combination_price_id, combination_sku_id),
	PRIMARY KEY (combination_price_id)
);


/** probation sku of resource **/
CREATE TABLE IF NOT EXISTS probation_sku (
	probation_sku_id       VARCHAR(50) NOT NULL UNIQUE,
	resource_attribute_id  VARCHAR(50) NOT NULL,
	attribute_ids          JSON        NOT NULL
	COMMENT 'sku attributes of resource_attribute: [attributeId, ...]',
	metering_attribute_ids JSON        NOT NULL
	COMMENT 'sku attributes of resource_attribute: [attributeId, ...]',
	limit_num              INT         NOT NULL DEFAULT 1,
	create_time            TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
	status_time            TIMESTAMP            DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	status                 VARCHAR(16)          DEFAULT 'active'
	COMMENT 'active, deleted',
	PRIMARY KEY (probation_sku_id)
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
