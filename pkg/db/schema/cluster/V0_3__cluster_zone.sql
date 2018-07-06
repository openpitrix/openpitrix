ALTER TABLE cluster
	ADD COLUMN zone VARCHAR(50) NOT NULL;

CREATE INDEX cluster_zone_idx
	ON cluster (zone ASC);
