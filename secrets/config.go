package secrets

import (
	"fmt"
	"os"

	"kiberon-labs/vault-env/v2/models"

	"gopkg.in/yaml.v2"
)

// We retrieve the YAML config file, expanding any environment variables we find in it using syntax like ${VAR_NAME}.
func getYamlConfig(inputFile string) (*string, error) {
	rawBytes, ioErr := os.ReadFile(inputFile)
	if ioErr != nil {
		return nil, fmt.Errorf("Error parsing template: %v", ioErr)
	}
	data := os.ExpandEnv(string(rawBytes[:]))
	return &data, nil
}

// GetConfig reads the YAML configuration file and returns a Secrets object with initialized defaults.
func getConfig(filepath string) (*models.Secrets, error) {

	var yamlValues models.SecretsRoot
	// Load the config file
	data, ioErr := getYamlConfig(filepath)
	if ioErr != nil {
		return nil, fmt.Errorf("Error parsing template: %v", ioErr)
	}

	unmarshalErr := yaml.Unmarshal([]byte(*data), &yamlValues)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("Failed to parse standard input: %v", unmarshalErr)
	}

	secrets, err := yamlValues.Secrets.InitDefaults()
	if err != nil {
		return nil, err
	}

	return secrets, nil
}
