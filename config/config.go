package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// 定义配置结构体

type Config struct {
	Server     ServerConfig     `yaml:"server" json:"server"`
	Database   DatabaseConfig   `yaml:"database" json:"database"`
	Redis      RedisConfig      `yaml:"redis" json:"redis"`
	JWT        JWTConfig        `yaml:"jwt" json:"jwt"`
	WhiteList  []string         `yaml:"white_list" json:"whiteList"`
	ApiPermits []ApiPermitsItem `yaml:"api_permits" json:"apiPermits"`
}

type ServerConfig struct {
	Port int `yaml:"port" json:"port"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver" json:"driver"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"dbname" json:"DBName"`
}

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret" json:"secret"`
	ExpireTime int    `yaml:"expire-time" json:"expireTime"` // 过期时间（小时）
}

type ApiPermitsItem struct {
	Method  string `yaml:"method" json:"method"`
	Path    string `yaml:"path" json:"path"`
	Permits string `yaml:"permits" json:"permits"`
}

// Load 加载配置文件并返回配置对象
func Load() *Config {
	data, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("yaml配置解析失败: %v", err)
	}

	return &config
}
