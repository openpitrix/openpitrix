/**
Price
**/
CREATE TABLE IF NOT EXISTS price (
	price_id     VARCHAR(50) NOT NULL UNIQUE,
	binding_id   VARCHAR(50) NOT NULL,
	prices       JSON COMMENT '{upto: price1, ...}',
	currency     VARCHAR(10) NOT NULL  DEFAULT 'cny',
	status       VARCHAR(16)           DEFAULT 'active'
	COMMENT 'active, deleted, disabled',
	start_time   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	end_time     TIMESTAMP   NULL,
	create_time  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	status_time  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (price_id)
);

/** Billing **/
CREATE TABLE IF NOT EXISTS leasing_contract (
	id              VARCHAR(50)   NOT NULL UNIQUE,
	leasing_id      VARCHAR(50)   NOT NULL,
	sku_id          VARCHAR(50)   NOT NULL,
	user_id         VARCHAR(50)   NOT NULL,
	metering_values JSON          NOT NULL,
	start_time      TIMESTAMP     NULL
	COMMENT 'same as leasing_time',
	status_time     TIMESTAMP     NULL,
	create_time     TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	fee_info        TEXT,
	fee             DECIMAL(8, 2) NOT NULL
	COMMENT 'total fee from starting cluster to now',
	due_fee         DECIMAL(8, 2) NOT NULL default 0,
	done_fee        DECIMAL(8, 2) NOT NULL default 0,
	before_bill_fee DECIMAL(8, 2) NOT NULL DEFAULT 0
	COMMENT 'the total fee of the before bills ',
	currency        VARCHAR(10)   NOT NULL DEFAULT 'cny',
	remark          TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS leased_contract (
	contract_id     VARCHAR(50)   NOT NULL UNIQUE,
	leasing_id      VARCHAR(50)   NOT NULL,
	sku_id          VARCHAR(50)   NOT NULL,
	user_id         VARCHAR(50)   NOT NULL,
	metering_values JSON          NOT NULL,
	start_time      TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	end_time        TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	create_time     TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	fee_info        TEXT,
	fee             DECIMAL(8, 2) NOT NULL
	COMMENT 'total fee from starting cluster to now',
	due_fee         DECIMAL(8, 2) NOT NULL default 0,
	done_fee        DECIMAL(8, 2) NOT NULL default 0,
	before_bill_fee DECIMAL(8, 2) NOT NULL DEFAULT 0
	COMMENT 'the total fee of the before bills ',
	currency        VARCHAR(10)   NOT NULL DEFAULT 'cny',
	remark          TEXT,
	PRIMARY KEY (contract_id)
);


/** Charge **/
CREATE TABLE IF NOT EXISTS charge (
	id          VARCHAR(50)   NOT NULL UNIQUE,
	user_id     VARCHAR(50)   NOT NULL,
	contract_id VARCHAR(50)   NOT NULL,
	fee         DECIMAL(8, 2) NOT NULL
	COMMENT 'total fee from starting cluster to now',
	currency    VARCHAR(10)   NOT NULL DEFAULT 'cny',
	info        JSON COMMENT '{couponReceivedID: fee}',
	status      VARCHAR(16)            DEFAULT 'successful'
	COMMENT 'successful, failed',
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	remark      TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS recharge (
	id          VARCHAR(50)   NOT NULL UNIQUE,
	user_id     VARCHAR(50)   NOT NULL,
	contract_id VARCHAR(50),
	fee         DECIMAL(8, 2) NOT NULL
	COMMENT 'total fee from starting cluster to now',
	currency    VARCHAR(10)   NOT NULL DEFAULT 'cny',
	info        JSON COMMENT '{couponReceivedID: fee}',
	status      VARCHAR(16)            DEFAULT 'successful'
	COMMENT 'successful, failed',
	operator    VARCHAR(50)   NOT NULL DEFAULT 'system',
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	remark      TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS income (
	id                  VARCHAR(50)   NOT NULL UNIQUE,
	user_id             VARCHAR(50)   NOT NULL,
	contract_id         VARCHAR(50)   NOT NULL,
	resource_version_id VARCHAR(50)   NOT NULL,
	balance             DECIMAL(9, 2) NOT NULL DEFAULT 0,
	currency            VARCHAR(10)   NOT NULL DEFAULT 'cny',
	create_time         TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	remark              TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS account (
	user_id            VARCHAR(50)   NOT NULL UNIQUE,
	user_type          TINYINT       NOT NULL DEFAULT 1
	COMMENT '0: deleted, 1: normal',
	balance            DECIMAL(9, 2) NOT NULL DEFAULT 0,
	income             JSON COMMENT '{cny: balance, ..}',
	currency           VARCHAR(10)   NOT NULL DEFAULT 'cny',
	credit_mode        TINYINT       NOT NULL DEFAULT 0
	COMMENT '0: close, 1: open',
	credit_amount      DECIMAL(8, 2),
	credit_duration    INT COMMENT 'unit: hour',
	first_in_debt_time TIMESTAMP     NULL,
	PRIMARY KEY (user_id)
);
