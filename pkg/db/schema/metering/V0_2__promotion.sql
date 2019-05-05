/**
Promotion
**/

/** the combination **/
CREATE TABLE IF NOT EXISTS combination (
	combination_id VARCHAR(50)  NOT NULL UNIQUE,
	name           VARCHAR(255) NOT NULL,
	description    TEXT,
	owner          VARCHAR(50)  NOT NULL,
	status         VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted',
	create_time    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (combination_id)
);


/** the combination_sku **/
CREATE TABLE IF NOT EXISTS combination_sku (
	combination_sku_id VARCHAR(50) NOT NULL UNIQUE,
	combination_id     VARCHAR(50) NOT NULL,
	sku_id             VARCHAR(50) NOT NULL,
	status             VARCHAR(16) DEFAULT 'active'
	COMMENT 'active, deleted',
	create_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
	status_time        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
	ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (combination_sku_id)
);
