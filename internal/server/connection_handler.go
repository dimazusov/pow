package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"

	"github.com/dimazusov/pow/internal/connection"
	"github.com/dimazusov/pow/internal/server/handler"
	"github.com/dimazusov/pow/internal/storage"
)

type Handler interface {
	Handle(requestMsg []byte) error
}

type ConnectionHandler struct {
	conn            *connection.Connection
	pow             POW
	dict            *storage.Dict
	initReqHandler  *handler.InitReqHandler
	chlngReqHandler *handler.ChlngReqHandler
}

func NewConnectionHandler(
	conn *connection.Connection,
	pow POW,
	chlngStrg *storage.Value,
	dict *storage.Dict,
) *ConnectionHandler {
	return &ConnectionHandler{
		conn:            conn,
		pow:             pow,
		dict:            dict,
		initReqHandler:  handler.NewInitReqHandler(conn, pow, chlngStrg),
		chlngReqHandler: handler.NewChlngReqHandler(conn, pow, chlngStrg, dict),
	}
}

func (m *ConnectionHandler) Handle(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("[PANIC]", "error", err)
		}
	}()

	msgCh := make(chan []byte)
	defer close(msgCh)

	wg := sync.WaitGroup{}
	wg.Add(2) //nolint:mnd,nolintlint

	go func() {
		defer wg.Done()

		for {
			b, err := m.conn.Read()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				slog.Error("[PANIC]", "error", err)

				return
			}
			msgCh <- b
		}
	}()

	go func() {
		defer wg.Done()
		defer m.conn.Close()

		handlers := []Handler{
			m.initReqHandler,
			m.chlngReqHandler,
		}
		for _, handler := range handlers {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgCh:
				err := handler.Handle(msg)
				if err != nil {
					slog.Error("failed handle msg", "error", err)
				}
			}
		}
	}()

	wg.Wait()
}
