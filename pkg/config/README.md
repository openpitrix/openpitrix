# config
`import "openpitrix.io/openpitrix/pkg/config"`

* [Overview](#pkg-overview)
* [Imported Packages](#pkg-imports)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

## <a name="pkg-imports">Imported Packages</a>

- [github.com/BurntSushi/toml](https://godoc.org/github.com/BurntSushi/toml)
- [github.com/golang/glog](https://godoc.org/github.com/golang/glog)
- [github.com/koding/multiconfig](https://godoc.org/github.com/koding/multiconfig)
- [github.com/pkg/errors](https://godoc.org/github.com/pkg/errors)
- [openpitrix.io/openpitrix/pkg/logger](https://godoc.org/openpitrix.io/openpitrix/pkg/logger)

## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func GetHomePath() string](#GetHomePath)
* [func IsUserConfigExists() bool](#IsUserConfigExists)
* [func RunInDocker() bool](#RunInDocker)
* [func UseDockerLinkedEnvironmentVariables()](#UseDockerLinkedEnvironmentVariables)
* [type ApiService](#ApiService)
* [type AppService](#AppService)
* [type ClusterService](#ClusterService)
* [type Config](#Config)
  * [func Default() \*Config](#Default)
  * [func Load(path string) (\*Config, error)](#Load)
  * [func MustLoad(path string) \*Config](#MustLoad)
  * [func MustLoadUnittestConfig() \*Config](#MustLoadUnittestConfig)
  * [func MustLoadUserConfig() \*Config](#MustLoadUserConfig)
  * [func Parse(content string) (\*Config, error)](#Parse)
  * [func (p \*Config) ActiveGlogFlags()](#Config.ActiveGlogFlags)
  * [func (p \*Config) Clone() \*Config](#Config.Clone)
  * [func (p \*Config) Save(path string) error](#Config.Save)
  * [func (p \*Config) String() string](#Config.String)
* [type Database](#Database)
  * [func (p \*Database) GetUrl() string](#Database.GetUrl)
* [type Glog](#Glog)
  * [func (p \*Glog) ActiveFlags()](#Glog.ActiveFlags)
* [type OpenPitrix_Config](#OpenPitrix_Config)
* [type RepoService](#RepoService)
* [type RuntimeService](#RuntimeService)
* [type Unittest](#Unittest)

#### <a name="pkg-files">Package files</a>
[config.go](./config.go) [default.go](./default.go) [docker.go](./docker.go) [glog.go](./glog.go) [unittest.go](./unittest.go) 

## <a name="pkg-constants">Constants</a>
``` go
const DefaultConfigContent = `
# OpenPitrix configuration
# https://openpitrix.io/

[Glog]
LogToStderr       = false
AlsoLogTostderr   = false
StderrThreshold   = "ERROR" # INFO, WARNING, ERROR, FATAL
LogDir            = ""

LogBacktraceAt    = ""
V                 = 0
VModule           = ""

CopyStandardLogTo = "INFO"

[DB]
Type         = "mysql"
Host         = "127.0.0.1"
Port         = 3306
Encoding     = "utf8"
Engine       = "InnoDB"
DbName       = "openpitrix"
RootPassword = "password"

[Api]
Host = "127.0.0.1"
Port = 9100

[App]
Host = "127.0.0.1"
Port = 9101

[Runtime]
Host = "127.0.0.1"
Port = 9102

[Cluster]
Host = "127.0.0.1"
Port = 9103

[Repo]
Host = "127.0.0.1"
Port = 9104

`
```
DefaultConfigContent is the default config file content.

``` go
const DefaultConfigPath = "~/.openpitrix/config.toml"
```
DefaultConfigFile is the default config file.

``` go
const UnittestConfigContent = `
# OpenPitrix configuration
# https://openpitrix.io/

[Api]
Host = "127.0.0.1"
Port = 9100

# Valid log levels are "debug", "info", "warn", "error", and "fatal".
LogLevel = "warn"

[DB]
Type         = "mysql"
Host         = "127.0.0.1"
Port         = 3306
Encoding     = "utf8"
Engine       = "InnoDB"
DbName       = "openpitrix"
RootPassword = "password"

[Unittest]
EnableDbTest = false

`
```
DefaultConfigContent is the default config file content.

``` go
const UnittestConfigPath = "~/.openpitrix/config_unittest.toml"
```
DefaultConfigFile is the default config file.

## <a name="GetHomePath">func</a> [GetHomePath](./config.go#L243)
``` go
func GetHomePath() string
```

## <a name="IsUserConfigExists">func</a> [IsUserConfigExists](./config.go#L258)
``` go
func IsUserConfigExists() bool
```

## <a name="RunInDocker">func</a> [RunInDocker](./docker.go#L15)
``` go
func RunInDocker() bool
```

## <a name="UseDockerLinkedEnvironmentVariables">func</a> [UseDockerLinkedEnvironmentVariables](./docker.go#L25)
``` go
func UseDockerLinkedEnvironmentVariables()
```

## <a name="ApiService">type</a> [ApiService](./config.go#L42-L45)
``` go
type ApiService struct {
    Host string `default:"127.0.0.1"`
    Port int    `default:"9100"`
}
```

## <a name="AppService">type</a> [AppService](./config.go#L47-L50)
``` go
type AppService struct {
    Host string `default:"127.0.0.1"`
    Port int    `default:"9101"`
}
```

## <a name="ClusterService">type</a> [ClusterService](./config.go#L57-L60)
``` go
type ClusterService struct {
    Host string `default:"127.0.0.1"`
    Port int    `default:"9103"`
}
```

## <a name="Config">type</a> [Config](./config.go#L25-L27)
``` go
type Config struct {
    OpenPitrix_Config
}
```

### <a name="Default">func</a> [Default](./config.go#L98)
``` go
func Default() *Config
```

### <a name="Load">func</a> [Load](./config.go#L111)
``` go
func Load(path string) (*Config, error)
```

### <a name="MustLoad">func</a> [MustLoad](./config.go#L130)
``` go
func MustLoad(path string) *Config
```

### <a name="MustLoadUnittestConfig">func</a> [MustLoadUnittestConfig](./config.go#L155)
``` go
func MustLoadUnittestConfig() *Config
```

### <a name="MustLoadUserConfig">func</a> [MustLoadUserConfig](./config.go#L142)
``` go
func MustLoadUserConfig() *Config
```

### <a name="Parse">func</a> [Parse](./config.go#L174)
``` go
func Parse(content string) (*Config, error)
```

### <a name="Config.ActiveGlogFlags">func</a> (\*Config) [ActiveGlogFlags](./glog.go#L26)
``` go
func (p *Config) ActiveGlogFlags()
```

### <a name="Config.Clone">func</a> (\*Config) [Clone](./config.go#L198)
``` go
func (p *Config) Clone() *Config
```

### <a name="Config.Save">func</a> (\*Config) [Save](./config.go#L214)
``` go
func (p *Config) Save(path string) error
```

### <a name="Config.String">func</a> (\*Config) [String](./config.go#L235)
``` go
func (p *Config) String() string
```

## <a name="Database">type</a> [Database](./config.go#L67-L75)
``` go
type Database struct {
    Type         string `default:"mysql"`
    Host         string `default:"127.0.0.1"`
    Port         int    `default:"3306"`
    Encoding     string `default:"utf8"`
    Engine       string `default:"InnoDB"`
    DbName       string `default:"openpitrix"`
    RootPassword string `default:"password"`
}
```

### <a name="Database.GetUrl">func</a> (\*Database) [GetUrl](./config.go#L94)
``` go
func (p *Database) GetUrl() string
```

## <a name="Glog">type</a> [Glog](./config.go#L81-L92)
``` go
type Glog struct {
    LogToStderr     bool   `default:"false"`
    AlsoLogTostderr bool   `default:"false"`
    StderrThreshold string `default:"ERROR"` // INFO, WARNING, ERROR, FATAL
    LogDir          string `default:""`

    LogBacktraceAt string `default:""`
    V              int    `default:"0"`
    VModule        string `default:""`

    CopyStandardLogTo string `default:""`
}
```

### <a name="Glog.ActiveFlags">func</a> (\*Glog) [ActiveFlags](./glog.go#L30)
``` go
func (p *Glog) ActiveFlags()
```

## <a name="OpenPitrix_Config">type</a> [OpenPitrix_Config](./config.go#L29-L40)
``` go
type OpenPitrix_Config struct {
    Glog Glog

    DB      Database
    Api     ApiService
    App     AppService
    Runtime RuntimeService
    Cluster ClusterService
    Repo    RepoService

    Unittest Unittest
}
```

## <a name="RepoService">type</a> [RepoService](./config.go#L62-L65)
``` go
type RepoService struct {
    Host string `default:"127.0.0.1"`
    Port int    `default:"9104"`
}
```

## <a name="RuntimeService">type</a> [RuntimeService](./config.go#L52-L55)
``` go
type RuntimeService struct {
    Host string `default:"127.0.0.1"`
    Port int    `default:"9102"`
}
```

## <a name="Unittest">type</a> [Unittest](./config.go#L77-L79)
``` go
type Unittest struct {
    EnableDbTest bool `default:"false"`
}
```

- - -
Generated by [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)