package connection

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/dimazusov/pow/internal/config"
	"github.com/stretchr/testify/require"
)

func TestReadTimeout(t *testing.T) {
	address := "0.0.0.0:3000"
	ls := net.ListenConfig{KeepAlive: time.Second * 3}

	listener, err := ls.Listen(context.Background(), "tcp", address)
	require.NoError(t, err)

	defer listener.Close()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		conn, err := listener.Accept()
		require.NoError(t, err)

		c := New(conn)
		require.NoError(t, c.Write([]byte{1}))
	}()

	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)

		c := New(conn)

		c.SetTimeouts(&config.ConnectionConfig{
			ReadTimeout: time.Millisecond * 100,
		})

		time.Sleep(200 * time.Millisecond)
		_, err = c.Read()
		require.ErrorContains(t, err, "i/o timeout")
	}()

	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", address)
		require.NoError(t, err)

		c := New(conn)

		c.SetTimeouts(&config.ConnectionConfig{
			ReadTimeout: time.Millisecond * 200,
		})

		time.Sleep(100 * time.Millisecond)
		b, err := c.Read()
		require.NoError(t, err)
		require.Len(t, b, 1)
	}()

	wg.Wait()
}
