<p align="center"><a href="http://openpitrix.io" target="_blank"><img src="https://raw.githubusercontent.com/openpitrix/openpitrix/master/docs/images/logo.png" alt="OpenPitrix"></a></p>

# OpenPitrix

[![Build Status](https://travis-ci.org/openpitrix/openpitrix.svg)](https://travis-ci.org/openpitrix/openpitrix)
[![Docker Build Status](https://img.shields.io/docker/build/openpitrix/openpitrix.svg)](https://hub.docker.com/r/openpitrix/openpitrix/)
[![Go Report Card](https://goreportcard.com/badge/openpitrix.io/openpitrix)](https://goreportcard.com/report/openpitrix.io/openpitrix)
[![GoDoc](https://godoc.org/openpitrix.io/openpitrix?status.svg)](https://godoc.org/openpitrix.io/openpitrix)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/openpitrix/openpitrix/blob/master/LICENSE)

----

OpenPitrix is an open platform to package and deploy applications into multiple cloud environments such as QingCloud, AWS, Kubernetes etc. Pitrix _['paitriks]_ means the matrix of PaaS and IaaS which makes it easy to develop, deploy, manage applications including PaaS on various runtime environments, i.e., Pitrix = **P**aaS + **I**aaS + Ma**trix**. It also means a matrix that contains endless (PI - the Greek letter "π") applications. 

----

## Motivation

The project originates from [QingCloud AppCenter](https://appcenter.qingcloud.com) which helps developers to create cloud-based enterprise applications in a few days and sell them on the center. In addition, the learning curve of how to [develop such applications](https://appcenter-docs.qingcloud.com/developer-guide/) is extremely low. Usually it takes a couple of hours for a developer to understand the working flow. Since QingCloud AppCenter was launched, many customers and partners have been asking us if it supports IaaS other than QingCloud such as AWS, Vmware. That is where the project comes from. Please read [OpenPitrix Insight](https://github.com/openpitrix/openpitrix/wiki/OpenPitrix-Insight) for details.

## Design

Basic idea is to decouple application repository and runtime environment. The runtime environment that an application can run is by matching the labels of runtime environment and the selectors of the repository where the application is from besides the provider. Please check out how to [design the project](docs/design/README.md).

## Roadmap

The [Roadmap](docs/Roadmap.md) demonstrates a list of open source product development plans and features being split by the edition and role modules, as well as our development team's anticipate of OpenPitrix. Obviously, it details the future's direction of OpenPitrix, but may change over time. We hope that can help you to get familiar with our open source plans and vision through the Roadmap. Of course, if you have any better ideas, welcome to [Issues](https://github.com/openpitrix/openpitrix/issues).

## Installation

Please follow the [Installation Guide](https://docs.openpitrix.io/v0.3/zh-CN/openpitrix-install-guide/) to install OpenPitrix.

## To start using OpenPitrix

To get started with OpenPitrix, please read the [Quick start guide](https://docs.openpitrix.io/v0.3/zh-CN/user-quick-start/).

For more information, please go to [openpitrix.io](http://openpitrix.io).

## Contributing to the project

All [members](docs/members.md) of the OpenPitrix community must abide by the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md). Only by respecting each other can we develop a productive, collaborative community.

You can then check out how to [setup for development](docs/development.md).
