package storage

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDictFromFile(t *testing.T) {
	filename := "quotes_test.txt"

	f, err := os.Create(filename)
	require.NoError(t, err)

	defer os.Remove(filename)

	expectedQuotes := []string{
		"first",
		"second",
	}

	_, err = f.Write([]byte(strings.Join(expectedQuotes, "\n")))
	require.NoError(t, err)

	dict, err := NewDictFromFile(filename)
	require.NoError(t, err)

	require.Equal(t, expectedQuotes, dict.quotes)
}

func TestDict_GetRandomPhrase(t *testing.T) {
	filename := "quotes_test.txt"

	f, err := os.Create(filename)
	require.NoError(t, err)

	defer os.Remove(filename)

	expectedQuotes := []string{
		"first",
		"second",
	}

	_, err = f.Write([]byte(strings.Join(expectedQuotes, "\n")))
	require.NoError(t, err)

	dict, err := NewDictFromFile(filename)
	require.NoError(t, err)

	phrase, err := dict.GetRandomPhrase()
	require.NoError(t, err)
	require.Condition(t, func() bool {
		return expectedQuotes[0] == phrase || expectedQuotes[1] == phrase
	}, "Must be equal first or second qoutes")
}
