package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

const configFile = "config.yaml"

type Metric struct {
	FuncName  string `yaml:"funcName"`
	FieldName string `yaml:"fieldName"`
	TableName string `yaml:"tableName"`
	AggRule   string `yaml:"aggRule"`
}

func readYaml(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("read yaml file failed: %v", err)
	}

	return string(bytes), nil
}

func main() {
	fn := filepath.Join("config/yaml", configFile)
	content, err := readYaml(fn)
	if err != nil {
		log.Fatal(err)
	}

	metrics := make(map[string]Metric)
	if err := yaml.Unmarshal([]byte(content), &metrics); err != nil {
		log.Fatalf("unmarshal yaml failed: %v", err)
	}

	bytes, _ := json.Marshal(metrics)
	log.Println(string(bytes))
}
