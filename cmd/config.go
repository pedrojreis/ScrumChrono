package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"strings"

	"github.com/pedrojreis/ScrumChrono/core"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var ErrInvalidEmail = errors.New("invalid email")

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
			if configFile != "" && viper.IsSet("teams") {
				d.Printf("Configuration file is located at: %s\n", configFile)
			} else {
				d.Println("Can't detect any configuration file at $HOME/.scrumchrono.yaml or the configuration is empty. Please run the wizard to create one.")
				return
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

	wizardCmd := &cobra.Command{
		Use:   "wizard",
		Short: "Run the configuration wizard",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig()
			runWizard()
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete team from configuration",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig()
			teams := getTeams()
			teamName, err := prompt.New().Ask("Team Name:").Choose(teams, choose.WithDefaultIndex(0), choose.WithHelp(true))
			core.CheckErr(err)

			configMap := viper.AllSettings()
			delete(configMap["teams"].(map[string]interface{}), teamName)
			encodedConfig, err := yaml.Marshal(&configMap)

			if err != nil {
				log.Fatalf("Failed to marshal config to YAML: %v", err)
			}

			if err = viper.ReadConfig(bytes.NewReader(encodedConfig)); err != nil {
				log.Fatalf("Failed to delete team from configuration: %v", err)
			}

			err = viper.WriteConfig()

			if err != nil {
				log.Fatalf("Failed to write configuration file: %v", err)
			}
		},
	}

	command.AddCommand(viewCmd, wizardCmd, deleteCmd)

	return &command
}

// initConfig initializes the configuration for the application.
// It sets the name and type of the config file, adds the config path,
// and reads the config file into the viper configuration object.
// If there is an error reading the config file, it logs a fatal error.
func initConfig() {
	viper.SetConfigName(".scrumchrono")          // name of config file (without extension)
	viper.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/")                // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := viper.SafeWriteConfig()
			if err != nil {
				log.Fatalf("Can't detect any configuration file at $HOME/.scrumchrono.yaml. Please run the wizard to create one.")
			}
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}
}

// runWizard is a function that runs a configuration wizard to set up Jira settings and optionally add a team.
func runWizard() {
	// 1. Jira settings
	jiraWizard()

	for i := 0; ; i++ {
		// Ask if we want to add a team
		promptText := "Do you want to add a team?"

		if i >= 1 {
			promptText = "Do you want to add another team?"
		}

		continueWizard, err := prompt.New().Ask(promptText).
			Choose(
				[]string{"Yes", "No"},
				choose.WithDefaultIndex(0),
				choose.WithHelp(true),
			)

		core.CheckErr(err)

		if continueWizard == "No" {
			return
		}

		// 2. Team settings
		teamWizard(i)
	}
}

// jiraWizard is a function that prompts the user to enter Jira configuration details such as the Jira URL, username, and token.
// It uses the prompt package to display input prompts and validate user input.
// The entered configuration details are then stored in the viper configuration and written to the config file.
func jiraWizard() {
	jiraUrl, err := prompt.New().Ask("Jira URL:").Input("https://my-company.atlassian.net/")
	core.CheckErr(err)

	jiraUsername, err := prompt.New().Ask("Jira Username:").Input("username@company.com", input.WithValidateFunc(validateEmail))
	core.CheckErr(err)

	jiraToken, err := prompt.New().Ask("Jira Token:").Input("token", input.WithEchoMode(input.EchoPassword))
	core.CheckErr(err)

	viper.Set("Jira.url", jiraUrl)
	viper.Set("Jira.username", jiraUsername)
	viper.Set("Jira.token", jiraToken)

	err = viper.WriteConfig()

	if err != nil {
		log.Fatalf("Failed to write configuration file: %v", err)
	}
}

func teamWizard(i int) {
	// 2. Team settings
	teamName, err := prompt.New().Ask("Team Name:").Input("TeamName#" + strconv.Itoa(i+1))
	core.CheckErr(err)

	teamBoardID, err := prompt.New().Ask("Board ID:").Input("1")
	core.CheckErr(err)

	teamMaxTime, err := prompt.New().Ask("Max Time:").Input("180")
	core.CheckErr(err)

	teamFont, err := prompt.New().Ask("Font:").Input("doom")
	core.CheckErr(err)

	teamMembers, err := prompt.New().Ask("Members:").Input("Zagreus, Link, Arthur, GLaDOS")
	core.CheckErr(err)

	viper.Set("Teams."+teamName+".BoardID", teamBoardID)
	viper.Set("Teams."+teamName+".MaxTime", teamMaxTime)
	viper.Set("Teams."+teamName+".Font", teamFont)

	teamMembersSplit := strings.Split(teamMembers, ",")
	for i, member := range teamMembersSplit {
		teamMembersSplit[i] = strings.TrimSpace(member)
	}

	viper.Set("Teams."+teamName+".Members", teamMembersSplit)

	err = viper.WriteConfig()

	if err != nil {
		log.Fatalf("Failed to write configuration file: %v", err)
	}
}

// validateEmail validates the given email address.
// It uses the mail.ParseAddress function to parse the email address.
// If the email address is invalid, it returns an error of type ErrInvalidEmail.
// Otherwise, it returns nil.
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%s: %w", email, ErrInvalidEmail)
	} else {
		return nil
	}
}

// getTeams retrieves the list of teams from the configuration file.
// It returns a slice of strings representing the team names.
func getTeams() []string {
	teams := viper.GetStringMap("Teams")
	var teamNames []string
	for k := range teams {
		teamNames = append(teamNames, k)
	}
	return teamNames
}
