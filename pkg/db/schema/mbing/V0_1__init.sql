/**  Price  **/
CREATE TABLE IF NOT EXISTS attribute
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	name 									VARCHAR(255) 		NOT NULL,
	display_name 					VARCHAR(255),
	create_time 					TIMESTAMP 			DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP				DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	remark  							TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS attribute_unit
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	name 									VARCHAR(30)		 	NOT NULL,
	display_name 					VARCHAR(30)		 	NOT NULL,
	create_time 					TIMESTAMP 			DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP				DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS attribute_value
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	attribute_id 					VARCHAR(50) 		NOT NULL,
	attribute_unit_id 		VARCHAR(50),
	min_value 						INT 						NOT NULL,
	max_value					  	INT 						NULL COMMENT 'the attribute value, support scope of value (min_value, max_value]; NULL: max',
	create_time 					TIMESTAMP 			DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP				DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS resource_attribute
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE,
	resource_version_id 	VARCHAR(50)			NOT NULL,
	attributes 						JSON 						NOT NULL COMMENT 'sku attribute ids',
	metering_attributes 	JSON 						NOT NULL COMMENT 'the attribute ids need to metering and billing',
	create_time 					TIMESTAMP 			DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP				DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS sku
(
	id 										VARCHAR(50) 	NOT NULL UNIQUE,
	resource_attribute_id VARCHAR(50)		NOT NULL,
	attribute_values 			JSON 					NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {attribute: value, ...}',
	create_time 					TIMESTAMP 		DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS price
(
	id 										VARCHAR(50) 	NOT NULL UNIQUE,
	sku_id 								VARCHAR(50) 	NOT NULL,
	attribute_id	 				VARCHAR(50)		NOT NULL,
	prices 								JSON 					COMMENT '{attribute_value1: price1, ...}',
	currency            	VARCHAR(10)		NOT NULL  DEFAULT 'cny',
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	update_time       		TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
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
	other_info						VARCHAR(255)	COMMENT 'used for distinguish when resource_id and sku_id are same with others',
	metering_values				JSON					COMMENT 'the values of metering_attributes, {att_id: value, ..}',
	lease_time		    		TIMESTAMP			NULL,
	update_duration_time  TIMESTAMP			NULL,
	renewal_time       		TIMESTAMP			NULL,
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	update_time       		TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
	lease_time		    		TIMESTAMP			NULL,
	end_time       				TIMESTAMP			NULL,
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
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
	start_time		    		TIMESTAMP 		NULL COMMENT 'same as leasing_time',
	update_time       		TIMESTAMP			NULL,
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	fee_info							TEXT,
	fee										DECIMAL(8,2) 	NOT NULL COMMENT 'total fee from starting cluster to now',
	due_fee								DECIMAL(8,2) 	NOT NULL default 0,
	done_fee							DECIMAL(8,2) 	NOT NULL default 0,
	before_bill_fee				DECIMAL(8,2) 	NOT NULL DEFAULT 0 COMMENT 'the total fee of the before bills ',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	remark  							TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS leased_contract
(
	contract_id						VARCHAR(50)		NOT NULL UNIQUE,
	leasing_id						VARCHAR(50)		NOT NULL,
	sku_id								VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	metering_values				JSON					NOT NULL,
	start_time		    		TIMESTAMP			DEFAULT CURRENT_TIMESTAMP,
	end_time		       		TIMESTAMP			DEFAULT CURRENT_TIMESTAMP,
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	fee_info							TEXT,
	fee										DECIMAL(8,2) 	NOT NULL COMMENT 'total fee from starting cluster to now',
	due_fee								DECIMAL(8,2) 	NOT NULL default 0,
	done_fee							DECIMAL(8,2) 	NOT NULL default 0,
	before_bill_fee				DECIMAL(8,2) 	NOT NULL DEFAULT 0 COMMENT 'the total fee of the before bills ',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	remark  							TEXT,
	PRIMARY KEY (contract_id)
);


/** Charge **/
CREATE TABLE IF NOT EXISTS charge
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	user_id								VARCHAR(50)		NOT NULL,
	contract_id						VARCHAR(50)		NOT NULL,
	fee										DECIMAL(8,2)  NOT NULL COMMENT 'total fee from starting cluster to now',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	info									JSON					COMMENT '{couponReceivedID: fee}',
	status								TINYINT				DEFAULT 1 NOT NULL COMMENT '0: fail, 1: success',
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	remark  							TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS recharge
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	user_id								VARCHAR(50)		NOT NULL,
	contract_id						VARCHAR(50),
	fee										DECIMAL(8,2)  NOT NULL COMMENT 'total fee from starting cluster to now',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	info									JSON					COMMENT '{couponReceivedID: fee}',
	status								TINYINT				DEFAULT 1 NOT NULL COMMENT '0: fail, 1: success',
	operator							VARCHAR(50)		NOT NULL DEFAULT 'system',
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	remark 								TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS income
(
	id										VARCHAR(50)		NOT NULL UNIQUE,
	user_id								VARCHAR(50)		NOT NULL,
	contract_id						VARCHAR(50)		NOT NULL,
	resource_version_id		VARCHAR(50)		NOT NULL,
	balance								DECIMAL(9,2)	NOT NULL DEFAULT 0,
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	create_time       		TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	remark 								TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS account
(
	user_id								VARCHAR(50)		NOT NULL UNIQUE,
	user_type							TINYINT				NOT NULL DEFAULT 1 COMMENT '0: deleted, 1: normal',
	balance								DECIMAL(9,2)	NOT NULL DEFAULT 0,
	income								JSON					COMMENT '{cny: balance, ..}',
	currency            	VARCHAR(10)		NOT NULL DEFAULT 'cny',
	credit_mode						TINYINT				NOT NULL DEFAULT  0 COMMENT '0: close, 1: open',
	credit_amount					DECIMAL(8, 2),
	credit_duration				INT						COMMENT 'unit: hour',
	first_in_debt_time		TIMESTAMP			NULL,
	PRIMARY KEY (user_id)
);

