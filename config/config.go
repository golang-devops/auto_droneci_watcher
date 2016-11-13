package config

import (
	"fmt"
	"io/ioutil"

	"os"

	yaml "gopkg.in/yaml.v2"
)

//LoadConfigFile will load yaml Config from the specified file path
func LoadConfigFile(configFile string) (*Config, error) {
	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to load config file '%s', error: %s", configFile, err.Error())
	}

	c := &Config{}
	if err = yaml.Unmarshal(fileContent, c); err != nil {
		return nil, fmt.Errorf("Unable to parse config file '%s' as yaml, error: %s", configFile, err.Error())
	}

	for _, proj := range c.Projects {
		proj.YamlFile = os.ExpandEnv(proj.YamlFile)
	}

	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

//Config holds the config
type Config struct {
	Projects []*Project
}

func (c *Config) validate() error {
	for _, proj := range c.Projects {
		if _, err := os.Stat(proj.YamlFile); err != nil {
			return fmt.Errorf("Error determining status of drone yaml file '%s', error: %s", proj.YamlFile, err.Error())
		}

		for _, secretLine := range proj.Secrets {
			if _, err := ParseSecretLine(secretLine); err != nil {
				return err
			}
		}
	}

	return nil
}
