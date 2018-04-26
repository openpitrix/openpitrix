CREATE TABLE IF NOT EXISTS cluster (
	cluster_id           VARCHAR(50)   NOT NULL,
	name                 VARCHAR(50)   NULL,
	description          VARCHAR(1000) NULL,
	app_id               VARCHAR(50)   NOT NULL,
	version_id           VARCHAR(50)   NOT NULL,
	subnet_id            VARCHAR(50)   NOT NULL,
	frontgate_id         VARCHAR(50)   NOT NULL,
	cluster_type         INT(11)       NOT NULL,
	endpoints            VARCHAR(1000) NULL,
	status               VARCHAR(50)   NOT NULL,
	transition_status    VARCHAR(50)   NOT NULL,
	create_time          TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time          TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	owner                VARCHAR(255)  NOT NULL,
	metadata_root_access BOOL          NOT NULL,
	global_uuid          MEDIUMTEXT    NOT NULL,
	upgrade_status       VARCHAR(50)   NOT NULL,
	upgrade_time         TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	runtime_id           VARCHAR(50)   NOT NULL,
	vpc_id               VARCHAR(50)   NOT NULL DEFAULT '',
	INDEX cluster_status_index (status ASC),
	INDEX cluster_create_time_index (create_time ASC),
	INDEX cluster_owner_index (owner ASC),
	PRIMARY KEY (cluster_id)
);

CREATE TABLE IF NOT EXISTS cluster_node (
	node_id           VARCHAR(50)    NOT NULL,
	name              VARCHAR(50)    NULL,
	cluster_id        VARCHAR(50)    NOT NULL,
	instance_id       VARCHAR(50)    NOT NULL,
	volume_id         VARCHAR(50)    NOT NULL,
	subnet_id         VARCHAR(50)    NOT NULL,
	private_ip        VARCHAR(50)    NOT NULL,
	server_id         INT(11)        NOT NULL,
	role              VARCHAR(50)    NOT NULL,
	status            VARCHAR(50)    NOT NULL,
	transition_status VARCHAR(50)    NOT NULL,
	device            VARCHAR(50)    NOT NULL DEFAULT '',
	create_time       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status_time       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
	owner             VARCHAR(255)   NOT NULL,
	group_id          INT(11)        NOT NULL,
	global_server_id  INT(11) UNIQUE NOT NULL AUTO_INCREMENT,
	custom_metadata   TEXT           NULL,
	is_backup         BOOL           NOT NULL DEFAULT 0,
	auto_backup       BOOL           NOT NULL DEFAULT 0,
	pub_key           TEXT           NULL,
	health_status     VARCHAR(50)    NOT NULL,
	PRIMARY KEY (node_id),
	INDEX cluster_node_cluster_id_index (cluster_id ASC),
	INDEX cluster_node_status_index (status ASC),
	INDEX cluster_node_create_time_index (create_time ASC),
	INDEX cluster_node_owner_index (owner ASC)
)
	AUTO_INCREMENT = 100000000;

CREATE TABLE IF NOT EXISTS cluster_common (
	cluster_id                   VARCHAR(50) NOT NULL,
	role                         VARCHAR(50) NOT NULL,
	server_id_upper_bound        INT(11)     NOT NULL,
	advanced_actions             TEXT        NULL,
	init_service                 TEXT        NULL,
	start_service                TEXT        NULL,
	stop_service                 TEXT        NULL,
	scale_out_service            TEXT        NULL,
	scale_in_service             TEXT        NULL,
	restart_service              TEXT        NULL,
	destroy_service              TEXT        NULL,
	upgrade_service              TEXT        NULL,
	custom_service               TEXT        NULL,
	health_check                 TEXT        NULL,
	monitor                      TEXT        NULL,
	passphraseless               TEXT        NULL,
	vertical_scaling_policy      VARCHAR(50) NOT NULL DEFAULT 'parallel',
	agent_installed              BOOL        NOT NULL,
	custom_metadata_script       TEXT        NULL,
	image_id                     TEXT        NOT NULL,
	backup_service               TEXT        NULL,
	backup_policy                VARCHAR(50) NULL,
	restore_service              TEXT        NULL,
	delete_snapshot_service      TEXT        NULL,
	incremental_backup_supported BOOL        NOT NULL DEFAULT 0,
	hypervisor                   VARCHAR(50) NOT NULL DEFAULT 'docker',
	PRIMARY KEY (cluster_id, role)
);

CREATE TABLE IF NOT EXISTS cluster_snapshot (
	snapshot_id        INT          NOT NULL,
	role               VARCHAR(50)  NOT NULL,
	server_ids         VARCHAR(255) NOT NULL,
	count              INT(11)      NOT NULL,
	app_id             VARCHAR(50)  NOT NULL,
	app_version        VARCHAR(50)  NOT NULL,
	child_snapshot_ids TEXT         NOT NULL,
	size               INT(11)      NOT NULL,
	PRIMARY KEY (snapshot_id, role, server_ids)
);

CREATE TABLE IF NOT EXISTS cluster_upgrade_audit (
	cluster_upgrade_audit_id VARCHAR(50) NOT NULL,
	cluster_id               VARCHAR(50) NOT NULL,
	from_app_version         VARCHAR(50) NOT NULL,
	to_app_version           VARCHAR(50) NOT NULL,
	service_params           TEXT        NULL,
	create_time              TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	upgrade_time             TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status                   VARCHAR(50) NOT NULL,
	owner                    VARCHAR(50) NOT NULL,
	PRIMARY KEY (cluster_upgrade_audit_id),
	INDEX cluster_upgrade_audit_cluster_index (cluster_id ASC),
	INDEX cluster_upgrade_audit_owner_index (owner ASC)
);

CREATE TABLE IF NOT EXISTS cluster_link (
	cluster_id          VARCHAR(50)  NOT NULL,
	name                VARCHAR(50)  NOT NULL,
	external_cluster_id VARCHAR(50)  NOT NULL,
	owner               VARCHAR(255) NOT NULL,
	PRIMARY KEY (cluster_id, name),
	INDEX cluster_link_name_index (name ASC),
	INDEX cluster_link_owner_index (owner ASC)
);

CREATE TABLE IF NOT EXISTS cluster_role (
	cluster_id    VARCHAR(50)  NOT NULL,
	role          VARCHAR(50)  NOT NULL,
	cpu           INT(11)      NOT NULL,
	gpu           INT(11)      NOT NULL,
	memory        INT(11)      NOT NULL,
	instance_size INT(11)      NOT NULL,
	storage_size  INT(11)      NOT NULL,
	env           TEXT         NULL,
	mount_point   VARCHAR(100) NOT NULL DEFAULT '',
	mount_options VARCHAR(100) NOT NULL DEFAULT '',
	file_system   VARCHAR(50)  NOT NULL DEFAULT '',
	PRIMARY KEY (cluster_id, role)
);

CREATE TABLE IF NOT EXISTS cluster_loadbalancer (
	cluster_id               INT         NOT NULL,
	role                     VARCHAR(50) NOT NULL,
	loadbalancer_listener_id VARCHAR(50) NOT NULL,
	loadbalancer_port        INT(11)     NOT NULL,
	loadbalancer_policy_id   VARCHAR(50) NOT NULL,
	PRIMARY KEY (cluster_id, role, loadbalancer_listener_id),
	INDEX cluster_loadbalancer_loadbalancer_listener_id_index (loadbalancer_listener_id ASC),
	INDEX cluster_loadbalancer_loadbalancer_policy_id_index (loadbalancer_policy_id ASC)
);
