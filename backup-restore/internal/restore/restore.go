package restore

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"backup-restore/internal/db"
	"backup-restore/internal/minio"
)

type RestoreManager struct {
	Provider    *db.DataBaseProvider
	MinioClient *minio.MinioClient
}

func NewRestoreManager(db *db.DataBaseProvider, minio *minio.MinioClient) *RestoreManager {
	return &RestoreManager{
		Provider:    db,
		MinioClient: minio,
	}
}

func (m *RestoreManager) MakeRestore() error {
	bucketName := m.Provider.Config.Minio.MinioBucket
	objectName := "testdb_buckup.dump"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	object, err := m.MinioClient.GetFileFromBucket(ctx, bucketName)
	defer object.Close()

	backupFile, err := os.Create(objectName)
	if err != nil {
		return fmt.Errorf("Failed to create local file: %v", err)
	}
	defer backupFile.Close()

	if _, err := io.Copy(backupFile, object); err != nil {
		return fmt.Errorf("Failed to save backup file: %v", err)
	}

	log.Print(backupFile.Name())

	dbHost := m.Provider.Config.Db.DbHost
	dbPort := m.Provider.Config.Db.DbPort
	dbName := m.Provider.Config.Db.DbName
	dbUser := m.Provider.Config.Db.DbUsername
	containerName := m.Provider.Config.Db.ContainerName

	copyCmd := exec.Command("docker", "cp", backupFile.Name(), fmt.Sprintf("%s:/", containerName))
	err = copyCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to copy dump to db: %v", err)
	}

	deleteCmd := exec.Command("docker", "exec", containerName, "psql", "-U", dbUser, "-h", dbHost,
		"-p", dbPort, "-d", dbName, "-c", "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	err = deleteCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete db: %v", err)
	}

	restoreCmd := exec.Command("docker", "exec", containerName,
		"psql",
		"-U", dbUser, "-h", dbHost,
		"-p", dbPort, "-d", dbName, "-f", backupFile.Name())
	err = restoreCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run restore: %v", err)
	}

	return nil
}
