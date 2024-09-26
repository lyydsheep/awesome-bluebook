package startup

import "github.com/spf13/viper"

func InitViper() {
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
