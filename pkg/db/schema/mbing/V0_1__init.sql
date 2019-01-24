/**  Price  **/
CREATE TABLE IF NOT EXISTS attribute
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE ,
	name 									VARCHAR(255) 		NOT NULL,
	display_name 					VARCHAR(255),
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	mark 									TEXT,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS attribute_unit
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE ,
	name 									VARCHAR(30)		 	NOT NULL,
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
)

CREATE TABLE IF NOT EXISTS attribute_value
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE ,
	attribute_id 					VARCHAR(50) 		NOT NULL,
	attribute_unit_id 		VARCHAR(50) 		NOT NULL,
	min_value 						INT 						NOT NULL ,
	max_value					  	INT 						NOT NULL COMMENT 'the attribute value, support scope of value [min_value, max_value);',
	create_time 					TIMESTAMP 			NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT 				DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS resource_attribute
(
	id 										VARCHAR(50) 		NOT NULL UNIQUE ,
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
	id 											VARCHAR(50) 	NOT NULL UNIQUE ,
	resource_attribute_id 	VARCHAR(50)		NOT NULL,
	values 									JSON 					NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {attribute: value, ...}',
	create_time 						TIMESTAMP 		NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 						TIMESTAMP,
	status									TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS price
(
	id 											VARCHAR(50) 	NOT NULL UNIQUE ,
	sku_id 									VARCHAR(50) 	NOT NULL,
	billing_attribute_id	 	VARCHAR(50)		NOT NULL ,
	prices 									JSON 					commment '{attribute_value1: price1, ...}',
	currency            		VARCHAR(50)		NOT NULL  DEFAULT 'cny',
	create_time       			TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time       			TIMESTAMP,
	status									TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
	INDEX price_sku_index (sku_id, id),
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS combination_resource_attribute
(
	id 										VARCHAR(50) NOT NULL UNIQUE ,
	resource_version_ids 	JSON NOT NULL COMMENT 'combination resource version id',
	attributes 						JSON NOT NULL COMMENT 'sku attribute ids: {resource_version_id:[],...}',
	metering_attributes 	JSON NOT NULL COMMENT 'the attribute ids need to metering: {resource_version_id:{}, ..}',
	billing_attributes	 	JSON NOT NULL COMMENT 'the attribute ids for billing: {resource_version_id:{}, ..}',
	create_time 					TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 					TIMESTAMP,
	status								TINYINT DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS combination_sku
(
	id 											VARCHAR(50) NOT NULL UNIQUE ,
	resource_attribute_id 	VARCHAR(50)		NOT NULL,
	values 									JSON NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {resource_version_id:{}, ..}',
	create_time 						TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 						TIMESTAMP,
	status									TINYINT DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS combination_price
(
	id 											VARCHAR(50) 	NOT NULL UNIQUE ,
	combination_sku_id 			VARCHAR(50) 	NOT NULL,
	billing_attribute_id	 	JSON					COMMENT '{resource_version_id: attribute_id}' ,
	prices 									JSON 					COMMENT '{attribute_value1: price1, ...}',
	currency            		VARCHAR(50)		NOT NULL  DEFAULT 'cny',
	create_time       			TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time       			TIMESTAMP,
	status									TINYINT 			DEFAULT 1 COMMENT '1: using, 0: deleted',
	INDEX price_sku_index (combination_sku_id, id),
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS probation_sku
(
	id 											VARCHAR(50) NOT NULL UNIQUE ,
	resource_attribute_id 	VARCHAR(50)		NOT NULL,
	values 									JSON NOT NULL COMMENT 'sku attribute values for attributes in resource_attribute: {attribute: value, ...}',
	create_time 						TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	update_time 						TIMESTAMP,
	status									TINYINT DEFAULT 1 COMMENT '1: using, 0: deleted',
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS probation_record
(
	probation_sku_id 				VARCHAR(50) 	NOT NULL,
	user_id 								VARCHAR(50) 	NOT NULL,
	create_time 						TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (probation_id, user_id)
)


CREATE TABLE IF NOT EXISTS dicount
(
	id 									VARCHAR(50) 	NOT NULL,
	name 								VARCHAR(255)	NOT NULL,
	limit								JSON 					COMMENT '{resource:.., sku:.., price:.., user:..,}',
	discount_value 			DECIMAL(8, 2) COMMENT 'the price value to cut down',
	discount_percent 		DECIMAL(1, 2) COMMENT 'the price percent to cut down, there is only one of discount_value and discount_percent',
	start_time		  		TIMESTAMP 	  NOT NULL,
	end_time       			TIMESTAMP   	NOT NULL,
	create_time     		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	mark 								TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS coupon
(
	id 									VARCHAR(50) 	NOT NULL,
	name 								VARCHAR(50) 	NOT NULL,
	limit								JSON 					COMMENT '{resource:.., sku:.., price:.., user:.., regoin:...}',
	balance 						DECIMAL(8, 2) NOT NULL,
	count 							INT 					NOT NULL,
	limit_num 					INT 					NOT NULL DEFAULT 1,
	start_time		    	TIMESTAMP 	  NOT NULL,
	end_time       			TIMESTAMP   	NOT NULL,
	create_time       	TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	mark 								TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS coupon_received
(
	coupon_id 					VARCHAR(50) 	NOT NULL,
	user_id 						VARCHAR(50) 	NOT NULL,
	num 								INT 					NOT NULL ,
	left_num						INT 					NOT NULL ,
	status							TINYINT				NOT NULL commnet '0: overtime, 1: received, 2: using, 3: used',
	create_time       	TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (coupon_id, user_id)
);





/**  Metering  **/
CREATE TABLE IF NOT EXISTS leasing
(
	id										VARCHAR(50)		NOT NULL,
	resource_id						VARCHAR(50)		NOT NULL,
	resource_version_id		VARCHAR(50)		NOT NULL,
	user_id								VARCHAR(50)		NOT NULL,
	price_id							VARCHAR(50)		NOT NULL,
	group_id							VARCHAR(50)		NOT NULL,
	duration							INT(11)			NOT NULL DEFAULT 0,
	lease_time		    		TIMESTAMP 	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	renewal_time       		TIMESTAMP   	NOT NULL,
	update_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status								VARCHAR(50),
	PRIMARY KEY (id)
);