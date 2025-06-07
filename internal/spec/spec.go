package spec

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Specification struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Name     string    `json:"name"`
	Origin   string    `json:"origin"`
	Backends []Backend `json:"backends"`
}

type Backend struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Health string `json:"health"`
}

func LoadSpec(path string) (Specification, error) {
	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Specification{}, errors.New("config file not found in path")
		}
		return Specification{}, err
	}

	var s Specification
	if err := v.Unmarshal(&s); err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Specification{}, err
	}

	return s, nil
}
