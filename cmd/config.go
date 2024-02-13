package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// configCmd returns a *cobra.Command for managing configuration.
// It provides subcommands for viewing and modifying the configuration.
func configCmd() *cobra.Command {

	command := cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	// Add "view" subcommand to "config" command
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View configuration",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig()
			configFile := viper.ConfigFileUsed()
			fmt.Printf("Configuration file is located at: %s\n", configFile)

			config := viper.AllSettings()
			yamlConfig, err := yaml.Marshal(config)

			if err != nil {
				log.Fatalf("Failed to marshal config to YAML: %v", err)
			}

			fmt.Println(string(yamlConfig))
		},
	}

	command.AddCommand(viewCmd)

	return &command
}

// initConfig initializes the configuration for the application.
// It sets the name and type of the config file, adds the config path,
// and reads the config file into the viper configuration object.
// If there is an error reading the config file, it logs a fatal error.
func initConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}
