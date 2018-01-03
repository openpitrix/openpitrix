# Change Log
All notable changes to QingCloud SDK for Go will be documented in this file.

## [v2.0.0-alpha.22] - 2017-12-23

### Add

- Add api to get vpn certificate

### Fixed

- Add loadbalancer attr cluster 
- POST query escape
- CreateServerCertificate use POST

## [v2.0.0-alpha.21] - 2017-12-15

### Fixed

- Change loadbalancer attr cluster to eips


## [v2.0.0-alpha.20] - 2017-12-15

### Fixed

- Fix loadbalancer http_header_size location


## [v2.0.0-alpha.19] - 2017-12-15

### Fixed

- Add loadbalancer http_header_size

## [v2.0.0-alpha.18] - 2017-12-14

### Fixed

- Add loadbalancer node_count


## [v2.0.0-alpha.17] - 2017-12-13

### Fixed

- Fix new config with endpoint

## [v2.0.0-alpha.16] - 2017-12-12

### Added

- Add new config with endpoint

### Fixed

- Fix cluster api bugs
- Fix userdata_path. userdata_file default value error 

## [v2.0.0-alpha.15] - 2017-12-07

### Added

- Add service cluster for go sdk

### Fixed

- Add requset type POST support 

## [v2.0.0-alpha.14] - 2017-11-24

### Fixed

- Add router static missing field
- Update template to fit snips change

## [v2.0.0-alpha.13] - 2017-11-09

### Fixed

- Delete DescribeSecurityGroupRules SecurityGroup Required

## [v2.0.0-alpha.12] - 2017-11-09

### Fixed

- DescribeSecurityGroupRules direction default value

## [v2.0.0-alpha.11] - 2017-11-08

### Fixed

- Fix wrong type of job id for stop lb output
- Add Makefile generate target to generate service code
- Fix ModifySecurityGroupRuleAttributes field type error
- Fix SecurityGroupRule field missing error

## [v2.0.0-alpha.10] - 2017-08-28

### Fixed

- Fixed loadbalancers section in request of StopLoadBalancers

## [v2.0.0-alpha.9] - 2017-08-23

### Added

- Add vxnetid and loadbalancertype parms for load balancer

### Fixed

- Fixed vxnet section in response of DescribeInstances

## [v2.0.0-alpha.8] - 2017-08-13

### Added

- Add missing parameter for describenic
- Add vxnet parms for instances

## [v2.0.0-alpha.7] - 2017-08-02

### Added

- Add timeout parameter for http client
- Add missing parameters in nic, router and security groups

## [v2.0.0-alpha.6] - 2017-07-17

### Added

- Update advanced client. [@jolestar]
- Fix several APIs. [@jolestar]

## [v2.0.0-alpha.5] - 2017-03-27

### Added

- Add advanced client for simple instance management. [@jolestar]
- Add wait utils for waiting a job to finish. [@jolestar]

## [v2.0.0-alpha.4] - 2017-03-14

### Fixed

- Fix Features type in RouterVxNet.

## [v2.0.0-alpha.3] - 2017-01-15

### Changed

- Fix request signer.

## [v2.0.0-alpha.2] - 2017-01-05

### Changed

- Fix logger output format, don't parse special characters.
- Rename package "errs" to "errors".

### Added

- Add type converters.

### BREAKING CHANGES

- Change value type in input and output to pointer.

## v2.0.0-alpha.1 - 2016-12-03

### Added

- QingCloud SDK for the Go programming language.
[v2.0.0-alpha.22]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.21...v2.0.0-alpha.22    
[v2.0.0-alpha.21]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.20...v2.0.0-alpha.21  
[v2.0.0-alpha.20]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.19...v2.0.0-alpha.20  
[v2.0.0-alpha.19]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.18...v2.0.0-alpha.19  
[v2.0.0-alpha.18]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.17...v2.0.0-alpha.18  
[v2.0.0-alpha.17]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.16...v2.0.0-alpha.17  
[v2.0.0-alpha.16]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.15...v2.0.0-alpha.16  
[v2.0.0-alpha.15]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.14...v2.0.0-alpha.15  
[v2.0.0-alpha.14]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.13...v2.0.0-alpha.14  
[v2.0.0-alpha.13]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.12...v2.0.0-alpha.13  
[v2.0.0-alpha.12]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.11...v2.0.0-alpha.12  
[v2.0.0-alpha.11]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.10...v2.0.0-alpha.11  
[v2.0.0-alpha.10]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.9...v2.0.0-alpha.10  
[v2.0.0-alpha.9]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.8...v2.0.0-alpha.9  
[v2.0.0-alpha.8]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.7...v2.0.0-alpha.8  
[v2.0.0-alpha.7]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.6...v2.0.0-alpha.7  
[v2.0.0-alpha.6]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.5...v2.0.0-alpha.6  
[v2.0.0-alpha.5]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.4...v2.0.0-alpha.5  
[v2.0.0-alpha.4]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.3...v2.0.0-alpha.4  
[v2.0.0-alpha.3]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.2...v2.0.0-alpha.3  
[v2.0.0-alpha.2]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.1...v2.0.0-alpha.2  

[@jolestar]: https://github.com/jolestar
