/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"os"

	"github.com/TestAndWin/e2e-coverage/logger"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	JWTKey     string `mapstructure:"JWT_KEY"`
}

// Returns the config. When the DB_USER is set as env variable, all values will be read from the environment variables.
// Otherwise the config is read from the config.env file
func LoadConfig() (config Config, err error) {
	u, b := os.LookupEnv("DB_USER")
	if b {
		logger.Debugf("Read config from environment variables")
		var c Config
		c.DBUser = u
		c.DBPassword = os.Getenv("DB_PASSWORD")
		c.DBHost = os.Getenv("DB_HOST")
		c.JWTKey = os.Getenv("JWT_KEY")
		return c, nil
	} else {
		logger.Debugf("Read config from config.env")
		viper.SetConfigFile("config.env")
		viper.AutomaticEnv()
		err = viper.ReadInConfig()
		if err != nil {
			logger.Errorf("Error while reading config file %v", err)
			return
		}
		err = viper.Unmarshal(&config)
		return
	}

}
