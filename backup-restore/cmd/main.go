package main

import (
	"backup-restore/internal/backup"
	"backup-restore/internal/minio"
	"backup-restore/internal/restore"
	"flag"

	_ "github.com/lib/pq"

	"backup-restore/internal/config"
	"backup-restore/internal/db"
)

var (
	isBackup   bool
	isRestore  bool
	configPath string
)

func init() {
	flag.BoolVar(&isBackup, "backup", false, "Flag to init backup behavior")
	flag.BoolVar(&isRestore, "restore", false, "Flag to init restore behavior")
	flag.StringVar(&configPath, "config", "config/config.yaml", "Path to config file")
}

func main() {
	flag.Parse()

	appConfig, err := config.ParseConfig(configPath)
	if err != nil {
		panic(err)
	}

	provider, err := db.NewDataBaseProvider(appConfig)
	if err != nil {
		panic(err)
	}

	minio, err := minio.NewMinioClient(appConfig)
	if err != nil {
		panic(err)
	}

	if isBackup {
		/*if err = provider.AddNewProject(); err != nil {
			panic(err)
		}*/
		backup := backup.NewBackupManager(provider, minio)
		if err = backup.MakeDump(); err != nil {
			panic(err)
		}
	}

	if isRestore {
		restore := restore.NewRestoreManager(provider, minio)
		if err = restore.MakeRestore(); err != nil {
			panic(err)
		}
	}

}
