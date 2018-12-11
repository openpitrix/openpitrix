 
 

CREATE TABLE vendor_verify_info
(
	user_id varchar(50) NOT NULL,
	company_name varchar(255) NOT NULL,
	company_website varchar(255) NOT NULL,
	company_profile text,
	authorizer_name varchar(50) NOT NULL,
	authorizer_email varchar(100) NOT NULL,
	authorizer_phone varchar(50) NOT NULL, 
	bank_name varchar(200) NOT NULL,
	bank_account_name varchar(50) NOT NULL,
	bank_account_number varchar(100) NOT NULL,
	-- new，submitted,  passed, rejected
	status varchar(16) NOT NULL COMMENT 'new, submitted,  passed, rejected',
	reject_message  varchar(200) DEFAULT '',
	submit_time timestamp,
	status_time timestamp,
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (user_id)
);




CREATE INDEX vendor_name_idx
	ON vendor_verify_info (company_name);
	 
CREATE INDEX vendor_company_website_idx
	ON vendor_verify_info (company_website);

CREATE INDEX vendor_authorizer_name_idx
	ON vendor_verify_info (authorizer_name);

CREATE INDEX vendor_status_idx
	ON vendor_verify_info (status);

CREATE INDEX vendor_submit_time_idx
	ON vendor_verify_info (submit_time);

CREATE INDEX vendor_status_time_idx
	ON vendor_verify_info (status_time);

CREATE INDEX vendor_update_time_idx
	ON vendor_verify_info (update_time);



