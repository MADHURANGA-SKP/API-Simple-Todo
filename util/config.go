package uill

import (
	"time"

	"github.com/spf13/viper"
)

//config stores all configuration of the application
//the values are read by viper from config file or enviroment var's
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
    DBSource      string `mapstructure:"DB_SOURCE"`
    ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
    AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`

}

//load reads configuratin from fileor env var's
func LoadConfig(path string)(config Config, err error){
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err=viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal((&config))
	return
}