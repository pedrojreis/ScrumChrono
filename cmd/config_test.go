package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetTeams(t *testing.T) {
	// Set up test data
	viper.Set("Teams", map[string]interface{}{
		"Team1": "Value1",
		"Team2": "Value2",
		"Team3": "Value3",
	})

	// Call the function under test
	teamNames := getTeams()

	// Assert that the returned team names match the expected values
	expectedTeamNames := []string{"team1", "team2", "team3"}
	assert.Equal(t, expectedTeamNames, teamNames)
}

func TestValidateEmail_ValidEmail(t *testing.T) {
	err := validateEmail("test@example.com")

	assert.NoError(t, err)
}

func TestValidateEmail_InvalidEmail(t *testing.T) {
	email := "invalid_email"

	err := validateEmail(email)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("%s: %v", email, ErrInvalidEmail))
}
