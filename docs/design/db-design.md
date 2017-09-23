# Database Design

### Key Points
* Will use mysql as RDBM.
* Since the project is microservice-based, we use database-per-service pattern for the backend store.
* The throughput is expected not too high, so we will use private-tables-per-service and schema-per-service given the way is the lowest overhead.

### Schemas
Use GRANT to isolate tables for services
```
GRANT [type of permission] ON [database name].[table name] TO ‘[username]’@'localhost’;
```

```
create database apphub;
use database apphub;
create table repository {
    repoid varchar(50) PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(1000) DEFAULT ‘’ NOT NULL,
    url varchar(255) NOT NULL
);
CREATE UNIQUE INDEX repo_name ON repository (name);
CREATE USER 'repo_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON apphub.repository TO 'repo_user'@'localhost';
```
