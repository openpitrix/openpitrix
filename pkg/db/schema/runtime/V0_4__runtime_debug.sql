ALTER TABLE runtime ADD COLUMN debug BOOL NOT NULL DEFAULT 0;
ALTER TABLE runtime_credential ADD COLUMN debug BOOL NOT NULL DEFAULT 0;

CREATE INDEX runtime_debug_idx ON runtime (debug);
CREATE INDEX runtime_credential_debug_idx ON runtime_credential (debug);