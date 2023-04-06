package model

type Config struct {
	Db struct {
		DbName        string `yaml:"dbName"`
		DbPort        string `yaml:"dbPort"`
		DbHost        string `yaml:"dbHost"`
		DbUsername    string `yaml:"dbUsername"`
		DbPassword    string `yaml:"dbPassword"`
		ContainerName string `yaml:"containerName"`
	} `yaml:"DB"`
	Minio struct {
		MinioEndpoint  string `yaml:"minioEndpoint"`
		MinioPort      string `yaml:"minioPort"`
		MinioAccessKey string `yaml:"minioAccessKey"`
		MinioSecretKey string `yaml:"minioSecretKey"`
		MinioSSL       string `yaml:"minioSSL"`
		MinioBucket    string `yaml:"minioBucket"`
	} `yaml:"Minio"`
}
