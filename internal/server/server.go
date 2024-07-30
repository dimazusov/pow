package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/dimazusov/pow/internal/config"
	"github.com/dimazusov/pow/internal/connection"
	"github.com/dimazusov/pow/internal/hashcash"
	"github.com/dimazusov/pow/internal/storage"
)

const (
	tcpProtocol = "tcp"
)

type POW interface {
	GetDifficulty() uint8
	Challenge() ([]byte, error)
	Verify(challenge []byte, nonce uint64) bool
}

type Server struct {
	cfg       *config.ServerConfig
	dict      *storage.Dict
	powSolver *hashcash.HashCash
}

func New(cfg *config.ServerConfig) (*Server, error) {
	dictionary, err := storage.NewDictFromFile("../../data/word_of_wisdom.txt")
	if err != nil {
		return nil, fmt.Errorf("failed load from file %w", err)
	}

	return &Server{
		cfg:       cfg,
		dict:      dictionary,
		powSolver: hashcash.New(cfg.PowDifficulty),
	}, nil
}

func (m *Server) Start(ctx context.Context) error {
	ls := net.ListenConfig{
		KeepAlive: m.cfg.Connection.KeepAlive,
	}

	listener, err := ls.Listen(ctx, tcpProtocol, m.cfg.ServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen address: %s, %w", m.cfg.ServerAddress, err)
	}

	defer listener.Close()

	slog.Info("Listening", "address", m.cfg.ServerAddress)

	for {
		if err = ctx.Err(); err != nil {
			return fmt.Errorf("context error: %w", err)
		}

		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("failed  accepting tcp connection: %w", err)
		}

		go m.runConnectionHandler(ctx, conn)
	}
}

func (m *Server) runConnectionHandler(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			slog.Info("[PANIC]", "error", err)
		}
	}()

	c := connection.New(conn)
	defer conn.Close()

	c.SetTimeouts(&m.cfg.Connection)

	handler := NewConnectionHandler(c, m.powSolver, storage.NewValue(), m.dict)
	handler.Handle(ctx)
}
