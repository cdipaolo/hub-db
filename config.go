package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const (
	configPath string = "./config.json"
)

var (
	// Config stores the configuration
	// variable
	Config *Configuration

	// config holds the configuration
	// path
	config string
)

// Configuration holds config information
// necessary to run the PornHub crawler
//
// TimeDelaySeconds is the number of seconds
// to wait between crawling each search page
// (between each ~350 image batches)
//
// DumpPath is the path to dump each
// search page json file on the host
// machine
//
// StartPage and EndPage are where the
// crawler starts and stops
type Configuration struct {
	Port       int16 `json:"port,omitempty"`
	portString string

	StartPage uint64 `json:"start"`
	EndPage   uint64 `json:"end"`

	DumpPath         string `json:"dump_path"`
	TimeDelaySeconds uint64 `json:"time_delay,omitempty"`
}

// init takes in the config path from
// the command line flags
func init() {
	flag.StringVar(&config, "C", configPath, "Sets the server configuration filepath")
}

// ParseConfigFromFile takes the passed
// config file path and parses the config
// params from the file located there
func ParseConfigFromFile() error {
	path, err := filepath.Abs(config)
	if err != nil {
		return fmt.Errorf("ERROR: error generating absolute filepath from the given config path. Does this file exist? %v", err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ERROR: error reading the bytes from the configuration filename. %v", err)
	}

	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		return fmt.Errorf("ERROR: error unmarshalling config file from JSON to struct. Is it in json format? %v", err)
	}

	if Config.Port == 0 {
		Config.Port = 8080
	}
	Config.portString = fmt.Sprintf(":%v", Config.Port)

	return nil
}
