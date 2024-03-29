package cmd

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {

	command := cobra.Command{
		Use:   "version",
		Short: "Print version/build info",
		Long:  "Print version/build information",
		Run: func(cmd *cobra.Command, args []string) {
			d := color.New(color.FgHiBlue)
			d.Print(figure.NewFigure("Scrum Chrono", "speed", true).String())
			d.Println()
			d.Println("Version ", Version)
			d.Println("Commit ", Commit)
			d.Println("Date ", Date)
		},
	}

	return &command
}
