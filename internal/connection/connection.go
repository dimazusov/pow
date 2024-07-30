package connection

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/dimazusov/pow/internal/config"
)

var ErrWrongExpectedMsgLen = errors.New("wrong expected message length")

type Connection struct {
	net.Conn
}

func New(conn net.Conn) *Connection {
	return &Connection{conn}
}

func (m *Connection) SetTimeouts(cfg *config.ConnectionConfig) {
	if cfg.IdleTimeout > 0 {
		if err := m.SetDeadline(time.Now().Add(cfg.IdleTimeout)); err != nil {
			slog.Info("failed to set idle timeout", "error", err)
			return
		}
	}

	if cfg.ReadTimeout > 0 {
		if err := m.SetReadDeadline(time.Now().Add(cfg.ReadTimeout)); err != nil {
			slog.Info("failed to set reading timeout", "error", err)
			return
		}
	}

	if cfg.WriteTimeout > 0 {
		if err := m.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout)); err != nil {
			slog.Info("failed to set reading timeout", "error", err)
			return
		}
	}
}

func (m *Connection) Read() ([]byte, error) {
	b := make([]byte, 4) //nolint:mnd,nolintlint

	_, err := m.Conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed read bytes: %w", err)
	}

	length := binary.BigEndian.Uint32(b)
	b = make([]byte, length)

	n, err := m.Conn.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed read bytes: %w", err)
	}

	if uint32(n) != length {
		return nil, fmt.Errorf("%w, expected: %d, given: %d", ErrWrongExpectedMsgLen, length, n)
	}

	return b, nil
}

func (m *Connection) Write(b []byte) error {
	bufLength := &bytes.Buffer{}
	length := uint32(len(b))

	err := binary.Write(bufLength, binary.BigEndian, length)
	if err != nil {
		return fmt.Errorf("failed writing message length %w", err)
	}

	_, err = m.Conn.Write(bufLength.Bytes())
	if err != nil {
		return fmt.Errorf("failed writing bytes to connection: %w", err)
	}

	_, err = m.Conn.Write(b)
	if err != nil {
		return fmt.Errorf("failed write bytes to connection: %w", err)
	}

	return nil
}
