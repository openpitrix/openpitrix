# OpenPitrix Roadmap

[ [中文版](Roadmap-zh.md) ]

OpenPitrix Roadmap demonstrates a list of open source product development plans and features being split by the edition and role modules, as well as OpenPitrix development team's anticipate of OpenPitrix. Obviously, it details the future's direction of OpenPitrix, but may change over time. We hope that can help you to get familiar with the project plans and vision through the Roadmap. Of course, if you have any better ideas, welcome to [Issues](https://github.com/openpitrix/openpitrix/issues).

## Release Goals

| Edition  | Schedule |
|---|---|
| [Release v0.1](Roadmap.md#v01)| May, 2018 |
| [Release v0.2](Roadmap.md#v02)| June, 2018 | 
| [Release v0.3](Roadmap.md#v03)| October, 2018 | 
| [Release v0.4](Roadmap.md#v04)| January, 2019 | 
| [Release v0.5](Roadmap.md#v05)| March, 2019 | 
| [Release v1.0](Roadmap.md#v10)| May, 2019 | 
| [Release v2.0](Roadmap.md#v20)| October, 2019 | 

## Features

## Past releases

### v0.1：

- [x] Resource Overview: display the resource statistics including runtimes, repositories and applications
- [x] Runtime: support for deploying applications to QingCloud, Kubernetes, KubeSphere platform.
- [x] Repository: support to create repositories for S3 or HTTP/HTTPS protocols like QingStor object storage.
- [x] Cluster Management: support cluster management, such as creating, disabling and closing cluster.
- [x] Store: support users to search and browse application.
- [x] Application Lifecycle Management: support to create and upload application packages based on Helm Chart, which can be deployed to Kubernetes as well.

### v0.2：

- [x] Runtime: add support for AWS and enable creating and closing clusters.
- [x] Repository: add support for AWS object storage of S3 protocol.
- [x] Application Lifecycle Management: add support for application package upload and deployment based on VM runtimes.
- [x] Adding key pairs which can attach ssh key to cluster node.

### v0.3：

**Admin**

- [x] Store: add support for getting application available for users to browse and deploy, add category management.
- [x] Application Lifecycle Management: add application review, application release, application deployment, application takedown, etc.
- [x] Platform Management: add support for repository, runtime and cluster instance management.
- [x] User Management: 3 roles by default, supporting role-based permission management.


**ISV**

- [x] Adding independent portal for ISV.
- [x] Application Lifecycle Management: add support for version management, such as uploading and creating new versions, as well as deploy to runtimes for new version.
- [x] Platform Management: support to view and manage repositories, runtimes, and cluster instances.


**User**

- [x] Store: support to view applications and details in the store. Also, application deployment will be supported to multi-cloud runtimes.
- [x] Runtime Management: support to manage runtimes such as QingCloud, AWS, Kubernetes, etc.   
- [x] Platform Management: support to view and manage applications and cluster instances, as well as view Pods of Helm cluster nodes.

## Upcoming releases

### v0.4:

**Admin**

- [ ] Adding Support for platform setting and quick guide, which can help admin to get started quickly.
- [ ] Resource Overview: support to view to-do items and detailed resource statistics.
- [ ] Runtime: Adding support for Ali Cloud.
- [ ] Application Lifecycle Management: refine application reviewing operations, such as supporting channel and development department.
- [ ] User management: Optimize the management of user roles
- [ ] Platform Management: add support for service monitoring, message center, as well as real-time resource and health status monitoring.
- [ ] System Settings: support users to customize the system theme such as name, color, Logo, etc.
- [ ] Adding ISV management, such as service provider audit.

**ISV**

- [ ] Application Management: add support for checking and editing application packages.
- [ ] Testing Management: sandbox will be added to provide test environment for all ISVs and isolate the production environment.
- [ ] Adding certification and contract management, such as manageing basic information and account of ISVs.


**User**

- [ ] Adding user guide function to help user get started quickly.


### v0.5:

**Admin**

- [ ] Resource Overview: add support to view ISV statistics and details.
- [ ] Store: support the customization of store theme such as Logo and color, etc.
- [ ] Platform Management: add support for message notification management, such as SMS and mail server configuration.
- [ ] System Settings: add support for platform built-in test environment.


**ISV**

- [ ] Adding operation and maintenance management, providing operation and maintenance dashboard, as well as event list and strategy management.

### v1.0：

**Admin**

- [ ] Platform Management: refine user identification management.
- [ ] System Settings: add support for billing by time, introducing platform terms.


**ISV**

- [ ] Application Management: refine application audit, such as ISV certification audit and contract audit.
- [ ] User List: support to view and manage user resource instances and application notifications.


### v2.0：

**Admin**

- [ ] System Settings: refine billing system, add support for multiple billing items such as capacity, number of calls, etc.


**ISV**

- [ ] Application Management: support to create more application types, such as SaaS, API, native applications, which can be deployed and tested in the runtimes.
- [ ] User List: introduce user ticket management.
- [ ] Member Management: add support for adding new members such as development, testing, finance, etc.

