package config

import (
	"fmt"

	"os"

	"gopkg.in/yaml.v2"
)

const (
	Ptt      string = "ptt"
	Facebook string = "facebook"
	Twitter  string = "twitter"
	Ticker   string = "ticker"
)

var validSrcTypes = []string{Ptt, Facebook, Twitter, Ticker}

type ParsersConfig struct {
	SrcType       string `yaml:"src_type"`
	BoardName     string `yaml:"board_name"`
	Url           string `yaml:"url,omitempty"`
	TickerName    string `yaml:"ticker_name,omitempty"`
	Interval      string `yaml:"interval"`
	PushCntThresh int    `yaml:"push_cnt_thresh"`
	AliasName     string `yaml:"alias_name,omitempty"`
}

type ElasticSearchConfig struct {
	Port     int    `yaml:"port"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Parsers       []ParsersConfig     `yaml:parsers`
	ElasticSearch ElasticSearchConfig `yaml:elasticsearch`
}

func LoadConfig(configFile string) (Config, error) {
	var config Config
	// viper.SetConfigFile(configFile)
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	return Config{}, fmt.Errorf("Failed to read the configuration file: %v", err)
	// }
	// err = viper.Unmarshal(&config)
	// if err != nil {
	// 	return Config{}, fmt.Errorf("Failed to unmarshal the configuration file: %v", err)
	// }

	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	// check if src_type in parsers is valid
	for _, parser := range config.Parsers {
		if !validSrcType(parser.SrcType) {
			return Config{}, fmt.Errorf("Invalid src_type %v in parsers configuration", parser.SrcType)
		}
	}

	if esEnvVal, ok := os.LookupEnv("ELASTICSEARCH_HOST"); ok {
		config.ElasticSearch.URL = esEnvVal
	}

	return config, nil
}

// check if src_type is valid
func validSrcType(srcType string) bool {
	for _, t := range validSrcTypes {
		if t == srcType {
			return true
		}
	}
	return false
}
