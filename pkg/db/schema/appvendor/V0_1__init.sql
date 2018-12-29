/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2018/12/27 9:40:46                           */
/*==============================================================*/


alter table vendor_verify_info
   drop primary key;

drop table if exists vendor_verify_info;

/*==============================================================*/
/* Table: vendor_verify_info                                    */
/*==============================================================*/
create table vendor_verify_info
(
   user_id              varchar(50) not null,
   company_name         varchar(255) not null,
   company_website      varchar(255) not null,
   company_profile      text,
   authorizer_name      varchar(50) not null,
   authorizer_email     varchar(100) not null,
   authorizer_phone     varchar(50) not null,
   bank_name            varchar(200) not null,
   bank_account_name    varchar(50) not null,
   bank_account_number  varchar(100) not null,
   status               varchar(16) not null comment 'new,  submitted,  passed, rejected',
   reject_message       varchar(200) default '',
   approver             varchar(50),
   owner                varchar(50),
   owner_path           varchar(100),
   submit_time          timestamp,
   status_time          timestamp
);

alter table vendor_verify_info
   add primary key (user_id);



CREATE INDEX vendor_name_idx
	ON vendor_verify_info (company_name);

CREATE INDEX vendor_company_website_idx
	ON vendor_verify_info (company_website);

CREATE INDEX vendor_authorizer_name_idx
	ON vendor_verify_info (authorizer_name);

CREATE INDEX vendor_authorizer_email_idx
	ON vendor_verify_info (authorizer_email);

CREATE INDEX vendor_status_idx
	ON vendor_verify_info (status);

CREATE INDEX vendor_submit_time_idx
	ON vendor_verify_info (submit_time);

CREATE INDEX vendor_status_time_idx
	ON vendor_verify_info (status_time);




