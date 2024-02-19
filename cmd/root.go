package cmd

import (
	"log"
	"strings"

	"github.com/pedrojreis/ScrumChrono/cmd/countdown"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string = "dev"
	Commit  string
	Date    string

	rootCmd = &cobra.Command{
		Use:   "ScrumChrono",
		Short: "ScrumChrono is a simple CLI tool to manage the time of each member of a team during a Scrum meeting.",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig()
			teams := viper.GetStringMapString("Teams")

			if _, ok := teams[strings.ToLower(team)]; !ok {
				log.Fatalf("Team %s not found in configuration", team)
			}
			countdown.StartCountdown(team)
		},
	}
)

var team string

func init() {

	rootCmd.AddCommand(versionCmd(), configCmd())
	rootCmd.Flags().StringVarP(&team, "team", "t", "", "Team name")
	_ = rootCmd.MarkFlagRequired("team")

}

func Execute() error {
	return rootCmd.Execute()
}
