# Database Design

### Key Points
* Will use mysql as RDBM.
* Since the project is microservice-based, we use database-per-service pattern for the backend store.
* The throughput is expected not too high, so we will use private-tables-per-service or schema-per-service given the way is the lowest overhead. It means services access the tables data that belong to other services only by API.

### Schemas
Use GRANT to isolate tables for services
```
GRANT [type of permission] ON [database name].[table name] TO ‘[username]’@'localhost’;
```

### Name Conventions
1. All lower case including table, column
2. Multiple words should be separated by underscores, i.e. [snake case](https://en.wikipedia.org/wiki/Snake_case).
3. Full words, not abbreviations; use common abbreviations for long word.
4. All singular names, not plural
5. Single column primary key fields should be named **id**.
6. Foreign key fields should be a combination of the name of the referenced table and the name of the referenced fields such as app_id.
7. No prefix or suffix such as tb_
8. Indexes should be explicitly named and include both the table name and the column name(s) indexed. 

### Database Scripts

```sql
CREATE DATABASE apphub;
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
CREATE USER 'repo_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.repo TO 'repo_user'@'localhost';
```
```sql
CREATE TABLE repo_label {
    repo_id VARCHAR(50) NOT NULL,
    label_key VARCHAR(50) NOT NULL,
    label_value VARCHAR(255) NOT NULL,
    PRIMARY KEY(repo_id, label_key),
    FOREIGN KEY(repo_id) REFERENCES repo(id) ON DELETE CASCADE
);
GRANT ALL PRIVILEGES ON apphub.repo_label TO 'repo_user'@'localhost';
```
```sql
CREATE TABLE repo_selector {
    repo_id VARCHAR(50) NOT NULL,
    selector_key VARCHAR(50) NOT NULL,
    selector_value VARCHAR(255) NOT NULL,
    PRIMARY KEY(repo_id, selector_key),
    FOREIGN KEY(repo_id) REFERENCES repo(id) ON DELETE CASCADE
);
GRANT ALL PRIVILEGES ON apphub.repo_selector TO 'repo_user'@'localhost';
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
CREATE USER 'app_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.app TO 'app_user'@'localhost';
```

* For cluster service
```sql
CREATE TABLE cluster {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    app_id VARCHAR(50) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX cluster_ix_name ON cluster (name);
CREATE USER 'cluster_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.cluster TO 'cluster_user'@'localhost';
```
```sql
CREATE TABLE cluster_node {
    id VARCHAR(50) PRIMARY KEY NOT NULL,
    instance_id VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(1000) DEFAULT ‘’ NOT NULL,
    cluster_id VARCHAR(50) NOT NULL
    FOREIGN KEY(cluster_id) REFERENCES cluster(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX cluster_node_ix_name ON clusternode (name);
GRANT ALL PRIVILEGES ON apphub.cluster_node TO 'cluster_user'@'localhost';
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
CREATE USER 'appruntime_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.app_runtime TO 'appruntime_user'@'localhost';
```
```sql
CREATE TABLE app_runtime_label {
    app_runtime_id VARCHAR(50) NOT NULL,
    label_key VARCHAR(50) NOT NULL,
    label_value VARCHAR(255) NOT NULL,
    PRIMARY KEY(app_runtime_id, label_key),
    FOREIGN KEY(app_runtime_id) REFERENCES app_runtime(id) ON DELETE CASCADE
);
GRANT ALL PRIVILEGES ON apphub.app_runtime_label TO 'appruntime_user'@'localhost';
```
