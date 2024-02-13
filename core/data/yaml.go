package data

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var filePath = filepath.Join(".", ".data", "file.yaml")

// ReadYamlFile reads the contents of a YAML file and returns the data as a byte slice.
// If the file does not exist, it creates a new file and logs the error.
func ReadYamlFile() ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		file, err := os.Create("file.yaml")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
		log.Println(err)
		return nil, err
	}
	return data, nil
}

// WriteYamlFile writes the given map[string]time.Duration to a YAML file.
// If the file already exists, it will be overwritten.
func WriteYamlFile(data map[string]time.Duration) error {

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(yamlData)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
