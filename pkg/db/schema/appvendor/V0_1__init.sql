CREATE TABLE vendor_verify_info
(
	user_id varchar(50) NOT NULL,
	company_name varchar(255) NOT NULL,
	company_website varchar(255) NOT NULL,
	company_profile text NOT NULL,
	authorizer_name varchar(50) NOT NULL,
	authorizer_email varchar(100) NOT NULL,
	authorizer_phone varchar(50) NOT NULL,
	authorizer_position varchar(100),
	bank_name varchar(200) NOT NULL,
	bank_account_name varchar(50) NOT NULL,
	bank_account_number varchar(100) NOT NULL,
	
	remarks varchar(255),
	company_code varchar(255),
	-- personal, company
	verify_type varchar(20) COMMENT 'personal, company',
	business_license_code varchar(255),
	offical_fax varchar(50),
	offical_address varchar(255),
	zip_code varchar(50),	-- 三证合一;企业三证
	company_cert_type varchar(50)   COMMENT '三证合一;企业三证',
	bus_license_attchment_id varchar(50),
	taxreg_attachment_id varchar(50),
	comprep_attachment_id varchar(50),
	orgcode_attachment_id varchar(50),
	
	--  new, pending, passed, rejected
	status varchar(16) NOT NULL COMMENT ' new, pending, passed, rejected',
	reject_message  varchar(200),
	submit_time timestamp,
	status_time timestamp,
	update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (user_id)
);




CREATE INDEX vendor_name_idx
	ON vendor_verify_info (company_name);

CREATE INDEX vendor_status_idx
	ON vendor_verify_info (status);

CREATE INDEX vendor_submit_time_idx
	ON vendor_verify_info (submit_time);

CREATE INDEX vendor_status_time_idx
	ON vendor_verify_info (status_time);

CREATE INDEX vendor_update_time_idx
	ON vendor_verify_info (update_time);



