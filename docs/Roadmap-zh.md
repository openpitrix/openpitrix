# OpenPitrix Roadmap

OpenPitrix Roadmap 是 OpenPitrix 核心开发团队预期的产品开发计划和功能列表，按照角色和版本模块进行功能的划分，详细说明了 OpenPitrix 的未来走向，后续可能会随时间出现变动。我们希望通过 Roadmap 能够让您知悉我们的计划与愿景。当然，如果您有什么更好的意见，欢迎在 [Issues](https://github.com/openpitrix/openpitrix/issues) 中提出。

## 发布

- [Alpha 0.1](https://github.com/FeynmanZhou/openpitrix/blob/master/docs/Roadmap-zh.md#alpha-01)：2018 年 7 月
- [Beta 0.2](https://github.com/FeynmanZhou/openpitrix/blob/master/docs/Roadmap-zh.md#beta-02)：2018 年 7 月
- [Beta 0.3](https://github.com/FeynmanZhou/openpitrix/blob/master/docs/Roadmap-zh.md#beta-03)：2018 年 9 月底
- [v1.0](https://github.com/FeynmanZhou/openpitrix/blob/master/docs/Roadmap-zh.md#v10)：2019 年 1 月
- [v2.0](https://github.com/FeynmanZhou/openpitrix/blob/master/docs/Roadmap-zh.md#v20)：2019 年 6 月
- v3.0 2019 年 12 月

## 功能

### Alpha 0.1：

- [x] 用户资源概览
- [x] 应用仓库，支持创建 S3 或 http/https 协议的仓库
- [x] 运行环境支持 QingCloud，能在其部署应用、创建和关闭集群
- [x] 基于 Kubernetes 运行环境的 Helm 应用的创建、更新参数等
- [x] 应用商店，支持 ISV 上传应用通过审核后上架到应用商店。

### Beta 0.2：

- [x] 运行环境支持 AWS
- [x] 完成应用审核、集群管理、应用仓库、运行环境等核心功能

### Beta 0.3：
- [ ] 应用商店，支持 ISV 上传应用通过审核后上架到应用商店，供用户浏览和部署使用，以及应用分类
- [ ] 应用生命周期管理，如应用上传、应用发布、应用部署、版本控制、升级、应用下架等
- [ ] 应用部署支持 QingCloud、AWS、Kubernetes 等运行环境
- [ ] 优化应用审核、集群管理、应用仓库、运行环境等核心功能
- [ ] 用户管理、权限管理


### v1.0：

管理员

- [ ] 支持平台设置、快速引导，帮助管理员快速上手
- [ ] 工作台支持查看资源概览、待办事项
- [ ] ISV 管理，如服务商审核、合约管理
- [ ] 运行环境，引入阿里云、VMware、OpenStack、KubeSphere，支持对运行环境的管理
- [ ] 商店管理，如应用审核、应用预览、应用部署、分类管理、应用目录管理以及商店主题定制化
- [ ] 用户管理，如权限管理、用户角色分配、组织机构配置
- [ ] 财务管理，如财务报表、平台对账、消费明细、提现管理等功能
- [ ] 平台管理，如服务监控、消息中心，实时监控系统的资源健康状态
- [ ] 身份验证管理，对接 LDAP/AD，后续将支持第三方登录如 Github、微信
- [ ] 计量计费，支持一次性付费和订阅付费等方式；平台支付方式支持线上和线下
- [ ] 系统主题设置，名称、颜色、Logo 支持用户自定义

软件提供商

- [ ] 支持快速创建应用，支持 Helm 和 VM 类型的应用，并支持部署和测试
- [ ] 应用管理，如版本管理、应用提交审核
- [ ] 应用运维，提供运维看板、事件列表和策略管理等
- [ ] 沙箱，为所有 ISV 提供一个沙箱测试环境，并给出实例和用户列表
- [ ] 财务管理，如销售明细、平台对账、提现申请
- [ ] 认证与合约管理，如管理服务商基本信息、银行账户、合约。

用户

- [ ] 应用商店，查看应用商店所有应用及应用详情
- [ ] 支持应用部署到运行环境
- [ ] 环境管理，如添加 QingCloud、AWS、Kubernetes 这类运行环境，管理已部署的实例和应用
- [ ] 我的钱包，包括消费记录和账户余额
- [ ] 已购应用，支持实例的弹性伸缩配置和针对实例设置监控告警

## v2.0：

管理员

- [ ] 用户管理支持添加和管理用户组
- [ ] 财务管理支持发票管理
- [ ] 平台管理引入工单系统
- [ ] 身份验证管理支持第三方登录如 Github、微信
- [ ] 运行环境引入 EdgeWize，支持短信服务器配置

软件提供商

- [ ] 支持创建更多的应用类型如 SaaS 类、API 类、原生类应用，下一个版本将支持 Serverless 类应用和系统类应用
- [ ] 支持应用编排
- [ ] 集成 CI/CD，支持开发者构建 Pipeline
- [ ] 引入用户工单系统
- [ ] 营销管理，如创建促销活动和优惠券活动
- [ ] 成员管理，如添加开发、测试、财务等成员

用户

- [ ] 支持优惠券和发票管理

