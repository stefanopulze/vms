package config

import (
	"flag"
	"os"
	"vms-core/internal/infrastructure/exporter/clickhouse"
	"vms-core/internal/infrastructure/exporter/influx"
	"vms-core/internal/notifier"

	"github.com/stefanopulze/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging LogConfig    `yaml:"logging"`
	Server  ServerConfig `yaml:"server"`
	Serial  struct {
		PortName  string `yaml:"portName" env:"SERIAL_PORT_NAME"`
		BaudRate  int    `yaml:"baudRate" env:"SERIAL_BAUD_RATE"`
		QueueSize int    `yaml:"queueSize" env:"SERIAL_QUEUE_SIZE"`
	} `yaml:"serial"`
	Influx     influx.Options          `yaml:"influx" env:"INFLUX"`
	ClickHouse clickhouse.Options      `yaml:"clickhouse" env:"CLICKHOUSE"`
	Telegram   notifier.TelegramConfig `yaml:"telegram" env:"TELEGRAM"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		flagConfigPath := flag.String("config", "./.env", "config file path")
		flag.Parse()
		configPath = *flagConfigPath
	}

	return LoadConfigFrom(configPath)
}

func loadFromYaml(path string, cfg *Config) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, cfg)
}

func LoadConfigFrom(path string) (*Config, error) {
	cfg := new(Config)

	if err := loadFromYaml("./configs/config.yaml", cfg); err != nil {
		return nil, err
	}

	if err := envconfig.ReadDotEnv(path); err != nil {
		return nil, err
	}

	if err := envconfig.ReadEnv(cfg); err != nil {
		return nil, err
	}

	configureLogging(cfg.Logging)

	return cfg, nil
}
