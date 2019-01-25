/**  Price  **/
CREATE TABLE IF NOT EXISTS attribute
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	name 									VARCHAR(255) 		NOT NULL,
	display_name 					VARCHAR(255),
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	remark  							TEXT,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS attribute_unit
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	name 									VARCHAR(30)		 	NOT NULL,
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
)

CREATE TABLE IF NOT EXISTS attribute_value
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	attribute_id 					VARCHAR(50) 		NOT NULL,
	attribute_unit_id 		VARCHAR(50) 		NOT NULL,
	min_value 						INT 						NOT NULL,
	max_value					  	INT 						NOT NULL COMMENT 'the attribute value, support scope of value [min_value, max_value);',
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS resource_attribute
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	resource_version_id 	VARCHAR(50)			NOT NULL,
	attributes 						JSON 						NOT NULL COMMENT 'sku attribute ids',
	metering_attributes 	JSON 						NOT NULL COMMENT 'the attribute ids need to metering',
	billing_attributes	 	JSON 						NOT NULL COMMENT 'the attribute ids for billing',
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS sku
(
	id 											VARCHAR(50) 	NOT NULL UNIQUE,
	resource_attribute_id 	VARCHAR(50)		NOT NULL,
	values 									JSON 					NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {attribute: value, ...}',
	create_time 						TIMESTAMP 		NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 						TIMESTAMP,
	status									TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS price
(
	id 											VARCHAR(50) 	NOT NULL UNIQUE,
	sku_id 									VARCHAR(50) 	NOT NULL,
	billing_attribute_id	 	VARCHAR(50)		NOT NULL,
	prices 									JSON 					commment '{attribute_value1: price1, ...}',
	currency            		VARCHAR(10)		NOT NULL  DEFAULT 'cny',
	create_time       			TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time       			TIMESTAMP,
	status									TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
	INDEX price_sku_index (sku_id, id),
	PRIMARY KEY (id)
);


/**  Metering  **/
CREATE TABLE IF NOT EXISTS leasing
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	group_id							VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	resource_id						VARCHAR(50)		NOT NULL COMMENT 'the same as cluster_id',
	sku_id								VARCHAR(50)		NOT NULL,
	metering_values				JSON					COMMENT 'the values of metering_attributes, {att_id: value, ..}',
	lease_time		    		TIMESTAMP 	  NOT NULL,
	renewal_time       		TIMESTAMP,
	update_time       		TIMESTAMP,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	close_time						JSON					COMMENT '[{close_time: restart_time}, ..]',
	status								TINYINT				NOT NULL DEFAULT 2 COMMENT '0: handClosed, 1: forceClosed, 2: running',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS leased
(
	leasing_id						VARCHAR(50)		NOT NULL UNIQUE,
	group_id							VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	resource_id						VARCHAR(50)		NOT NULL COMMENT 'the same as cluster_id',
	sku_id								VARCHAR(50)		NOT NULL,
	metering_values				JSON					COMMENT 'the values of metering_attributes, {att_id: value, ..}',
	lease_time		    		TIMESTAMP 	  NOT NULL,
	update_time       		TIMESTAMP	    NOT NULL,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	close_time						JSON					COMMENT '[{close_time: restart_time}, ..]',
	PRIMARY KEY (leasing_id)
);


/** Billing **/
CREATE TABLE IF NOT EXISTS leasing_contract
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	leasing_id						VARCHAR(50)		NOT NULL,
	sku_id								VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	metering_values				JSON					NOT NULL,
	start_time		    		TIMESTAMP 	  NOT NULL COMMENT 'same as leasing_time',
	update_time       		TIMESTAMP	    NOT NULL,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	fee_info							TEXT,
	fee										DECIMAL(8,2) 	NOT NULL COMMENT 'total fee from starting cluster to now',
	due_fee								DECIMAL(8,2) 	NOT NULL default 0,
	before_bill_fee				DECIMAL(8,2) 	NOT NULL DEFAULT 0 COMMENT 'the total fee of the before bills ',
	coupon_fee						DECIMAL(8,2) 	NOT NULL default 0,
	real_fee							DECIMAL(8,2) 	NOT NULL default 0,
	PRIMARY KEY (id)
)


CREATE TABLE IF NOT EXISTS leased_contract
(
	contract_id						VARCHAR(50)		NOT NULL UNIQUE,
	leasing_id						VARCHAR(50)		NOT NULL,
	sku_id								VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	metering_values				JSON					NOT NULL,
	start_time		    		TIMESTAMP,
	end_time		       		TIMESTAMP,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	fee_info							TEXT,
	fee										DECIMAL(8,2) 	NOT NULL COMMENT 'total fee from starting cluster to now',
	due_fee								DECIMAL(8,2) 	NOT NULL default 0,
	before_bill_fee				DECIMAL(8,2) 	NOT NULL DEFAULT 0 COMMENT 'the total fee of the before bills ',
	coupon_fee						DECIMAL(8,2) 	NOT NULL default 0,
	real_fee							DECIMAL(8,2) 	NOT NULL default 0,
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	PRIMARY KEY (contract_id)
)


/** Charge **/
CREATE TABLE IF NOT EXISTS charge
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	contract_id						VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	fee										DECIMAL(8,2)  NOT NULL COMMENT 'total fee from starting cluster to now',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	PRIMARY KEY (id)
)


CREATE TABLE IF NOT EXISTS recharge
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	user_id								VARCHAR(50)		NOT NULL,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	fee										DECIMAL(8,2)  NOT NULL COMMENT 'total fee from starting cluster to now',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	operator							VARCHAR(50)		NOT NULL,
	contract_id						VARCHAR(50),
	remark 								TEXT,
	PRIMARY KEY (id)
)


CREATE TABLE IF NOT EXISTS income
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	user_id								VARCHAR(50)		NOT NULL,
	contract_id						VARCHAR(50)		NOT NULL,
	balance								DECIMAL(9,2)	NOT NULL DEFAULT 0,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	PRIMARY KEY (id)
)


CREATE TABLE IF NOT EXISTS account
(
	user_id								VARCHAR(50)		NOT NULL UNIQUE,
	user_type							TINYINT				NOT NULL DEFAULT 1 COMMENT '0: deleted, 1: normal',
	balance								DECIMAL(9,2)	NOT NULL DEFAULT 0,
	income								DECIMAL(9,2)	NOT NULL DEFAULT 0,
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	credit_mode						TINYINT				NOT NULL DEFAULT  0 COMMENT '0: close, 1: open',
	credit_amount					DECIMAL(8, 2),
	credit_duration				DECIMAL(8, 2),
	first_in_debt_time		TIMESTAMP,
	PRIMARY KEY (user_id)
)

