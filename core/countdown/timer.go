package countdown

import (
	"fmt"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/gizak/termui/v3"
	"github.com/spf13/viper"
)

func (app *App) InternalTimer(team string, quitCh chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			currentPerson := app.userList.Rows[app.userList.SelectedRow]
			minutes := int(app.userTimers[currentPerson] / time.Minute)
			seconds := int((app.userTimers[currentPerson] % time.Minute) / time.Second)
			maxTime := viper.GetInt("Teams." + team + ".MaxTime")
			switch {
			case app.countdownIsPaused:
				app.countdownText.TextStyle = termui.NewStyle(termui.Color(8))
				app.countdownText.BorderStyle = termui.NewStyle(termui.Color(8))
			case int(app.userTimers[currentPerson]/time.Second) < (maxTime / 3):
				app.countdownText.TextStyle = termui.NewStyle(termui.ColorGreen)
				app.countdownText.BorderStyle = termui.NewStyle(termui.ColorGreen)
			case int(app.userTimers[currentPerson]/time.Second) < (maxTime * 2 / 3):
				app.countdownText.TextStyle = termui.NewStyle(termui.ColorYellow)
				app.countdownText.BorderStyle = termui.NewStyle(termui.ColorYellow)
			case int(app.userTimers[currentPerson]/time.Second) < (maxTime):
				app.countdownText.TextStyle = termui.NewStyle(termui.Color(202))
				app.countdownText.BorderStyle = termui.NewStyle(termui.Color(202))
			default:
				app.countdownText.TextStyle = termui.NewStyle(termui.ColorRed)
				app.countdownText.BorderStyle = termui.NewStyle(termui.ColorRed)
			}

			app.userList.SelectedRowStyle = app.countdownText.TextStyle

			app.countdownText.Text = figure.NewFigure(fmt.Sprintf("%02d:%02d", minutes, seconds), viper.GetString("Teams."+team+".Font"), true).String()
			termui.Render(app.uiGrid)

			if app.countdownIsPaused {
				continue
			}

			app.userTimers[currentPerson] += time.Millisecond * 100
		case <-quitCh:
			ticker.Stop()
			return
		}
	}
}
