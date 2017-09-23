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

```
create database apphub;
```

```
use database apphub;
create table repo {
    repoid varchar(50) PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(1000) DEFAULT ‘’ NOT NULL,
    url varchar(255) NOT NULL
);
CREATE UNIQUE INDEX repo_name ON repo (name);
CREATE USER 'repo_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.repo TO 'repo_user'@'localhost';
```

```
use database apphub;
create table app {
    appid varchar(50) PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(1000) DEFAULT ‘’ NOT NULL,
    repoid varchar(50) NOT NULL,
    url varchar(255) NOT NULL
);
CREATE UNIQUE INDEX app_name ON app (name);
CREATE USER 'app_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.app TO 'app_user'@'localhost';
```

```
use database apphub;
create table cluster {
    clusterid varchar(50) PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(1000) DEFAULT ‘’ NOT NULL,
    appid varchar(50) NOT NULL
);
CREATE UNIQUE INDEX cluster_name ON cluster (name);
CREATE USER 'cluster_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.cluster TO 'cluster_user'@'localhost';
```

```
use database apphub;
create table appruntime {
    appruntimeid varchar(50) PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(1000) DEFAULT ‘’ NOT NULL,
    url varchar(255) NOT NULL
);
CREATE UNIQUE INDEX appruntime_name ON appruntime (name);
CREATE USER 'appruntime_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.appruntime TO 'appruntime_user'@'localhost';
```
