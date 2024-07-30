package hashcash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

var ErrFailedToFindNonce = errors.New("failed to find nonce")

type HashCash struct {
	reader     io.Reader
	difficulty uint8
}

func New(difficulty uint8) *HashCash {
	return &HashCash{
		reader:     rand.Reader,
		difficulty: difficulty,
	}
}

func (m *HashCash) Verify(challenge []byte, nonce uint64) bool {
	target := make([]byte, m.difficulty)

	return checkChallenge(challenge, nonce, target)
}

func (m *HashCash) Calculate(challenge []byte) (uint64, error) {
	target := make([]byte, m.difficulty)
	for i := uint64(0); i < math.MaxUint64; i++ {
		if checkChallenge(challenge, i, target) {
			return i, nil
		}
	}

	return 0, ErrFailedToFindNonce
}

func (m *HashCash) Challenge() ([]byte, error) {
	c := make([]byte, m.difficulty)
	if _, err := io.ReadFull(m.reader, c[:]); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}

	return c[:], nil
}

func (m *HashCash) GetDifficulty() uint8 {
	return m.difficulty
}

func checkChallenge(challenge []byte, nonce uint64, target []byte) bool {
	nonceBytes := make([]byte, 8) //nolint:mnd,nolintlint
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	challenge = append(challenge, nonceBytes...)
	h := sha256.New()
	h.Write(challenge)
	hash := h.Sum(nil)

	return bytes.Equal(hash[:len(target)], target)
}
