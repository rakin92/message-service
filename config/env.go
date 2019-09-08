package config

import "github.com/spf13/viper"

// SetupEnv : sets up project environment
func SetupEnv() {
	viper.AutomaticEnv()

	viper.SetDefault("ALLOWED_ORIGIN", "")

	// twilio
	viper.SetDefault("TWILIO_ACCOUNT_SID", "")
	viper.SetDefault("TWILIO_AUTH_TOKEN", "")

	// sendgrid
	viper.SetDefault("SENDGRID_API_KEY", "")

	// database
	viper.SetDefault("DB_HOST", "")
	viper.SetDefault("DB_PORT", "")
	viper.SetDefault("DB_USER", "")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "")
	viper.SetDefault("DB_SSL_MODE", "")

}
