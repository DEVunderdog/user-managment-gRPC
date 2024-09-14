package utils

import "github.com/spf13/viper"

type Config struct {
	HTTPServer string `mapstructure:"HTTP_SERVER"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	Issuer     string `mapstructure:"ISSUER"`
	Audience   string `mapstructure:"AUDIENCE"`
	Passphrase string `mapstructure:"PASSPHRASE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)
	
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}


