package minio

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"backup-restore/pkg/model"
)

type MinioClient struct {
	Config *model.Config
	Client *minio.Client
}

func NewMinioClient(config *model.Config) (*MinioClient, error) {
	minioEndpoint := fmt.Sprintf("%s:%s", config.Minio.MinioEndpoint, config.Minio.MinioPort)

	ssl, err := strconv.ParseBool(config.Minio.MinioSSL)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Minio.MinioAccessKey, config.Minio.MinioSecretKey, ""),
		Secure: ssl,
	})

	if err != nil {
		return nil, err
	}

	return &MinioClient{
		Client: minioClient,
		Config: config,
	}, nil
}

func (m *MinioClient) AddFileToBucket(bucketName, filePath string, file *os.File) error {
	_, err := m.Client.PutObject(context.Background(), bucketName, filePath, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioClient) GetFileFromBucket(ctx context.Context, bucketName string) (*minio.Object, error) {
	object, err := m.getTheLatestBackupFile(bucketName)
	if err != nil {
		return nil, err
	}

	downloadOpts := minio.GetObjectOptions{
		VersionID: object.VersionID,
	}
	result, err := m.Client.GetObject(ctx, bucketName, object.Key, downloadOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	return result, nil
}

func (m *MinioClient) getTheLatestBackupFile(bucketName string) (*minio.ObjectInfo, error) {
	opts := minio.ListObjectsOptions{
		Recursive: true,
	}

	objects := make([]minio.ObjectInfo, 0)
	for obj := range m.Client.ListObjects(context.Background(), bucketName, opts) {
		if obj.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", obj.Err)
		}
		objects = append(objects, obj)
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].LastModified.After(objects[j].LastModified)
	})

	dbPrefix := m.Config.Db.DbName
	var latestObject *minio.ObjectInfo
	for _, obj := range objects {
		if strings.HasPrefix(obj.Key, dbPrefix) {
			latestObject = &obj
			break
		}
	}
	if latestObject == nil {
		return nil, fmt.Errorf("file with prefix %s not found in bucket", dbPrefix)
	}
	return latestObject, nil
}
