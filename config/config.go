package config

import (
	"fmt"

	"github.com/spf13/viper"
)


type VaultConfig struct{
    Path string `mapstructure:"path"`
}

type Config struct{
    Vault VaultConfig `mapstructure:"vault"`
}

var Cfg Config

func InitConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")

    if err := viper.ReadInConfig(); err != nil {
        panic(fmt.Errorf("fatal error loading config file: %w", err))
    }

    if err := viper.Unmarshal(&Cfg); err != nil {
        panic(fmt.Errorf("unable to decode config into struct: %w", err))
    }

    fmt.Println("Config loaded successfully")
}



