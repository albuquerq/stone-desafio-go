package main

import (
	"net/http"

	"github.com/albuquerq/stone-desafio-go/pkg/application"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/persistence/mem"
	"github.com/albuquerq/stone-desafio-go/pkg/presentation/rest"
	"github.com/sirupsen/logrus"
)

func main() {

	appRegistry := application.NewRegistry(mem.NewRepositoryRegistry())

	handler := rest.New(appRegistry)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		logrus.Panic(err)
	}

}
