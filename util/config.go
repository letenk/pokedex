package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
type Config struct {
	APP_PORT        string `mapstructure:"APP_PORT"`
	JWT_SECRET_KEY  string `mapstructure:"JWT_SECRET_KEY"`
	DB_DRIVER       string `mapstructure:"DB_DRIVER"`
	DB_SOURCE       string `mapstructure:"DB_SOURCE"`
	DB_SOURCE_TEST  string `mapstructure:"DB_SOURCE_TEST"`
	AWS_REGION      string `mapstructure:"AWS_REGION"`
	AWS_KEY_ID      string `mapstructure:"AWS_KEY_ID"`
	AWS_SECRET_KEY  string `mapstructure:"AWS_SECRET_KEY"`
	AWS_BUCKET_NAME string `mapstructure:"AWS_BUCKET_NAME"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// Config viper
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Check all variable in env
	viper.AutomaticEnv()

	// Find and read variable the config file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Insert value config into object viper
	err = viper.Unmarshal(&config)
	return
}
