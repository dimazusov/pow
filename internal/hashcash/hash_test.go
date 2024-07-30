package hashcash_test

import (
	"testing"

	"github.com/dimazusov/pow/internal/hashcash"
	"github.com/stretchr/testify/require"
)

func TestHashCash_Verify(t *testing.T) {
	var difficulty uint8 = 3
	solver := hashcash.New(difficulty)
	chlng, err := solver.Challenge()
	require.NoError(t, err)
	require.Len(t, chlng, int(difficulty))

	nonce, err := solver.Calculate(chlng)
	require.NoError(t, err)

	require.True(t, solver.Verify(chlng, nonce))
}
