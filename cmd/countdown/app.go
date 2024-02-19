package countdown

import (
	"fmt"
	"strings"
	"time"

	"github.com/pedrojreis/ScrumChrono/core"
	"github.com/pedrojreis/ScrumChrono/core/jira"

	"github.com/common-nighthawk/go-figure"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/viper"
)

type App struct {
	users             []string
	userList          *widgets.List
	userTimers        map[string]time.Duration
	countdownText     *widgets.Paragraph
	countdownIsPaused bool
	ticketsList       *widgets.List
	uiGrid            *termui.Grid
}

func (app *App) UpdateGrid() {
	app.uiGrid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	app.uiGrid.SetRect(0, 0, termWidth, termHeight)

	if termWidth > 100 {
		//Countdown UI
		app.countdownText.PaddingLeft = int(termWidth*2.0/3.0/2.0) - (len(strings.Split(app.countdownText.Text, "\n")[0]) / 2.0) //width
		app.countdownText.PaddingTop = int(termHeight*2.0/3.0/2.0 - (len(strings.Split(app.countdownText.Text, "\n")) / 2.0))    //height

		//List UI
		app.userList.PaddingLeft = int(termWidth*1.0/3.0/2.0) - (len(strings.Split(app.users[0], "\n")[0]) / 2.0) //width

		//Set Grid
		app.uiGrid.Set(
			termui.NewCol(2.0/3.0, termui.NewRow(2.0/3.0, app.countdownText), termui.NewRow(1.0/3.0, app.ticketsList)),
			termui.NewCol(1.0/3.0, termui.NewRow(1.0, app.userList)),
		)
	} else {
		//Countdown UI
		app.countdownText.PaddingLeft = int(termWidth/2) - (len(strings.Split(app.countdownText.Text, "\n")[0]) / 2.0) //width
		app.countdownText.PaddingTop = int(termHeight/3/2) - 4                                                         //why 4? because the height of our text is 6, so we divide that by 2 and add 1 because of the title

		//List UI
		app.userList.PaddingLeft = int(termWidth/2) - (len(app.users[0]) / 2.0) //width

		//Set Grid
		app.uiGrid.Set(
			termui.NewCol(1.0,
				termui.NewRow(1.0/3.0, app.countdownText),
				termui.NewRow(1.0/3.0, app.userList),
				termui.NewRow(1.0/3.0, app.ticketsList),
			))
	}
}

// StartCountdown starts the countdown app, by creating the app UI
// and both the timer and keyboard events threads
func StartCountdown(team string) {
	if err := termui.Init(); err != nil {
		fmt.Printf("Failed to initialize termui: %v\n", err)
		return
	}

	// data.InitDB(viper.GetStringSlice("Teams." + team + ".Members"))

	stringTitle := make(chan string)
	go jira.GetIssuesForSprint(team, stringTitle)

	names := core.RandomizeOrder(viper.GetStringSlice("Teams." + team + ".Members"))
	app := App{
		users:         names,
		userList:      widgets.NewList(),
		countdownText: widgets.NewParagraph(),
		ticketsList:   widgets.NewList(),
		uiGrid:        termui.NewGrid(),
	}
	app.userTimers = make(map[string]time.Duration)
	app.countdownIsPaused = true

	// Init our List with names
	app.userList.Title = "[ " + strings.ToUpper(team) + " ]"
	app.userList.Rows = names
	app.userList.TextStyle = termui.NewStyle(termui.ColorWhite)

	// Init our Countdown
	app.countdownText.Text = figure.NewFigure("00:00", viper.GetString("Teams."+team+".Font"), true).String()
	app.countdownText.TextStyle = termui.NewStyle(termui.ColorBlue)
	app.countdownText.Title = "[ Countdown ]"

	// Init our sprintInfo
	app.ticketsList.Title = "[ " + <-stringTitle + " ]"
	app.ticketsList.Rows = jira.GetIssuesForSprintByUser(app.userList.Rows[app.userList.SelectedRow])
	app.ticketsList.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	app.ticketsList.WrapText = false

	app.UpdateGrid()
	termui.Render(app.uiGrid)

	events := termui.PollEvents()
	quit := make(chan struct{})

	// start internal Timer
	go app.InternalTimer(team, quit)

	// start keyboard handler
	go app.KeyboardHandler(events, quit)

	<-quit
}
