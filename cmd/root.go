package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pedrojreis/ScrumChrono/core"
	"github.com/pedrojreis/ScrumChrono/core/jira"

	"github.com/common-nighthawk/go-figure"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
			startCountdown()
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

// TODO This method must be moved to new files
func startCountdown() {
	if err := termui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v\n", err)
		return
	}

	// data.InitDB(viper.GetStringSlice("Teams." + team + ".Members"))

	stringTitle := make(chan string)
	go jira.GetIssuesForSprint(team, stringTitle)

	// Init our List with names
	names := core.RandomizeOrder(viper.GetStringSlice("Teams." + team + ".Members"))
	list := widgets.NewList()
	list.Title = "[ " + strings.ToUpper(team) + " ]"
	list.Rows = names
	list.TextStyle = termui.NewStyle(termui.ColorWhite)

	// Init our Countdown
	countdown := widgets.NewParagraph()
	countdown.Text = figure.NewFigure("00:00", viper.GetString("Teams."+team+".Font"), true).String()
	countdown.TextStyle = termui.NewStyle(termui.ColorBlue)
	countdown.Title = "[ Countdown ]"

	// Init our sprintInfo
	sprintInfo := widgets.NewList()
	sprintInfo.Title = "[ " + <-stringTitle + " ]"
	sprintInfo.Rows = jira.GetIssuesForSprintByUser(list.Rows[list.SelectedRow])
	sprintInfo.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	sprintInfo.WrapText = false

	//focusable stuff
	focus := []*widgets.List{list, sprintInfo}
	currentFocus := 0

	grid := termui.NewGrid()

	updateGrid := func() {
		grid = termui.NewGrid()
		termWidth, termHeight := termui.TerminalDimensions()
		grid.SetRect(0, 0, termWidth, termHeight)

		if termWidth > 100 {
			//Countdown UI
			countdown.PaddingLeft = int(termWidth*2.0/3.0/2.0) - (len(strings.Split(countdown.Text, "\n")[0]) / 2.0) //width
			countdown.PaddingTop = int(termHeight*2.0/3.0/2.0 - (len(strings.Split(countdown.Text, "\n")) / 2.0))    //height

			//List UI
			list.PaddingLeft = int(termWidth*1.0/3.0/2.0) - (len(strings.Split(names[0], "\n")[0]) / 2.0) //width

			//Set Grid
			grid.Set(
				termui.NewCol(2.0/3.0, termui.NewRow(2.0/3.0, countdown), termui.NewRow(1.0/3.0, sprintInfo)),
				termui.NewCol(1.0/3.0, termui.NewRow(1.0, list)),
			)
		} else {
			//Countdown UI
			countdown.PaddingLeft = int(termWidth/2) - (len(strings.Split(countdown.Text, "\n")[0]) / 2.0) //width
			countdown.PaddingTop = int(termHeight/3/2) - 4                                                 //why 4? because the height of our text is 6, so we divide that by 2 and add 1 because of the title

			//List UI
			list.PaddingLeft = int(termWidth/2) - (len(names[0]) / 2.0) //width

			//Set Grid
			grid.Set(
				termui.NewCol(1.0,
					termui.NewRow(1.0/3.0, countdown),
					termui.NewRow(1.0/3.0, list),
					termui.NewRow(1.0/3.0, sprintInfo),
				))
		}
	}

	updateGrid()

	termui.Render(grid)

	events := termui.PollEvents()
	quit := make(chan struct{})

	timers := make(map[string]time.Duration)
	isPaused := true

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				currentPerson := list.Rows[list.SelectedRow]
				minutes := int(timers[currentPerson] / time.Minute)
				seconds := int((timers[currentPerson] % time.Minute) / time.Second)
				maxTime := viper.GetInt("Teams." + team + ".MaxTime")
				switch {
				case isPaused:
					countdown.TextStyle = termui.NewStyle(termui.Color(8))
					countdown.BorderStyle = termui.NewStyle(termui.Color(8))
				case int(timers[currentPerson]/time.Second) < (maxTime / 3):
					countdown.TextStyle = termui.NewStyle(termui.ColorGreen)
					countdown.BorderStyle = termui.NewStyle(termui.ColorGreen)
				case int(timers[currentPerson]/time.Second) < (maxTime * 2 / 3):
					countdown.TextStyle = termui.NewStyle(termui.ColorYellow)
					countdown.BorderStyle = termui.NewStyle(termui.ColorYellow)
				case int(timers[currentPerson]/time.Second) < (maxTime):
					countdown.TextStyle = termui.NewStyle(termui.Color(202))
					countdown.BorderStyle = termui.NewStyle(termui.Color(202))
				default:
					countdown.TextStyle = termui.NewStyle(termui.ColorRed)
					countdown.BorderStyle = termui.NewStyle(termui.ColorRed)
				}

				list.SelectedRowStyle = countdown.TextStyle

				countdown.Text = figure.NewFigure(fmt.Sprintf("%02d:%02d", minutes, seconds), viper.GetString("Teams."+team+".Font"), true).String()
				termui.Render(grid)

				if isPaused {
					continue
				}

				timers[currentPerson] += time.Millisecond * 100
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		for {
			e := <-events
			switch e.ID {
			case "q", "<C-c>":
				termui.Close()
				quit <- struct{}{}
				return
			case "<Left>":
				currentFocus = (currentFocus + 1) % len(focus)
			case "<Right>":
				currentFocus = (currentFocus - 1 + len(focus)) % len(focus)
			case "<Up>":
				if len(focus[currentFocus].Rows) > 0 && focus[currentFocus].SelectedRow == 0 {
					focus[currentFocus].SelectedRow = len(focus[currentFocus].Rows) - 1
				} else {
					focus[currentFocus].ScrollUp()
				}

				sprintInfo.Rows = jira.GetIssuesForSprintByUser(list.Rows[list.SelectedRow])
			case "<Down>":
				if len(focus[currentFocus].Rows) > 0 {
					if focus[currentFocus].SelectedRow == len(focus[currentFocus].Rows)-1 {
						focus[currentFocus].SelectedRow = 0
					} else {
						focus[currentFocus].ScrollDown()
					}
				}

				sprintInfo.Rows = jira.GetIssuesForSprintByUser(list.Rows[list.SelectedRow])
			case "<Space>":
				isPaused = !isPaused
			case "<Resize>":
				updateGrid()
				termui.Render(grid)
			}
			for i, f := range focus {
				if i == currentFocus {
					f.TitleStyle = termui.NewStyle(termui.Color(termui.ModifierBold))
				} else {
					f.TitleStyle = termui.NewStyle(termui.ColorWhite)
				}
			}
		}
	}()

	<-quit
}
