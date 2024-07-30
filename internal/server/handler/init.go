package handler

import (
	"fmt"
	"log/slog"

	"github.com/dimazusov/pow/api/gen/pb-go/pb/scheme"
	"github.com/dimazusov/pow/internal/connection"
	"github.com/dimazusov/pow/internal/storage"
	"google.golang.org/protobuf/proto"
)

type POWChallanger interface {
	GetDifficulty() uint8
	Challenge() ([]byte, error)
}

type InitReqHandler struct {
	conn      *connection.Connection
	pow       POWChallanger
	chlngStrg *storage.Value
}

func NewInitReqHandler(conn *connection.Connection, pow POWChallanger, chlngStrg *storage.Value) *InitReqHandler {
	return &InitReqHandler{
		conn:      conn,
		pow:       pow,
		chlngStrg: chlngStrg,
	}
}

func (m *InitReqHandler) Handle(msgReq []byte) error {
	var initReq scheme.InitRequest

	err := proto.Unmarshal(msgReq, &initReq)
	if err != nil {
		return fmt.Errorf("failed unmarshal msg from pb: %w", err)
	}

	slog.Debug("read message:", "version", initReq.GetProtocolVersion())

	challenge, err := m.pow.Challenge()
	if err != nil {
		return fmt.Errorf("failed generating challenge: %w", err)
	}

	m.chlngStrg.Set(challenge)

	response := &scheme.InitResponse{
		Challenge:  challenge,
		Difficulty: int32(m.pow.GetDifficulty()),
	}

	err = writeInitResponse(m.conn, response)
	if err != nil {
		return fmt.Errorf("failed writing init response: %w", err)
	}

	return nil
}

func writeInitResponse(conn *connection.Connection, res *scheme.InitResponse) error {
	b, err := proto.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed marshal msg to pb: %w", err)
	}

	if err = conn.Write(b); err != nil {
		return fmt.Errorf("failed writing to connection: %w", err)
	}

	return nil
}
