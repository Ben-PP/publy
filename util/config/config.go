package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

var ErrConfigLoadFailed = errors.New("failed to load config")

type Config struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	ScriptDir string `mapstructure:"script-dir"`
	Pubs      map[string]struct {
		Script    string `mapstructure:"script"`
		TokenHash string `mapstructure:"token-hash"`
	} `mapstructure:"pubs"`
}

var globalConfig *Config

func loadConfig() (config *Config, err error) {
	if os.Getenv("GO_ENV") == "dev" {
		viper.SetConfigName("dev-config")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/publy")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	var localConfig Config

	err = viper.Unmarshal(&localConfig)
	config = &localConfig
	return
}

func Get() (config *Config, err error) {
	if globalConfig == nil {
		globalConfig, err = loadConfig()
		if err != nil {
			err = errors.Join(ErrConfigLoadFailed, err)
			return
		}
	}
	config = globalConfig
	return
}

// Returns the value of GO_ENV environment variable.
func GetGoEnv() string {
	return os.Getenv("GO_ENV")
}
