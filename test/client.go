package test

import (
	"flag"
	"log"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "openpitrix.io/openpitrix/test/client"
)

const UserSystem = "system"

type IgnoreLogger struct{}

func (IgnoreLogger) Printf(format string, args ...interface{}) {
}

func (IgnoreLogger) Debugf(format string, args ...interface{}) {
}

type ClientConfig struct {
	Host     string
	BasePath string
	Debug    bool
}

func GetClient(conf *ClientConfig) *apiclient.Openpitrix {
	transport := httptransport.New(conf.Host, conf.BasePath, []string{"http"})
	transport.SetDebug(conf.Debug)
	//transport.SetLogger(IgnoreLogger{})
	Client := apiclient.New(transport, strfmt.Default)
	return Client
}

func GetClientConfig() *ClientConfig {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	var (
		host     string
		basePath string
	)
	flag.StringVar(&host, "host", "localhost:9100", "specify api gateway host")
	flag.StringVar(&basePath, "base_path", "/", "specify http base path")
	flag.Parse()
	return &ClientConfig{
		Host:     host,
		BasePath: basePath,
		Debug:    false,
	}
}
