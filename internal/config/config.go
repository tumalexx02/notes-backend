package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string `mapstructure:"env"`
	MigrationsPath string `mapstructure:"migrations_path"`
	IsReload       bool   `mapstructure:"is_reload"`
	Postgres       `mapstructure:"postgres"`
	HTTPServer     `mapstructure:"http_server"`
	Authorization  `mapstructure:"authorization"`
	Image          `mapstructure:"image"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Authorization struct {
	JWTSecret  string        `mapstructure:"secret"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
	Salt       string        `mapstructure:"salt"`
}

type Image struct {
	ImagesDir string `mapstructure:"images_dir"`
	ImageSalt string `mapstructure:"image_salt"`
	MaxWidth  uint   `mapstructure:"max_width"`
}

func MustLoad() *Config {
	var cfgPath string

	// load config path from .env
	_ = godotenv.Load()
	cfgPath = os.Getenv("CONFIG_PATH")

	// validate config path
	if cfgPath == "" {
		log.Fatal("cannot load config")
		os.Exit(1)
	}

	cfgPath = strings.Trim(cfgPath, "\"")

	// load config file
	viper.SetConfigFile(cfgPath)

	// read config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("cannot read config", err)
		os.Exit(1)
	}

	// unmarshal config
	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("cannot unmarshal config" + err.Error())
	}

	return cfg
}
