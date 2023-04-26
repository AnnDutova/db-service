package db

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq"

	"easy-master-slave/pkg/model"
)

type DataBaseProvider struct {
	DBMaster *sql.DB
	DBSlave  *sql.DB
	Config   *model.Config
}

func NewDataBaseProvider(config *model.Config) (*DataBaseProvider, error) {
	conMaster, err := openConnection(config.DbMaster.DbHost,
		config.DbMaster.DbPort,
		config.DbCommon.DbName,
		config.DbCommon.DbUsername,
		config.DbCommon.DbPassword)
	if err != nil {
		return nil, err
	}
	log.Printf("Open Master connection")

	conSlave, err := openConnection(config.DbSlave.DbHost,
		config.DbSlave.DbPort,
		config.DbCommon.DbName,
		config.DbCommon.DbUsername,
		config.DbCommon.DbPassword)
	if err != nil {
		return nil, err
	}
	log.Printf("Open Slave connection")

	return &DataBaseProvider{
		DBMaster: conMaster,
		DBSlave:  conSlave,
		Config:   config,
	}, nil
}

func openConnection(host, port, dbName, username, password string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
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

func (db *DataBaseProvider) Read() *sql.DB {
	return db.DBSlave
}

func (db *DataBaseProvider) Write() *sql.DB {
	return db.DBMaster
}
