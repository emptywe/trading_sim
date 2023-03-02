package config

import "github.com/spf13/viper"

func Init(path, name string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		panic("error initializing configs: " + err.Error())
	}
}
