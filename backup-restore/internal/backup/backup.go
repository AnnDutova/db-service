package backup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"backup-restore/internal/db"
	"backup-restore/internal/minio"
)

type BackupManager struct {
	Provider    *db.DataBaseProvider
	MinioClient *minio.MinioClient
}

func NewBackupManager(provider *db.DataBaseProvider, minio *minio.MinioClient) *BackupManager {
	return &BackupManager{
		Provider:    provider,
		MinioClient: minio,
	}
}

func (m *BackupManager) MakeDump() error {
	containerName := m.Provider.Config.Db.ContainerName
	dbName := m.Provider.Config.Db.DbName
	dbPort := m.Provider.Config.Db.DbPort
	dbUser := m.Provider.Config.Db.DbUsername
	bucketName := m.Provider.Config.Minio.MinioBucket

	timeOfBackup := time.Now()
	timeStamp := fmt.Sprintf("%d-%d-%d_%d-%d", timeOfBackup.Day(), timeOfBackup.Month(), timeOfBackup.Year(),
		timeOfBackup.Hour(), timeOfBackup.Minute())
	outputFile := fmt.Sprintf("testdb_backup_%s.dump", timeStamp)

	// Build the command to run pg_dump
	cmd := exec.Command("docker", "exec", containerName, "pg_dump", "-d", dbName, "-p", dbPort, "-U", dbUser, "-wf", outputFile)
	err := cmd.Run()
	if err != nil {
		return err
	}

	log.Println("Backup succeeded!")

	cmd = exec.Command("docker", "exec", containerName, "cat", outputFile)

	outfile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmd.Stdout = outfile

	err = cmd.Run()
	if err != nil {
		return err
	}

	log.Println("Write to file")

	// Open the backup file
	file, err := os.Open(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload the file to MinIO
	err = m.MinioClient.AddFileToBucket(bucketName, outputFile, file)
	if err != nil {
		return err
	}

	log.Printf("Backup uploaded to MinIO bucket '%s' with object name '%s'\n", bucketName, outputFile)
	return nil
}
