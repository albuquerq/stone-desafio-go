package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/albuquerq/stone-desafio-go/pkg/application"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/persistence/mem"
	"github.com/albuquerq/stone-desafio-go/pkg/presentation/rest"
)

func main() {

	const addr = ":8080"

	appRegistry := application.NewRegistry(mem.NewRepositoryRegistry())

	handler := rest.New(appRegistry)

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			common.Logger().WithField("source", "main").Fatal(err)
		}
	}()
	common.Logger().Infof("Running on %s", server.Addr)

	<-done

	common.Logger().Info("Stoping server")

	err := server.Shutdown(context.Background())
	if err != nil {
		common.Logger().WithError(err).Error("Server shutdown failed")
	}

	common.Logger().Info("Server exited")

}
