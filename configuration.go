package main

import (
	"fmt"
	"log"
	"path"

	"github.com/spf13/viper"
)

func configure() {

	// configuration type
	viper.SetConfigType("toml")
	viper.SetConfigName("septum")

	// file locations
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	// bind environment
	viper.SetEnvPrefix("septum")
	viper.AutomaticEnv()

	// setup defaults
	viper.SetDefault("database.path", ".")
	viper.SetDefault("database.name", "septum.db")
	viper.SetDefault("server.port", 79)
	viper.SetDefault("server.api.path", "/api")
	viper.SetDefault("server.web.path", "/web")

	// load configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("configure: %s", err)
	}

	config.dbPath = path.Join(viper.GetString("database.path"), viper.GetString("database.name"))

	dumpConfig()
}

func dumpConfig() {
	keys := viper.AllKeys()
	for i := range keys {
		fmt.Printf("%s: %s\n", keys[i], viper.GetString(keys[i]))
	}
}
