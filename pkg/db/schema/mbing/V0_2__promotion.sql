/**
Promotion
**/

/** the attributes of all combination resources **/
CREATE TABLE IF NOT EXISTS combination_resource_attribute
(
	cra_id 									VARCHAR(50)		NOT NULL UNIQUE ,
	resource_attribute_ids 	JSON					NOT NULL COMMENT 'combination resource version id',
	create_time 						TIMESTAMP 	  DEFAULT CURRENT_TIMESTAMP,
	update_time 						TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								  VARCHAR(16) 	DEFAULT 'in_use' COMMENT 'in_use, deleted',
	PRIMARY KEY (cra_id)
);

/** the sku of combination resources **/
CREATE TABLE IF NOT EXISTS combination_sku
(
	com_sku_id 						VARCHAR(50) 	NOT NULL UNIQUE ,
	cra_id								VARCHAR(50)		NOT NULL COMMENT 'the id of combination_resource_attribute',
	attribute_values 			JSON 					NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {resource_version_id:{}, ..}',
	create_time 					TIMESTAMP 		DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status								VARCHAR(16) 	DEFAULT 'in_use' COMMENT 'in_use, deleted',
	PRIMARY KEY (com_sku_id)
);


CREATE TABLE IF NOT EXISTS combination_price
(
	com_price_id 						VARCHAR(50) 	NOT NULL UNIQUE ,
	com_sku_id 							VARCHAR(50) 	NOT NULL,
	resource_version_id 		VARCHAR(50) 	NOT NULL,
	attribute_id					 	VARCHAR(50) 	NOT NULL,
	prices 									JSON 					COMMENT '{attribute_value1: price1, ...}',
	currency            		VARCHAR(50)		NOT NULL  DEFAULT 'cny',
	create_time       			TIMESTAMP	    DEFAULT CURRENT_TIMESTAMP,
	update_time       			TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status									VARCHAR(16) 	DEFAULT 'in_use' COMMENT 'in_use, deleted',
	INDEX price_sku_index (com_price_id, com_sku_id),
	PRIMARY KEY (com_price_id)
);


/** probation sku of resource **/
CREATE TABLE IF NOT EXISTS probation_sku
(
	pro_sku_id 							VARCHAR(50) 	NOT NULL UNIQUE ,
	resource_attribute_id 	VARCHAR(50)		NOT NULL,
	attribute_values 				JSON 					NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {attribute: value, ...}',
	limit_num								INT						NOT NULL DEFAULT 1,
	create_time 						TIMESTAMP 		DEFAULT CURRENT_TIMESTAMP,
	update_time 						TIMESTAMP			DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status									VARCHAR(16) 	DEFAULT 'in_use' COMMENT 'in_use, deleted',
	PRIMARY KEY (pro_sku_id)
);


/** the records of probation resource used by user **/
CREATE TABLE IF NOT EXISTS probation_record
(
	probation_sku_id 				VARCHAR(50) 	NOT NULL,
	user_id 								VARCHAR(50) 	NOT NULL,
	num											INT						NOT NULL DEFAULT 1,
	create_time 						TIMESTAMP 		DEFAULT CURRENT_TIMESTAMP,
	probation_times 				JSON			 		COMMENT '[probation_time1, ...]',
	PRIMARY KEY (probation_sku_id, user_id)
);


CREATE TABLE IF NOT EXISTS dicount
(
	id 										VARCHAR(50) 		NOT NULL,
	name 									VARCHAR(255)		NOT NULL,
	limits								JSON 						COMMENT '{resource:.., sku:.., price:.., user:.., regoin:..}',
	discount_value 				DECIMAL(8, 2) 	COMMENT 'the price value to cut down',
	discount_percent 			DECIMAL(2, 2) 	COMMENT 'the price percent to cut down, there is only one of discount_value and discount_percent',
	start_time		  			TIMESTAMP 	  	DEFAULT CURRENT_TIMESTAMP,
	end_time       				TIMESTAMP   		DEFAULT CURRENT_TIMESTAMP,
	create_time     			TIMESTAMP	    	DEFAULT CURRENT_TIMESTAMP,
	status								VARCHAR(16) 	  DEFAULT 'in_use' COMMENT 'in_use, deleted, overtime',
	mark 									TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS coupon
(
	id 										VARCHAR(50) 		NOT NULL,
	name 									VARCHAR(50) 		NOT NULL,
	limits								JSON 						COMMENT '{resource:.., sku:.., price:.., user:.., regoin:...}',
	balance 							DECIMAL(8, 2) 	NOT NULL,
	count 								INT 						NOT NULL,
	limit_num 						INT 						NOT NULL DEFAULT 1,
	start_time		    		TIMESTAMP 	  	DEFAULT CURRENT_TIMESTAMP,
	end_time       				TIMESTAMP   		DEFAULT CURRENT_TIMESTAMP,
	create_time       		TIMESTAMP	    	DEFAULT CURRENT_TIMESTAMP,
	status								VARCHAR(16) 	  DEFAULT 'in_use' COMMENT 'in_use, deleted, overtime',
	mark 									TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS coupon_received
(
	id 										VARCHAR(50) 		NOT NULL,
	coupon_id 						VARCHAR(50) 		NOT NULL,
	user_id 							VARCHAR(50) 		NOT NULL,
	balance 							DECIMAL(8, 2) 	NOT NULL,
	status								VARCHAR(16) 	  DEFAULT 'received' COMMENT 'received, in_use, overtime',
	create_time       		TIMESTAMP	    	DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);
