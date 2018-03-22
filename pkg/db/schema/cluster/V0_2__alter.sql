ALTER TABLE cluster
	CHANGE COLUMN vxnet_id subnet_id VARCHAR(50) NOT NULL;
ALTER TABLE cluster_node
	CHANGE COLUMN vxnet_id subnet_id VARCHAR(50) NOT NULL;
