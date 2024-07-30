package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	expectedVal := []byte{1, 2}

	strg := NewValue()
	strg.Set(expectedVal)

	require.Equal(t, expectedVal, strg.Get())
}
