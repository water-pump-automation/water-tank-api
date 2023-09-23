package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"water-tank-api/controllers"
	"water-tank-api/core/entity/logs"
	database_mock "water-tank-api/infra/database/mock"
	"water-tank-api/infra/stdout"
	"water-tank-api/infra/web"
	"water-tank-api/infra/web/routes"

	kingpin "github.com/alecthomas/kingpin/v2"
	iris "github.com/kataras/iris/v12"
	"golang.org/x/sync/errgroup"
)

var (
	port  = kingpin.Flag("port", "Server's port").Short('p').Default("8080").Envar("SERVER_PORT").Int()
	route = kingpin.Flag("vport", "Internal or External route").Short('r').Default("external").Envar("SERVER_ROUTE").String()
)

func main() {
	kingpin.Parse()

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logs.SetLogger(stdout.NewSTDOutLogger())

	app := iris.New()

	internalRouter := routes.ExternalRouter{}
	externalRouter := routes.InternalRouter{}

	switch r := strings.ToUpper(*route); r {
	case "INTERNAL":
		internalRouter.Route(app)
		web.SetControllers(controllers.NewInternalController(database_mock.NewWaterTankMockData()), nil)
		break
	case "EXTERNAL":
		web.SetControllers(nil, controllers.NewExternalController(database_mock.NewWaterTankMockData()))
		externalRouter.Route(app)
		break
	default:
		web.SetControllers(nil, controllers.NewExternalController(database_mock.NewWaterTankMockData()))
		externalRouter.Route(app)
		break
	}

	go func() {
		if err := app.Run(iris.Addr(fmt.Sprintf(":%d", *port))); err != nil {
			logs.Gateway().Fatal(fmt.Sprintf("Error on starting http listener: %s", err.Error()))
		}
	}()

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		<-gCtx.Done()

		app.Shutdown(context.Background())

		return nil
	})

	if err := g.Wait(); err != nil {
		logs.Gateway().Fatal(fmt.Sprintf("exit reason: %s \n", err))
	}
}
