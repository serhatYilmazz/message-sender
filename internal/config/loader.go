package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Load(logger *logrus.Logger, configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logger.WithError(err).Errorf("error while reading the config")
		return nil, err
	}

	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		logger.Errorf("error while decoding the config")
		return nil, err
	}

	logger.Infof("Config loaded from: %s", v.ConfigFileUsed())
	return &cfg, nil
}
