package config

import (
	"os"

	"gopkg.in/yaml.v3"
	"rnd-git.valsun.cn/ebike-server/go-common/common"
	"rnd-git.valsun.cn/ebike-server/go-common/logs"
)

type AppConfig struct {
	EnableDebug   bool         `yaml:"enableDebug"`
	ServerAddress string       `yaml:"serverAddress"`
	MaxBodySize   int          `yaml:"maxBodySize"`
	Log           logs.ConfLog `yaml:"log"`

	GptApiKey string `yaml:"gptApiKey"`
	GptUrl    string `yaml:"gptUrl"`
	// GptDeploymentName string `yaml:"gptDeploymentName"`

	QdAddr string `yaml:"qdAddr"`
	QdPort string `yaml:"qdPort"`

	MongoDbAddr       string `yaml:"mongoDbAddr"`
	MongoDbUser       string `yaml:"mongoDbUser"`
	MongoDbPwd        string `yaml:"mongoDbPwd"`
	MongoDbDatabase   string `yaml:"mongoDbDatabase"`
	MongoDbCollection string `yaml:"mongoDbCollection"`
}

// InitDefault set default while not configure
func (c *AppConfig) InitDefault() {
	common.EmptyInitDefault(&c.ServerAddress, ":8080")
	common.ZeroInitDefault(&c.MaxBodySize, 5<<20) // set 5M
}

// PrepareEnv prepare to start server
func (c *AppConfig) PrepareEnv() error {
	// initial log output
	c.Log.InitLogger()
	logs.Infof("PrepareEnv finish")
	return nil
}

var appConf = &AppConfig{}

// LoadContent init config
func LoadContent(content []byte) error {
	conf := &AppConfig{}
	err := yaml.Unmarshal(content, conf)
	if err != nil {
		return err
	}

	// set default
	conf.InitDefault()

	SetAppConf(conf)
	return nil
}

// init from file
func LoadFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadContent(content)
}

// GetAppConf get app config
func GetAppConf() *AppConfig {
	return appConf
}

// GConf set app config for test
func SetAppConf(cfg *AppConfig) {
	appConf = cfg
}
