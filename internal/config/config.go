package config

import "time"

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	SslMode  string `mapstructure:"sslmode"`
}

type WebhookConfig struct {
	URL     string        `mapstructure:"url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type SchedulerConfig struct {
	Interval    time.Duration `mapstructure:"interval"`
	BatchSize   int           `mapstructure:"batch_size"`
	SendTimeout time.Duration `mapstructure:"send_timeout"`
	Enabled     bool          `mapstructure:"enabled"`
}

type Config struct {
	DbConfig        DbConfig        `mapstructure:"database"`
	WebhookConfig   WebhookConfig   `mapstructure:"webhook"`
	SchedulerConfig SchedulerConfig `mapstructure:"scheduler"`
}
