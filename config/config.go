package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// 定义配置结构体
type Config struct {
	MainConfig struct {
		AppName  string `yaml:"appName"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		WorkerId int    `yaml:"workerId"`
	} `yaml:"mainConfig"`

	MySQLConfig struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"databaseName"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxIdleTime  int    `yaml:"maxIdleTime"`
	} `yaml:"mysqlConfig"`

	RedisConfig struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redisConfig"`

	LogConfig struct {
		LogPath string `yaml:"logPath"`
	} `yaml:"logConfig"`

	KafkaConfig struct {
		MessageMode string `yaml:"messageMode"`
		HostPort    string `yaml:"hostPort"`
		OrderTopic  string `yaml:"orderTopic"`
		Partition   int    `yaml:"partition"`
		Timeout     int    `yaml:"timeout"`
	} `yaml:"kafkaConfig"`
}

var MyConfig Config

func ConfigInit() {
	// 读取 YAML 文件
	data, err := os.ReadFile("/home/doovvv/code/golang/doovvvDP/config/config.yaml") // 请确保 config.yaml 路径正确
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// 解析 YAML 文件
	err = yaml.Unmarshal(data, &MyConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
}
