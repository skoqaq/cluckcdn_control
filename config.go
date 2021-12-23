package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	Control string `yaml:"control"`
	Token   string `yaml:"token"`
}

var adminConfig map[string]interface{}
var vhostListConfig []map[string]interface{}

var nodeConfig yamlConfig

func init() {
	// NodeConfig
	yamlFile, err := ioutil.ReadFile("static/config.yaml")
	if err != nil {
		fmt.Printf("yamlFile.Get err #%v ", err)
	}
	yamlFileR := strings.ReplaceAll(string(yamlFile), "{ctrlServer}", "")
	err = yaml.Unmarshal([]byte(yamlFileR), &nodeConfig)
	if err != nil {
		fmt.Println("Unmarshal", err)
	}
	upAdminConfig()
	upVhostConfig()
}

func upAdminConfig() {
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("jsonFile.Get err #%v ", err)
	}
	err = json.Unmarshal(jsonFile, &adminConfig)
	if err != nil {
		fmt.Println("Unmarshal", err)
	}
}

func upVhostConfig() {
	jsonFile, err := ioutil.ReadFile("vhost.json")
	if err != nil {
		fmt.Printf("jsonFile.Get err #%v ", err)
	}
	err = json.Unmarshal(jsonFile, &vhostListConfig)
	if err != nil {
		fmt.Println("Unmarshal", err)
	}
}
