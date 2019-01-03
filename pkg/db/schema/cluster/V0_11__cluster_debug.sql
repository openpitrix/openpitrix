ALTER TABLE cluster ADD COLUMN debug BOOL NOT NULL DEFAULT 0;
CREATE INDEX cluster_debug_index ON cluster (debug);
