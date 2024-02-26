package jira

import (
	"fmt"
	"unicode"

	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var sprintIssues []jira.Issue

// We have to promise this
func GetIssuesForSprint(team string, out chan<- string) {
	tp := jira.BasicAuthTransport{
		Username: viper.GetString("Jira.username"),
		Password: viper.GetString("Jira.token"),
	}

	url := viper.GetString("Jira.url")

	jiraClient, err := jira.NewClient(tp.Client(), url)

	if err != nil {
		out <- "No active sprint detected"
		return
	}

	sprintlist, _, err := jiraClient.Board.GetAllSprintsWithOptions(viper.GetInt("Teams."+team+".BoardID"), &jira.GetAllSprintsOptions{State: "active"})

	if sprintlist.Values == nil || err != nil {
		out <- "No active sprint detected"
		return
	}

	issues, _, _ := jiraClient.Sprint.GetIssuesForSprint((sprintlist.Values[0].ID))
	sprintIssues = issues
	out <- sprintlist.Values[0].Name
}

func GetIssuesForSprintByUser(name string) []string {
	sumo := []string{}
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	alteredname, _, _ := transform.String(t, name)

	for _, issue := range sprintIssues {
		if issue.Fields.Assignee != nil && issue.Fields.Assignee.DisplayName == alteredname {
			sumo = append(sumo, PaintAccordingToStatus(issue.Fields.Summary, issue.Fields.Status.Name))
		}
	}

	return sumo
}

func PaintAccordingToStatus(text string, status string) string {
	//to add some color just add the following to the string
	// fg:color	- foreground color
	// bg:color	- background color
	switch status {
	case "To Do", "New":
		return fmt.Sprintf("😴 %s", text)
	case "In Progress", "In Development":
		return fmt.Sprintf("🚧 %s", text)
	case "In Tests":
		return fmt.Sprintf("🧪 %s", text)
	case "Done", "In Deployment":
		return fmt.Sprintf("✅ %s", text)
	case "Rejected":
		return fmt.Sprintf("❌ %s", text)
	}

	return text
}
