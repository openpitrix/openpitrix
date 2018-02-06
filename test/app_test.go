package test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/go-openapi/runtime"

	config "openpitrix.io/openpitrix/pkg/config/unittest"
	"openpitrix.io/openpitrix/test/client/app_service"
	"openpitrix.io/openpitrix/test/common"
	"openpitrix.io/openpitrix/test/models"
)

var (
	tShowEnvFlag = flag.Bool("show-env-flag", false, "show env flags")

	tConfig *config.Config
)

func TestMain(m *testing.M) {
	flag.Parse()

	if *tShowEnvFlag {
		config.PrintEnvs()
		os.Exit(0)
	}

	if conf, err := config.LoadConfig(); err == nil {
		tConfig = conf
	} else {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestApp(t *testing.T) {
	if !tConfig.Unittest.Rest.Enabled {
		t.Skip()
	}

	client := common.GetClient()
	// delete old app
	testAppId := "app-xxxxxxxy"
	testAppName := "test1"
	deleteAppParams := app_service.NewDeleteAppParams()
	deleteAppParams.SetID(testAppId)
	deleteAppRet, err := client.AppService.DeleteApp(deleteAppParams)
	if err != nil {
		if serr, ok := err.(*runtime.APIError); ok {
			if resp, ok := serr.Response.(runtime.ClientResponse); ok {
				t.Fatal(resp.Body())
			}
			t.Fatal(serr.Response)
		}
	}
	t.Log(deleteAppRet)
	// create new app
	appModel := models.OpenpitrixApp{ID: testAppId, Name: testAppName}
	createAppParams := app_service.NewCreateAppParams()
	createAppParams.WithBody(&appModel)
	createAppRet, err := client.AppService.CreateApp(createAppParams)
	if err != nil {
		if serr, ok := err.(*runtime.APIError); ok {
			if resp, ok := serr.Response.(runtime.ClientResponse); ok {
				t.Fatal(resp.Body())
			}
			t.Fatal(serr.Response)
		}
	}
	t.Log(createAppRet)
	// get new app
	getAppParams := app_service.NewGetAppParams()
	getAppParams.SetID(testAppId)
	getAppRet, err := client.AppService.GetApp(getAppParams)
	if err != nil {
		t.Fatal(err)
	}
	app := getAppRet.Payload
	if app == nil {
		t.Fatalf("failed to get app [%s]", testAppId)
	}
	t.Log(getAppRet)
	t.Log("test app finish, all test is ok")
}
