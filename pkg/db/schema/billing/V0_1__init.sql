/**
Price
**/
CREATE TABLE IF NOT EXISTS price (
	price_id    VARCHAR(50) NOT NULL UNIQUE,
	binding_id  VARCHAR(50) NOT NULL,
	prices      JSON COMMENT '{upto: price1, ...}',
	currency    VARCHAR(10) NOT NULL  DEFAULT 'cny',
	status      VARCHAR(16)           DEFAULT 'active'
	COMMENT 'active, deleted, disabled',
	start_time  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	end_time    TIMESTAMP   NULL,
	create_time TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (price_id)
);

/** Billing **/
CREATE TABLE IF NOT EXISTS leasing_contract (
	contract_id        VARCHAR(50)   NOT NULL UNIQUE,
	leasing_id         VARCHAR(50)   NOT NULL,
	resource_id        VARCHAR(50)   NOT NULL,
	sku_id             VARCHAR(50)   NOT NULL,
	user_id            VARCHAR(50)   NOT NULL,
	metering_values    JSON          NOT NULL,
	fee_info           TEXT,
	fee                DECIMAL(8, 2) NOT NULL
	COMMENT 'total fee from starting cluster to now',
	due_fee            DECIMAL(8, 2) NOT NULL default 0,
	other_contract_fee DECIMAL(8, 2) NOT NULL DEFAULT 0
	COMMENT 'the total fee of the other LeasingContracts splited',
	coupon_fee         DECIMAL(8, 2) NOT NULL default 0,
	real_fee           DECIMAL(8, 2) NOT NULL default 0,
	currency           VARCHAR(10)   NOT NULL DEFAULT 'cny',
	status             VARCHAR(16)   NOT NULL,
	create_time        TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	start_time         TIMESTAMP     NULL
	COMMENT 'same as leasing_time',
	status_time        TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (contract_id)
);


CREATE TABLE IF NOT EXISTS leased_contract (
	contract_id        VARCHAR(50)   NOT NULL UNIQUE,
	leasing_id         VARCHAR(50)   NOT NULL,
	resource_id        VARCHAR(50)   NOT NULL,
	sku_id             VARCHAR(50)   NOT NULL,
	user_id            VARCHAR(50)   NOT NULL,
	metering_values    JSON          NOT NULL,
	fee_info           TEXT,
	fee                DECIMAL(8, 2) NOT NULL
	COMMENT 'total fee from starting cluster to now',
	due_fee            DECIMAL(8, 2) NOT NULL default 0,
	other_contract_fee DECIMAL(8, 2) NOT NULL DEFAULT 0
	COMMENT 'the total fee of the other LeasingContracts splited',
	coupon_fee         DECIMAL(8, 2) NOT NULL default 0,
	real_fee           DECIMAL(8, 2) NOT NULL default 0,
	currency           VARCHAR(10)   NOT NULL DEFAULT 'cny',
	start_time         TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	end_time           TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	create_time        TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (contract_id)
);


CREATE TABLE IF NOT EXISTS account (
	user_id     VARCHAR(50)   NOT NULL UNIQUE,
	user_type   VARCHAR(16)   NOT NULL,
	balance     DECIMAL(9, 2) NOT NULL DEFAULT 0,
	currency    VARCHAR(10)   NOT NULL DEFAULT 'cny',
	income      JSON COMMENT '{cny: balance, ..}',
	status      VARCHAR(16)   NOT NULL,
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	status_time TIMESTAMP     NULL
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id)
);


/** Charge **/
CREATE TABLE IF NOT EXISTS charge (
	charge_id   VARCHAR(50)   NOT NULL UNIQUE,
	contract_id VARCHAR(50)   NOT NULL,
	user_id     VARCHAR(50)   NOT NULL,
	fee         DECIMAL(8, 2) NOT NULL
	COMMENT 'due_fee in contract',
	currency    VARCHAR(10)   NOT NULL DEFAULT 'cny',
	status      VARCHAR(16)            DEFAULT 'successful'
	COMMENT 'successful, failed',
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (charge_id)
);


CREATE TABLE IF NOT EXISTS refund (
	refund_id   VARCHAR(50)   NOT NULL UNIQUE,
	contract_id VARCHAR(50)   NOT NULL,
	fee         DECIMAL(8, 2) NOT NULL
	COMMENT 'due_fee in contract',
	currency    VARCHAR(10)   NOT NULL DEFAULT 'cny',
	status      VARCHAR(16)            DEFAULT 'successful'
	COMMENT 'successful, failed',
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (refund_id)
);


CREATE TABLE IF NOT EXISTS recharge (
	recharge_id VARCHAR(50)   NOT NULL UNIQUE,
	balance     DECIMAL(8, 2) NOT NULL
	COMMENT 'the money that recharged',
	currency    VARCHAR(10)   NOT NULL DEFAULT 'cny',
	status      VARCHAR(16)            DEFAULT 'successful'
	COMMENT 'successful, failed',
	create_time TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	description TEXT,
	PRIMARY KEY (recharge_id)
);


CREATE TABLE IF NOT EXISTS income (
	income_id                  VARCHAR(50)   NOT NULL UNIQUE,
	contract_id         VARCHAR(50)   NOT NULL,
	balance             DECIMAL(9, 2) NOT NULL DEFAULT 0,
	currency            VARCHAR(10)   NOT NULL DEFAULT 'cny',
	create_time         TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (income_id)
);
