# Database Design

## Key Points

* Will use mysql as RDBM.
* Since the project is microservice-based, we use database-per-service pattern for the backend store.
* The throughput is expected not too high, so we will use private-tables-per-service or schema-per-service given the way is the lowest overhead. It means services access the tables data that belong to other services only by API.

## Schemas

Use GRANT to isolate tables for services

```sql
GRANT [type of permission] ON [database name].[table name] TO '[username]'@'%';
```

## Name Conventions

1. All lower case including table, column
2. Multiple words should be separated by underscores, i.e. [snake case](https://en.wikipedia.org/wiki/Snake_case).
3. Full words, not abbreviations; use common abbreviations for long word.
4. All singular names, not plural
5. Single column primary key fields should be named **id**.
6. Foreign key fields should be a combination of the name of the referenced table and the name of the referenced fields such as app_id.
7. No prefix or suffix such as tb_
8. Indexes should be explicitly named and include both the table name and the column name(s) indexed. 

## Database Scripts

```sql
CREATE DATABASE IF NOT EXISTS openpitrix
	DEFAULT CHARACTER SET utf8
	DEFAULT COLLATE utf8_general_ci
;
```

* For repo service

```sql
CREATE TABLE repo {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    url VARCHAR(255) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX repo_ix_name ON repo (name);
CREATE USER IF NOT EXISTS 'openpitrix_repo_user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON openpitrix.repo TO 'openpitrix_repo_user'@'%';
```

```sql
CREATE TABLE repo_label {
    repo_id VARCHAR(50) NOT NULL,
    label_key VARCHAR(50) NOT NULL,
    label_value VARCHAR(255) NOT NULL,
    PRIMARY KEY(repo_id, label_key),
    FOREIGN KEY(repo_id) REFERENCES repo(id) ON DELETE CASCADE
);
GRANT ALL PRIVILEGES ON openpitrix.repo_label TO 'openpitrix_repo_user'@'%';
```

```sql
CREATE TABLE repo_selector {
    repo_id VARCHAR(50) NOT NULL,
    selector_key VARCHAR(50) NOT NULL,
    selector_value VARCHAR(255) NOT NULL,
    PRIMARY KEY(repo_id, selector_key),
    FOREIGN KEY(repo_id) REFERENCES repo(id) ON DELETE CASCADE
);
GRANT ALL PRIVILEGES ON openpitrix.repo_selector TO 'openpitrix_repo_user'@'%';
```

* For app service
```sql
CREATE TABLE app {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    repo_id VARCHAR(50) NOT NULL,
    url VARCHAR(255) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX app_ix_name ON app (name);
CREATE USER IF NOT EXISTS 'openpitrix_app_user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON openpitrix.app TO 'openpitrix_app_user'@'%';
```

* For cluster service
```sql
CREATE TABLE cluster {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    app_id VARCHAR(50) NOT NULL,
    app_version VARCHAR(50) NOT NULL,
    app_version VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    transition_status VARCHAR(50) NOT NULL DEFAULT '',
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX cluster_ix_name ON cluster (name);
CREATE USER IF NOT EXISTS 'openpitrix_cluster_user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON openpitrix.cluster TO 'openpitrix_cluster_user'@'%';
```
```sql
CREATE TABLE cluster_node {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    instance_id VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    cluster_id VARCHAR(50) NOT NULL,
    private_ip VARCHAR(50) DEFAULT '' NOT NULL,
    FOREIGN KEY(cluster_id) REFERENCES cluster(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX cluster_node_ix_name ON clusternode (name);
GRANT ALL PRIVILEGES ON openpitrix.cluster_node TO 'openpitrix_cluster_user'@'%';
```

* For app runtime service
```sql
CREATE TABLE app_runtime {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    url VARCHAR(255) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX app_runtime_ix_name ON appruntime (name);
CREATE USER IF NOT EXISTS 'openpitrix_appruntime_user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON openpitrix.app_runtime TO 'openpitrix_appruntime_user'@'%';
```
```sql
CREATE TABLE app_runtime_label {
    app_runtime_id VARCHAR(50) NOT NULL,
    label_key VARCHAR(50) NOT NULL,
    label_value VARCHAR(255) NOT NULL,
    PRIMARY KEY(app_runtime_id, label_key),
    FOREIGN KEY(app_runtime_id) REFERENCES app_runtime(id) ON DELETE CASCADE
);
GRANT ALL PRIVILEGES ON openpitrix.app_runtime_label TO 'openpitrix_appruntime_user'@'%';
```
