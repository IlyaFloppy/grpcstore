package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadFile(path string) (Config, error) {
	b, err := ioutil.ReadFile(path) //nolint:gosec
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to read config")
	}

	var conf Config
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to parse config")
	}

	return conf, nil
}
