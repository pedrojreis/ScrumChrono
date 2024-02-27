package core

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cqroot/prompt"
)

// randomizeOrder randomizes the order of elements in the given slice of names.
// It uses the Fisher-Yates algorithm to shuffle the elements.
// The function takes a slice of strings as input and returns the shuffled slice.
func RandomizeOrder(names []string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(names), func(i, j int) {
		names[i], names[j] = names[j], names[i]
	})
	return names
}

func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}
