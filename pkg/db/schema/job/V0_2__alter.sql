ALTER TABLE job
	CHANGE COLUMN runtime provider VARCHAR(50) NOT NULL;
ALTER TABLE job
	DROP INDEX job_runtime_index;
ALTER TABLE job
	ADD INDEX job_provider_index (provider ASC);
