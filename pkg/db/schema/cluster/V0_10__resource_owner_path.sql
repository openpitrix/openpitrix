
ALTER TABLE cluster
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX cluster_owner_path_idx
	ON cluster (owner_path);
UPDATE cluster
SET owner_path = CONCAT(':', owner);

ALTER TABLE cluster_link
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX cluster_link_owner_path_idx
	ON cluster_link (owner_path);
UPDATE cluster_link
SET owner_path = CONCAT(':', owner);

ALTER TABLE cluster_node
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX cluster_node_owner_path_idx
	ON cluster_node (owner_path);
UPDATE cluster_node
SET owner_path = CONCAT(':', owner);

ALTER TABLE cluster_upgrade_audit
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX cluster_upgrade_audit_owner_path_idx
	ON cluster_upgrade_audit (owner_path);
UPDATE cluster_upgrade_audit
SET owner_path = CONCAT(':', owner);

ALTER TABLE key_pair
	ADD COLUMN owner_path VARCHAR(255) NOT NULL;
CREATE INDEX key_pair_owner_path_idx
	ON key_pair (owner_path);
UPDATE key_pair
SET owner_path = CONCAT(':', owner);
