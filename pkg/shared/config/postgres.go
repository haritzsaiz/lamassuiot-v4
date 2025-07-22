package config

type PostgresConfig struct {
	Hostname string   `mapstructure:"hostname"`
	Port     int      `mapstructure:"port"`
	Username string   `mapstructure:"username"`
	Password Password `mapstructure:"password"`
}
