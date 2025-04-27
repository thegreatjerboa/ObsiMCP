package config

import (
	"fmt"
	"os"
	"path/filepath"

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
    
    dir, _ := os.Executable()
    dir = filepath.Dir(dir) // 获取当前程序所在目录

    viper.SetConfigName("config-dev")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(dir)
    viper.AddConfigPath(filepath.Join(dir, "config"))
    viper.AddConfigPath(".") // 再保险一层，当前运行目录

    if err := viper.ReadInConfig(); err != nil {
        panic(fmt.Errorf("fatal error loading config file: %w", err))
    }

    if err := viper.Unmarshal(&Cfg); err != nil {
        panic(fmt.Errorf("unable to decode config into struct: %w", err))
    }

    fmt.Println("Config loaded successfully")
}



