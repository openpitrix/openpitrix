update vendor_verify_info set approver="" where approver is null;
alter table vendor_verify_info modify column approver varchar(50) not null default "" ;

update vendor_verify_info set owner="" where owner is null;
alter table vendor_verify_info modify column owner varchar(50) not null default "" ;

update vendor_verify_info set owner_path="" where owner_path is null;
alter table vendor_verify_info modify column owner_path varchar(100) not null default "" ;