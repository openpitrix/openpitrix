ALTER TABLE runtime_credential ADD COLUMN name VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE runtime_credential ADD COLUMN description TEXT NULL;
ALTER TABLE runtime_credential ADD COLUMN provider VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE runtime_credential ADD COLUMN runtime_url TEXT NULL;
ALTER TABLE runtime_credential ADD COLUMN owner VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE runtime_credential ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE runtime_credential ADD COLUMN status_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL;
ALTER TABLE runtime_credential CHANGE column content runtime_credential_content JSON NOT NULL;
CREATE INDEX runtime_credential_name_idx ON runtime_credential (name);
CREATE INDEX runtime_credential_provider_idx ON runtime_credential (provider);
CREATE INDEX runtime_credential_owner_idx ON runtime_credential (owner);
CREATE INDEX runtime_credential_status_idx ON runtime_credential (status);

UPDATE runtime_credential a INNER JOIN (SELECT runtime_credential_id,name,description,runtime_url,owner,provider,status,status_time FROM runtime) b
SET
a.name=b.name,
a.description=b.description,
a.runtime_url=b.runtime_url,
a.owner=b.owner,
a.provider=b.provider,
a.status=b.status,
a.status_time=b.status_time
WHERE a.runtime_credential_id=b.runtime_credential_id;

ALTER TABLE runtime drop column runtime_url;
DROP TABLE runtime_label;





