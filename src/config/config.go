package config

import (
	"fmt"
	"path"
	"runtime"

	// "os"
	// "path/filepath"

	"github.com/spf13/viper"
)

type VaultConfig struct {
	Path string `mapstructure:"path"`
}

type BackupConfig struct {
	Path string `mapstructure:"path"`
}

type TemplateConfig struct {
	Path string `mapstructure:"path"`
}

type PluginsConfig struct {
	Rest struct {
		BaseUrl   string `mapstructure:"base_url"`
		AuthToken string `mapstructure:"auth_token"`
	} `mapstructure:"rest_api"`
}

type Config struct {
	Vault    VaultConfig    `mapstructure:"vault"`
	Backup   BackupConfig   `mapstructure:"backup"`
	Template TemplateConfig `mapstructure:"template"`
	Plugins  PluginsConfig  `mapstructure:"plugins"`
}

var Cfg Config

func InitConfig() {

	// dir, _ := os.Getwd()
	_, filename, _, _ := runtime.Caller(0)
	confPath := path.Dir(filename)
	viper.AddConfigPath(confPath)
	viper.SetConfigName("config-dev")
	viper.SetConfigType("yaml")
	// viper.AddConfigPath("./src/config")
	// viper.AddConfigPath(filepath.Join(dir, "src/config"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error loading config file: %w", err))
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(fmt.Errorf("Unable to decode config into struct: %w", err))
	}

	fmt.Println("Config loaded successfully")
}
