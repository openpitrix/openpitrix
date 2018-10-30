<p align="center"><a href="http://openpitrix.io" target="_blank"><img src="https://raw.githubusercontent.com/openpitrix/openpitrix/master/docs/images/logo.png" alt="OpenPitrix"></a></p>

# OpenPitrix

[![Build Status](https://travis-ci.org/openpitrix/openpitrix.svg)](https://travis-ci.org/openpitrix/openpitrix)
[![Docker Build Status](https://img.shields.io/docker/build/openpitrix/openpitrix.svg)](https://hub.docker.com/r/openpitrix/openpitrix/)
[![GoDoc](https://godoc.org/openpitrix.io/openpitrix?status.svg)](https://godoc.org/openpitrix.io/openpitrix)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/openpitrix/openpitrix/blob/master/LICENSE)

----

OpenPitrix 是一个开放的平台，致力于在多个云环境中(青云 QingCloud、AWS、kubernetes 等）开发和部署应用程序，从而能够让应用程序无缝的运行在各个云环境中。

Pitrix 的发音是 _['paitriks]_，它意指将 IaaS 和 Paas 纵横交错连接起来，从而能够让用户更加轻松的在多种运行环境中开发、部署和管理应用。即 Pitrix = **P**aaS + **I**aaS + Ma**trix**。同时它也有 PI （希腊语中的"π"）的含义，即包含无限应用的巨大矩阵。

----

## 背景

项目的灵感来自于[青云 QingCloud AppCenter](https://appcenter.qingcloud.com)，青云 QingCloud AppCenter 是一款帮助开发者快速创建企业级应用程序的平台，可以做到将项目周期缩短到按日来计算，且开发者可以在此对自己的产品进行销售。同时，开发的学习门槛是非常低的，遵照[开发者文档](https://appcenter-docs.qingcloud.com/developer-guide/)，通常花几个小时就能理解所有的工作流。但是，从产品发布以来，来自用户和合作伙伴的呼声——支持其它的IaaS平台——却越来越高，于是我们以开源协作的方式将平台开放，以支持正如 AWS、VMware、Kubernetes 等。

## 设计

设计的基本思路就是解耦应用程序的仓库和运行时环境。应用程序可以通过匹配运行时环境的标签和应用来源的仓库选择器来运行。有关项目的设计详情请移步[项目设计](docs/design/README.md)。

## 路线图

[OpenPitrix Roadmap](docs/Roadmap-zh.md) 是 OpenPitrix 核心开发团队预期的产品开发计划和功能列表，按照版本和角色模块进行功能的划分，详细说明了 OpenPitrix 的计划与未来走向以及愿景。OpenPitrix 将会逐步发展为多云环境下应用程序管理系统的全方位的解决方案。当然，如果您有什么更好的建议或意见，欢迎在 [Issues](https://github.com/openpitrix/openpitrix/issues) 中提出。

## 安装

OpenPitrix 支持 all-in-one 和 Kubernetes 两种安装模式，详见 [安装指南](https://docs.openpitrix.io/v0.3/zh-CN/installation/installation-guide/)。

## 快速入门

安装 OpenPitrix 后，[快速入门](https://docs.openpitrix.io/v0.3/zh-CN/getting-start/introduction/) 帮助您快速上手 OpenPitrix，了解 OpenPitrix 平台管理应用生命周期的基本工作流程。更多信息，请参阅 [OpenPitrix 官网](https://openpitrix.io) 和 [文档](https://docs.openpitrix.io)。

## 为项目做贡献

OpenPitrix 社区所有[成员](docs/members.md) 均必须遵守 [CNCF 行为准则](https://github.com/cncf/foundation/blob/master/code-of-conduct.md)，我们以为只有彼此尊重对方，构建高效的、协作的社区才有可能。

关于开发的说明，请移步[开发指南](docs/development.md).
