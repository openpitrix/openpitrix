# OpenPitrix

OpenPitrix is an open platform to package and deploy applications into multiple environment such as QingCloud, AWS, Kubernetes etc. 

### Motivation

The project originates from [QingCloud AppCenter](https://appcenter.qingcloud.com) which helps developers to create a cloud-based enterprise application in a few days and sell it on the center. In addition, the learning curve of how to [develop such application](https://appcenter-docs.qingcloud.com/developer-guide/) is extremely low. Usually it takes a couple of hours for a developer to understand the working flow. Since QingCloud AppCenter was launched, many customers and partners have been asking us if it supports IaaS other than QingCloud such as AWS, Vmware. That is where the project comes from. 

### Design

Basic idea is to decouple application repository and runtime environment. The runtime environment an application can run is by matching the labels of runtime environment and the selector of the repository where the application is from. Please check out how to [design the project](docs/design/README.md)

### Contributing to the project

All [members](docs/members.md) of the OpenPitrix community must abide by the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md). Only by respecting each other can we develop a productive, collaborative community.

You can then check out how to [setup for development](docs/development.md).
