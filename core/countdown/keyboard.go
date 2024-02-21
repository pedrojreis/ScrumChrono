package countdown

import (
	"github.com/pedrojreis/ScrumChrono/core/jira"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func (app *App) KeyboardHandler(events <-chan termui.Event, quitCh chan struct{}) {
	//focusable stuff
	focus := []*widgets.List{app.userList, app.ticketsList}
	currentFocus := 0
	for {
		e := <-events
		switch e.ID {
		case "q", "<C-c>":
			termui.Close()
			quitCh <- struct{}{}
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

			app.ticketsList.Rows = jira.GetIssuesForSprintByUser(app.userList.Rows[app.userList.SelectedRow])
		case "<Down>":
			if len(focus[currentFocus].Rows) > 0 {
				if focus[currentFocus].SelectedRow == len(focus[currentFocus].Rows)-1 {
					focus[currentFocus].SelectedRow = 0
				} else {
					focus[currentFocus].ScrollDown()
				}
			}

			app.ticketsList.Rows = jira.GetIssuesForSprintByUser(app.userList.Rows[app.userList.SelectedRow])
		case "<Space>":
			app.countdownIsPaused = !app.countdownIsPaused
		case "<Resize>":
			//updateGrid()
			app.UpdateGrid()
			termui.Render(app.uiGrid)
		}
		for i, f := range focus {
			if i == currentFocus {
				f.TitleStyle = termui.NewStyle(termui.Color(termui.ModifierBold))
			} else {
				f.TitleStyle = termui.NewStyle(termui.ColorWhite)
			}
		}
	}
}
