package db

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq"

	"backup-restore/pkg/model"
)

type DataBaseProvider struct {
	DB     *sql.DB
	Config *model.Config
}

func NewDataBaseProvider(config *model.Config) (*DataBaseProvider, error) {
	con, err := openConnection(config)
	if err != nil {
		return nil, err
	}

	return &DataBaseProvider{
		DB:     con,
		Config: config,
	}, nil
}

func openConnection(config *model.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db.DbHost, config.Db.DbPort, config.Db.DbUsername, config.Db.DbPassword, config.Db.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DataBaseProvider) AddNewProject() error {
	if _, err := db.DB.Query(`INSERT INTO project(id, title) VALUES (1,'Default'),(2, 'NewProjectBeforeBackup');`); err != nil {
		return err
	}
	log.Print("Add info")
	return nil
}
