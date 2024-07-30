package storage

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
)

type Dict struct {
	quotes []string
}

func NewDictFromFile(filepath string) (*Dict, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %w", err)
	}

	return &Dict{
		quotes: append([]string{}, strings.Split(string(b), "\n")...),
	}, nil
}

func (m *Dict) GetRandomPhrase() (string, error) {
	max := big.NewInt(int64(len(m.quotes)))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed reading random int number %w", err)
	}

	return m.quotes[int(n.Int64())], nil
}
