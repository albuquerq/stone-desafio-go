package config

// Config common values for system configuration.
type Config struct {
	DBHost string `json:"db_host"`
	DBPort string `json:"db_port"`
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_password"`

	JwtSecret string `json:"jwt_secret,omitempty"`

	// Port is the listen port for bankingd server.
	Port string `json:"port"`
}
