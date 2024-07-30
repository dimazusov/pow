package handler

import (
	"fmt"
	"log/slog"

	"github.com/dimazusov/pow/api/gen/pb-go/pb/scheme"
	"github.com/dimazusov/pow/internal/connection"
	"github.com/dimazusov/pow/internal/storage"
	"google.golang.org/protobuf/proto"
)

type POWVerifier interface {
	GetDifficulty() uint8
	Verify(challenge []byte, nonce uint64) bool
}

type ChlngReqHandler struct {
	conn         *connection.Connection
	powVerifier  POWVerifier
	chlngStorage *storage.Value
	quotesDict   *storage.Dict
}

func NewChlngReqHandler(
	conn *connection.Connection,
	powVerifier POWVerifier,
	storage *storage.Value,
	dict *storage.Dict,
) *ChlngReqHandler {
	return &ChlngReqHandler{
		conn:         conn,
		powVerifier:  powVerifier,
		chlngStorage: storage,
		quotesDict:   dict,
	}
}

func (m *ChlngReqHandler) Handle(msgReq []byte) error {
	var chlngReq scheme.ChallengeRequest

	err := proto.Unmarshal(msgReq, &chlngReq)
	if err != nil {
		return fmt.Errorf("failed unmarshal msg from pb: %w", err)
	}

	if !m.powVerifier.Verify(m.chlngStorage.Get(), chlngReq.GetNonce()) {
		slog.Debug("verification failed: ", "nonce", chlngReq.GetNonce())
		return fmt.Errorf("wrong: %w", err)
	}

	randomQuote, err := m.quotesDict.GetRandomPhrase()
	if err != nil {
		return fmt.Errorf("failed getting random quote %w", err)
	}

	res := &scheme.ChallengeResponse{
		Quote: randomQuote,
	}

	b, err := proto.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed marshal msg to pb: %w", err)
	}

	if err = m.conn.Write(b); err != nil {
		return fmt.Errorf("failed writing to connection: %w", err)
	}

	return nil
}
