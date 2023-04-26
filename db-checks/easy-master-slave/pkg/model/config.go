package model

type Config struct {
	DbMaster struct {
		DbPort string `yaml:"dbPort"`
		DbHost string `yaml:"dbHost"`
	} `yaml:"DB-master"`
	DbSlave struct {
		DbPort string `yaml:"dbPort"`
		DbHost string `yaml:"dbHost"`
	} `yaml:"DB-slave"`
	DbCommon struct {
		DbName     string `yaml:"dbName"`
		DbUsername string `yaml:"dbUsername"`
		DbPassword string `yaml:"dbPassword"`
	} `yaml:"DB-common"`
}
