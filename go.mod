module openpitrix.io/openpitrix

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Masterminds/semver v1.5.0
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/aws/aws-sdk-go v1.25.21
	github.com/bitly/go-simplejson v0.5.0
	github.com/chai2010/jsonmap v1.0.0
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/disintegration/imaging v1.6.1
	github.com/elazarl/goproxy v0.0.0-20200315184450-1f3cb6622dad // indirect
	github.com/fatih/camelcase v1.0.0
	github.com/fatih/structs v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.7
	github.com/go-openapi/spec v0.19.4
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gocraft/dbr v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	github.com/google/gops v0.3.6
	github.com/gorilla/websocket v1.4.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.11.3
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron v1.2.0
	github.com/sony/sonyflake v1.0.0
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.1
	github.com/xeipuuv/gojsonschema v1.2.0
	go.etcd.io/etcd v3.3.17+incompatible
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/tools v0.0.0-20191029041327-9cc4af7d6b2c
	google.golang.org/genproto v0.0.0-20191028173616-919d9bdd9fe6
	google.golang.org/grpc v1.24.0
	gopkg.in/square/go-jose.v2 v2.4.0
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/helm/v3 v3.0.1
	k8s.io/api v0.0.0-20200131112707-d64dbec685a4
	k8s.io/apiextensions-apiserver v0.0.0-20191016113550-5357c4baaf65
	k8s.io/apimachinery v0.0.0-20200131112342-0c9ec93240c9
	k8s.io/client-go v0.0.0-20200228043759-e3bfc0127563
	kubesphere.io/im v0.1.0
	openpitrix.io/iam v0.1.0
	openpitrix.io/notification v0.2.2
)

replace github.com/ugorji/go => github.com/ugorji/go v0.0.0-20190128213124-ee1426cffec0

replace github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0
