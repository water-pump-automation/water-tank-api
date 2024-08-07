package webserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"water-tank-api/app/controllers"
	"water-tank-api/app/core/entity/logs"
	"water-tank-api/app/core/usecases/get_group"
	"water-tank-api/app/core/usecases/get_tank"
	mongodb "water-tank-api/infra/database/mongoDB"
	"water-tank-api/infra/logs/stdout"
	"water-tank-api/infra/web"
)

func External() {
	mainCtx := context.Background()

	logs.SetLogger(stdout.NewSTDOutLogger())

	mongoClient, err := mongodb.InitClient(mainCtx, databaseURI)
	if err != nil {
		logs.Gateway().Fatal(fmt.Sprintf("Error on starting mongo DB client: %s", err.Error()))
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}
	externalRouter := web.ExternalRouter{}
	collection := mongodb.NewCollection(mainCtx, mongoClient, databaseName, databaseCollection)

	externalRouter.Route(
		mux,
		controllers.NewExternalController(
			get_tank.NewGetWaterTank(collection),
			get_group.NewGetGroupWaterTank(collection),
		),
	)

	go func() {
		logs.Gateway().Info("Started internal server on port:" + serverPort)
		if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
			logs.Gateway().Fatal(fmt.Sprintf("Error on starting http listener: %s", err.Error()))
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	ctx, shutdownRelease := context.WithTimeout(mainCtx, 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(ctx); err != nil {
		logs.Gateway().Fatal(fmt.Sprintf("Shutdown error: %s", err.Error()))
	}
	logs.Gateway().Info("Graceful shutdown complete")
}
