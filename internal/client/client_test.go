package client

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/dimazusov/pow/api/gen/pb-go/pb/scheme"
	"github.com/dimazusov/pow/internal/config"
	"github.com/dimazusov/pow/internal/connection"
	"github.com/dimazusov/pow/internal/hashcash"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestClient_GetRandomQuote(t *testing.T) {
	address := "0.0.0.0:3000"
	var difficulty uint8 = 3
	expectedQuote := "Never break your promises. Keep every promise;"

	ls := net.ListenConfig{KeepAlive: time.Second * 10}

	listener, err := ls.Listen(context.Background(), "tcp", address)
	require.NoError(t, err)
	defer listener.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		runServer(t, listener, difficulty, expectedQuote)
	}()

	go func() {
		defer wg.Done()

		cl := New(&config.ClientConfig{
			ClientAddress: address,
			Connection:    config.ConnectionConfig{},
		})

		givenQuote, err := cl.GetRandomQuote()
		require.NoError(t, err)
		require.Equal(t, expectedQuote, givenQuote)
	}()

	wg.Wait()
}

func runServer(t *testing.T, listener net.Listener, difficulty uint8, expectedQuote string) {
	conn, err := listener.Accept()
	require.NoError(t, err)

	defer conn.Close()

	c := connection.New(conn)

	b, err := c.Read()
	require.NoError(t, err)

	initReq := &scheme.InitRequest{}
	require.NoError(t, proto.Unmarshal(b, initReq))
	require.Equal(t, int(initReq.GetProtocolVersion()), ProtocolVersion)

	hc := hashcash.New(difficulty)

	challenge, err := hc.Challenge()
	require.NoError(t, err)
	require.Len(t, challenge, int(difficulty))

	initRes := &scheme.InitResponse{
		Challenge:  challenge,
		Difficulty: int32(difficulty),
	}
	b, err = proto.Marshal(initRes)
	require.NoError(t, err)

	require.NoError(t, c.Write(b))

	b, err = c.Read()
	require.NoError(t, err)

	req := &scheme.ChallengeRequest{}
	require.NoError(t, proto.Unmarshal(b, req))

	require.True(t, hc.Verify(challenge, req.GetNonce()))

	chlngRes := &scheme.ChallengeResponse{Quote: expectedQuote}
	b, err = proto.Marshal(chlngRes)
	require.NoError(t, err)

	require.NoError(t, c.Write(b))
}
