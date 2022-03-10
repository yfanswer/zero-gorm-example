package main

import (
	perrors "github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var conf Config

type Config struct {
	DSN string `yaml:"DSN"`
}

func Configuration() *Config {
	return &conf
}

func LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return perrors.WithStack(err)
	}

	if err = yaml.Unmarshal(bytes, &conf); err != nil {
		return perrors.WithStack(err)
	}
	return nil
}
