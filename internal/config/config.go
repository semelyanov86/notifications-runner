package config

import (
	_ "github.com/octoper/go-ray"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type (
	Config struct {
		Environment string         `yaml:"environment"`
		Cache       CacheConfig    `yaml:"cache"`
		Db          DatabaseConfig `yaml:"db"`
		Slack       SlackConfig    `yaml:"slack"`
		Telegram    TelegramConfig `yaml:"telegram"`
	}
	CacheConfig struct {
		TTL time.Duration `yaml:"ttl"`
	}
	DatabaseConfig struct {
		Dsn          string
		Host         string
		Login        string
		Password     string
		Dbname       string
		MaxOpenConns int    `yaml:"maxOpenConns"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxIdleTime  string `yaml:"maxIdleTime"`
	}
	SlackConfig struct {
		Token     string `yaml:"token"`
		ChannelId string `yaml:"channelId"`
	}
	TelegramConfig struct {
		Token string `yaml:"token"`
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) *Config {
	var cfg *Config
	cfg = readConfigFile(cfg, configsDir)
	cfg.Db.Dsn = cfg.Db.Login + ":" + cfg.Db.Password + "@" + cfg.Db.Host + "/" + cfg.Db.Dbname + "?parseTime=true"
	return cfg
}

func readConfigFile(cfg *Config, configsDir string) *Config {
	bytesOut, err := os.ReadFile(configsDir + "/notification.yaml")

	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bytesOut, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
