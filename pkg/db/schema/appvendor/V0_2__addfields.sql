
alter table vendor_verify_info add approver varchar(50);
alter table vendor_verify_info add owner varchar(50);
alter table vendor_verify_info add owner_path varchar(100);

CREATE INDEX vendor_authorizer_email_idx
	ON vendor_verify_info (authorizer_email);

