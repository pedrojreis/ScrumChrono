package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
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
			d := color.New(color.FgHiBlue)
			if configFile != "" {
				d.Printf("Configuration file is located at: %s\n", configFile)
			} else {
				d.Println("I'll generate a new configuration file for you at $HOME/.scrumchrono.yaml.")
			}

			config := viper.AllSettings()
			yamlConfig, err := yaml.Marshal(config)

			if err != nil {
				log.Fatalf("Failed to marshal config to YAML: %v", err)
			}

			fmt.Println()
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
	viper.SetConfigName(".scrumchrono") // name of config file (without extension)
	viper.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/")       // optionally look for config in the working directory
	setDefaults()
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := viper.SafeWriteConfig()
			if err != nil {
				log.Fatalf("Fatal error config file: %s \n", err)
			}
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}
}

// setDefaults sets the default values for the configuration settings.
// It initializes the default values for Jira URL, username, and token,
// as well as the default values for team-specific settings such as
// board ID, maximum time, font, and members.
func setDefaults() {
	// Jira Defaults
	viper.SetDefault("Jira.url", "https://my-company.atlassian.net/")
	viper.SetDefault("Jira.username", "username")
	viper.SetDefault("Jira.token", "token")

	// Team Defaults
	viper.SetDefault("Teams.TeamName.BoardID", "1")
	viper.SetDefault("Teams.TeamName.MaxTime", "180")
	viper.SetDefault("Teams.TeamName.Font", "doom")
	viper.SetDefault("Teams.TeamName.Members", []string{"Zagreus", "Link", "Arthur", "GLaDOS"})
}
