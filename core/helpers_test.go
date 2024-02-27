package core

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/cqroot/prompt"
	"github.com/stretchr/testify/assert"
)

func TestRandomizeOrder(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie", "Dave", "Eve"}
	namesCopy := make([]string, len(names))
	copy(namesCopy, names)
	randomizedNames := RandomizeOrder(namesCopy)

	// Assert that the length of the randomized names is the same as the original names
	assert.Len(t, randomizedNames, len(names))

	// Assert that the randomized names contain the same elements as the original names
	for _, name := range names {
		assert.Contains(t, randomizedNames, name)
	}

	// Assert that the order of the names has been randomized
	assert.NotEqual(t, names, randomizedNames)
}

func TestCheckErr(t *testing.T) {
	// Test case 1: Error is nil, should not panic
	CheckErr(nil)

	// Test case 2: Error is not prompt.ErrUserQuit, should panic
	err := errors.New("Some error")
	assert.Panics(t, func() { CheckErr(err) }, "Expected panic")
}

func TestCheckErr_ErrUserQuit(t *testing.T) {
	err := prompt.ErrUserQuit

	if os.Getenv("BE_CRASHER") == "1" {
		CheckErr(err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErr_ErrUserQuit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
