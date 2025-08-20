package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Load(logger *logrus.Logger, configPath string) (*Config, error) {
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("error while reading the config")
		return nil, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		logger.Errorf("error while decoding the config")
		return nil, err
	}

	logger.Infof("Config loaded from: %s", viper.ConfigFileUsed())
	return &cfg, nil
}
