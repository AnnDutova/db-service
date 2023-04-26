package apiserver

import (
	"easy-master-slave/internal/config"
	"easy-master-slave/internal/db"
	"net/http"
)

func Start(configPath string) error {
	appConfig, err := config.ParseConfig(configPath)
	if err != nil {
		return err
	}

	provider, err := db.NewDataBaseProvider(appConfig)
	if err != nil {
		return err
	}

	router := &Router{
		Provider: provider,
		Config:   appConfig,
	}

	configureRouters(router)

	return http.ListenAndServe(":8000", nil)
}
