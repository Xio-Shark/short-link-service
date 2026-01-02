package config

import "github.com/spf13/viper"

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	HTTP  HTTPConfig  `mapstructure:"http"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type HTTPConfig struct {
	Addr string `mapstructure:"addr"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
