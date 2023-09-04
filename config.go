package main

import (
	"flag"
	"fmt"
	"os"

	"io/ioutil"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type Web struct {
	Port string `default:"3000" yaml:"port"`
	User string `default:"webadmin" yaml:"user"`
	Pass string `default:"password" yaml:"pass"`
}

func (w *Web) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(w)

	type plain Web
	if err := unmarshal((*plain)(w)); err != nil {
		return err
	}

	return nil
}

type Db struct {
	Port string `default:"5432" yaml:"port"`
	Host string `default:"127.0.0.1" yaml:"host"`
	Name string `default:"mydb" yaml:"name"`
	User string `default:"user" yaml:"user"`
	Pass string `default:"password" yaml:"pass"`
}

func (d *Db) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(d)

	type plain Db
	if err := unmarshal((*plain)(d)); err != nil {
		return err
	}

	return nil
}

type Config struct {
	Web Web `yaml:"web"`
	Db  Db  `yaml:"db"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

func ParseConfig(configPath string) (*Config, error) {
	config := &Config{}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(content), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func WriteConfig(configPath string, config *Config) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(content)
	f.Sync()

	return nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if os.IsNotExist(err) {
		defaultWeb := &Web{}
		defaults.Set(defaultWeb)
		defaultDb := &Db{}
		defaults.Set(defaultDb)
		defaultConfig := &Config{Web: *defaultWeb, Db: *defaultDb}

		err = WriteConfig(path, defaultConfig)
		if err != nil {
			return err
		}

		return nil
	} else if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}
