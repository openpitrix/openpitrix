insert into repo (repo_id, name, description, url, credential, type, visibility, owner, status, create_time, status_time)
values ('repo-vmbased', 'default vmbased repo', '', 's3://minio.kubesphere-system.svc:9000/openpitrix-internal-repo/vmbased',
	'{\"access_key_id\": \"openpitrixminioaccesskey\", \"secret_access_key\": \"openpitrixminiosecretkey\"}', 's3',
	'public', 'system', 'active', '2018-01-01 00:00:00', '2018-01-01 00:00:00'),
	('repo-helm', 'default helm repo', '', 's3://minio.kubesphere-system.svc:9000/openpitrix-internal-repo/helm',
		'{\"access_key_id\": \"openpitrixminioaccesskey\", \"secret_access_key\": \"openpitrixminiosecretkey\"}', 's3',
		'public', 'system', 'active', '2018-01-01 00:00:00', '2018-01-01 00:00:00');

insert into repo_provider (repo_id, provider)
values ('repo-vmbased', 'aws'), ('repo-vmbased', 'qingcloud'), ('repo-helm', 'kubernetes');