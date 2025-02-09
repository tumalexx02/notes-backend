package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string `mapstructure:"env"`
	MigrationsPath string `mapstructure:"migrations_path"`
	Postgres       `mapstructure:"postgres"`
	HTTPServer     `mapstructure:"http_server"`
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

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot load .env")
	}

	viper.AutomaticEnv()

	cfgPath := viper.GetString("CONFIG_DIR")
	cfgName := viper.GetString("CONFIG_NAME")

	viper.SetConfigName(cfgName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("cannot read config")
	}

	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("cannot unmarshal config")
	}

	return cfg
}
