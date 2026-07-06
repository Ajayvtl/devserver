package config

import (
	"context"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service  ServiceConfig  `yaml:"service"`
	Security SecurityConfig `yaml:"security"`
	Logging  LoggingConfig  `yaml:"logging"`
	State    StateConfig    `yaml:"state"`
	Database DatabaseConfig `yaml:"database"`
	Modules  []string       `yaml:"modules"`
	Commands CommandsConfig `yaml:"commands"`
}

type ServiceConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type StateConfig struct {
	Path string `yaml:"path"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type SecurityConfig struct {
	MasterKey string `yaml:"masterKey"`
}

type CommandsConfig struct {
	Default string `yaml:"default"`
}

func Default() Config {
	return Config{
		Service: ServiceConfig{
			Name: "devserver",
			Host: "localhost",
		},
		Security: SecurityConfig{
			MasterKey: "00000000000000000000000000000000", // 32 byte default key
		},
		Logging: LoggingConfig{
			Level: "info",
		},
		State: StateConfig{
			Path: "configs/state.yaml",
		},
		Database: DatabaseConfig{
			DSN: "root:12345678@tcp(127.0.0.1:3306)/devops",
		},
		Modules: []string{},
		Commands: CommandsConfig{
			Default: "doctor",
		},
	}
}

type Loader struct {
	defaults Config
}

func NewLoader() *Loader {
	return &Loader{defaults: Default()}
}

func (l *Loader) Load(ctx context.Context, path string) (Config, error) {
	if err := ctx.Err(); err != nil {
		return Config{}, err
	}

	cfg := l.defaults
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return Config{}, err
			}
		} else if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, err
		}
	}

	applyEnv(&cfg)
	cfg.normalize()

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func applyEnv(cfg *Config) {
	if value := os.Getenv("DEVSERVER_SERVICE_NAME"); value != "" {
		cfg.Service.Name = value
	}
	if value := os.Getenv("DEVSERVER_SERVICE_HOST"); value != "" {
		cfg.Service.Host = value
	}
	if value := os.Getenv("DEVSERVER_MASTER_KEY"); value != "" {
		cfg.Security.MasterKey = value
	}
	if value := os.Getenv("DEVSERVER_LOG_LEVEL"); value != "" {
		cfg.Logging.Level = value
	}
	if value := os.Getenv("DEVSERVER_STATE_PATH"); value != "" {
		cfg.State.Path = value
	}
	if value := os.Getenv("DEVSERVER_DATABASE_DSN"); value != "" {
		cfg.Database.DSN = value
	}
	if value := os.Getenv("DEVSERVER_DEFAULT_COMMAND"); value != "" {
		cfg.Commands.Default = value
	}
	if value := os.Getenv("DEVSERVER_MODULES"); value != "" {
		cfg.Modules = splitCSV(value)
	}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func (c *Config) normalize() {
	if c.Service.Name == "" {
		c.Service.Name = "devserver"
	}
	if c.Service.Host == "" {
		c.Service.Host = "localhost"
	}
	if c.Security.MasterKey == "" {
		c.Security.MasterKey = "00000000000000000000000000000000"
	}
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.State.Path == "" {
		c.State.Path = "configs/state.yaml"
	}
	if c.Database.DSN == "" {
		c.Database.DSN = "root:12345678@tcp(127.0.0.1:3306)/devops"
	}
	if c.Commands.Default == "" {
		c.Commands.Default = "doctor"
	}
}

func (c Config) Validate() error {
	if c.Service.Name == "" {
		return errors.New("service.name is required")
	}
	if c.Logging.Level == "" {
		return errors.New("logging.level is required")
	}
	if c.State.Path == "" {
		return errors.New("state.path is required")
	}
	return nil
}
