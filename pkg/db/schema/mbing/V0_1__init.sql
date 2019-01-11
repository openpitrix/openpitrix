/**  Metering  **/
CREATE TABLE IF NOT EXISTS leasing
(
	id						VARCHAR(50)		NOT NULL,
	resource_id				VARCHAR(50)		NOT NULL,
	resource_version_id		VARCHAR(50)		NOT NULL,
	user_id					VARCHAR(50)		NOT NULL,
	charge_id				VARCHAR(50)		NOT NULL,
	group_id				VARCHAR(50)		NOT NULL,
	duration				INT(11)			NOT NULL DEFAULT 0,
	lease_time		    	TIMESTAMP 	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	renewal_time       		TIMESTAMP   	NOT NULL,
	update_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	create_time       		TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);


/**  Price  **/
CREATE TABLE IF NOT EXISTS price
(
	id 					VARCHAR(50)		NOT NULL,
	resource_version_id VARCHAR(50) 	NOT NULL,
	charge_mode 		TINYINT			NOT NULL COMMENT 'ELASTIC:0, MONTHLY:1, YEARLY:2',
	price 				DECIMAL(8, 2) 	NOT NULL,
	duration 			INT(11),
	count 				INT(11),
	free_time 			INT(11),
	create_time       	TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	rule 				TEXT,
	INDEX price_resource_index (resource_version_id, charge_mode),
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS dicount
(
	id 					VARCHAR(50) 	NOT NULL,
	user_id				VARCHAR(50),
	name 				VARCHAR(255)	NOT NULL,
	price_id			VARCHAR(50) 	NOT NULL,
	new_price 			DECIMAL(8, 2),
	discount 			FLOAT(3, 2),
	type 				TINYINT		 	NOT NULL,
	start_time		    TIMESTAMP 	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	end_time       		TIMESTAMP   	NOT NULL,
	create_time       	TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	mark 				TEXT,
	PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS coupon
(
	id 					VARCHAR(50) 	NOT NULL,
	name 				VARCHAR(50) 	NOT NULL,
	sn 					VARCHAR(50),
	quota 				DECIMAL(8, 2) 	NOT NULL,
	balance 			DECIMAL(8, 2) 	NOT NULL,
	type 				TINYINT			NOT NULL COMMENT 'MONEY:0, TIME:1',
	resource_version_id VARCHAR(50),
	region				VARCHAR(50),
	status				TINYINT			NOT NULL COMMENT 'UNRECEIVE:0, RECEIVED:1, USED:2, EXPIRED:3',
	start_time		    TIMESTAMP 	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	end_time       		TIMESTAMP   	NOT NULL,
	create_time       	TIMESTAMP	    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	mark 				TEXT,
	PRIMARY KEY (id)
);
