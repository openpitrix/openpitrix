ALTER TABLE cluster_role
	ADD COLUMN service_type VARCHAR(50) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN service_cluster_ip VARCHAR(50) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN service_external_ip VARCHAR(50) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN service_ports VARCHAR(50) NOT NULL DEFAULT "";

ALTER TABLE cluster_role
	ADD COLUMN config_map_data_count INT(11) NOT NULL DEFAULT 0;
ALTER TABLE cluster_role
	ADD COLUMN secret_data_count INT(11) NOT NULL DEFAULT 0;

ALTER TABLE cluster_role
	ADD COLUMN pvc_status VARCHAR(50) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN pvc_volume VARCHAR(100) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN pvc_capacity VARCHAR(50) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN pvc_access_modes VARCHAR(50) NOT NULL DEFAULT "";

ALTER TABLE cluster_role
	ADD COLUMN ingress_hosts VARCHAR(100) NOT NULL DEFAULT "";
ALTER TABLE cluster_role
	ADD COLUMN ingress_address VARCHAR(50) NOT NULL DEFAULT "";
