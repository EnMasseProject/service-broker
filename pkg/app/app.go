package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/EnMasseProject/maas-service-broker/pkg/broker"
	"github.com/EnMasseProject/maas-service-broker/pkg/handler"
	"github.com/EnMasseProject/maas-service-broker/pkg/maas"
)

type App struct {
	broker   *broker.MaasBroker
	args     Args
	config   Config
	log      *Log
	client   *maas.MaasClient
}

func CreateApp() App {
	var err error

	fmt.Println("============================================================")
	fmt.Println("==           Starting MaaS Service Broker...              ==")
	fmt.Println("============================================================")

	app := App{}

	// Writing directly to stderr because log has not been bootstrapped
	if app.args, err = CreateArgs(); err != nil {
		os.Stderr.WriteString("ERROR: Failed to validate input\n")
		os.Stderr.WriteString(err.Error())
		ArgsUsage()
		os.Exit(127)
	}

	if app.config, err = CreateConfig(app.args.ConfigFile); err != nil {
		os.Stderr.WriteString("ERROR: Failed to read config file\n")
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	addressControllerHost := os.Getenv("ADDRESS_CONTROLLER_SERVICE_HOST")
	addressControllerPort := os.Getenv("ADDRESS_CONTROLLER_SERVICE_PORT")
	if addressControllerHost == "" || addressControllerPort == "" {
		fmt.Println("The following environment variables must point to the address controller: ADDRESS_CONTROLLER_SERVICE_HOST and ADDRESS_CONTROLLER_SERVICE_PORT")
		os.Exit(1)
	}
	app.config.Maas.Url = "http://" + addressControllerHost + ":" + addressControllerPort

	if app.log, err = NewLog(app.config.Log); err != nil {
		os.Stderr.WriteString("ERROR: Failed to initialize logger\n")
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	app.log.Debug("Connecting MaasClient")
	if app.client, err = maas.NewMaasClient(app.config.Maas, app.log.Logger); err != nil {
		app.log.Error("Failed to initialize MaasClient\n")
		app.log.Error(err.Error())
		os.Exit(1)
	}

	app.log.Debug("Creating MaaSBroker")
	if app.broker, err = broker.NewMaasBroker(app.log.Logger, app.client); err != nil {
		app.log.Error("Failed to create MaaSBroker\n")
		app.log.Error(err.Error())
		os.Exit(1)
	}

	return app
}

func (a *App) Start() {
	a.log.Notice("MaaS Service Broker Started")
	a.log.Notice("Listening on http://localhost:1338")
	err := http.ListenAndServe(":1338", handler.NewHandler(a.log.Logger, a.broker))
	if err != nil {
		a.log.Error("Failed to start HTTP server")
		a.log.Error(err.Error())
		os.Exit(1)
	}
}
