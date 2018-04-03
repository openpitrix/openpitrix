ALTER TABLE cluster
	CHANGE COLUMN runtime provider VARCHAR(50) NOT NULL;
ALTER TABLE cluster
	DROP INDEX job_runtime_index;
ALTER TABLE cluster
	ADD INDEX job_provider_index (provider ASC);
