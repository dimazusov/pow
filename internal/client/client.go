package client

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/dimazusov/pow/api/gen/pb-go/pb/scheme"
	"github.com/dimazusov/pow/internal/config"
	"github.com/dimazusov/pow/internal/connection"
	"github.com/dimazusov/pow/internal/hashcash"
	"google.golang.org/protobuf/proto"
)

const ProtocolVersion = 1

type Client struct {
	cfg *config.ClientConfig
}

func New(cfg *config.ClientConfig) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (m *Client) GetRandomQuote() (string, error) {
	conn, err := net.Dial("tcp", m.cfg.ClientAddress)
	if err != nil {
		return "", fmt.Errorf("failed to dial connection: %w", err)
	}

	c := connection.New(conn)
	c.SetTimeouts(&m.cfg.Connection)

	defer func() {
		if err = conn.Close(); err != nil {
			slog.Error("failed to close connection", "error", err)
		}
	}()

	challenge, difficulty, err := m.getChallenge(c)
	if err != nil {
		return "", fmt.Errorf("failed get challenge: %w", err)
	}

	hc := hashcash.New(difficulty)

	nonce, err := hc.Calculate(challenge)
	if err != nil {
		return "", fmt.Errorf("failed calculate hash: %w", err)
	}

	quote, err := m.getQuote(c, nonce)
	if err != nil {
		return "", fmt.Errorf("failed get quote: %w", err)
	}

	return quote, nil
}

func (m *Client) getChallenge(conn *connection.Connection) ([]byte, uint8, error) {
	req := &scheme.InitRequest{
		ProtocolVersion: ProtocolVersion,
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed marshal msg to pb: %w", err)
	}

	if err = conn.Write(b); err != nil {
		return nil, 0, fmt.Errorf("failed decode protobuf message: %w", err)
	}

	b, err = conn.Read()
	if err != nil {
		return nil, 0, fmt.Errorf("failed read from connection %w", err)
	}

	res := &scheme.InitResponse{}

	err = proto.Unmarshal(b, res)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot unmarshal message: %w", err)
	}

	return res.GetChallenge(), uint8(res.GetDifficulty()), nil
}

func (m *Client) getQuote(conn *connection.Connection, nonce uint64) (string, error) {
	challengeReq := &scheme.ChallengeRequest{
		Nonce: nonce,
	}

	b, err := proto.Marshal(challengeReq)
	if err != nil {
		return "", fmt.Errorf("failed marshal msg to pb: %w", err)
	}

	if err = conn.Write(b); err != nil {
		return "", fmt.Errorf("cannot write message: %w", err)
	}

	b, err = conn.Read()
	if err != nil {
		return "", fmt.Errorf("cannot read connection: %w", err)
	}

	challengeRes := &scheme.ChallengeResponse{}

	err = proto.Unmarshal(b, challengeRes)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal: %w", err)
	}

	return challengeRes.GetQuote(), nil
}
