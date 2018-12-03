CREATE TABLE app_version_audit (
	version_id  VARCHAR(50) NOT NULL,
	app_id      VARCHAR(50) NOT NULL,
	status      VARCHAR(50) NOT NULL,
	operator    VARCHAR(50) NOT NULL,
	role        VARCHAR(50) NOT NULL,
	message     TEXT        NOT NULL,
	status_time TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX app_version_audit_version_id_idx
	ON app_version_audit (version_id);
CREATE INDEX app_version_audit_app_id_idx
	ON app_version_audit (app_id);
CREATE INDEX app_version_audit_status_idx
	ON app_version_audit (status);
CREATE INDEX app_version_audit_operator_idx
	ON app_version_audit (operator);
CREATE INDEX app_version_audit_status_time_idx
	ON app_version_audit (status_time);
