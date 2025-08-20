package config

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SslMode  string
}

type Config struct {
	DbConfig DbConfig
}
