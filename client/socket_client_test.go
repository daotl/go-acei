package aceiclient_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ssrv "github.com/daotl/guts/service/suture"

	aceiclient "github.com/daotl/go-acei/client"
	"github.com/daotl/go-acei/server"
	"github.com/daotl/go-acei/types"
)

var ctx = context.Background()

func TestProperSyncCalls(t *testing.T) {
	app := slowApp{}

	s, c := setupClientServer(t, app)
	t.Cleanup(func() {
		if stopped, err := s.Stop(); err != nil {
			t.Error(err)
		} else {
			<-stopped
		}
	})
	t.Cleanup(func() {
		if stopped, err := c.Stop(); err != nil {
			t.Error(err)
		} else {
			<-stopped
		}
	})

	resp := make(chan error, 1)
	go func() {
		// This is BeginBlockSync unrolled....
		reqres, err := c.BeginBlockAsync(ctx, types.RequestBeginBlock{})
		assert.NoError(t, err)
		err = c.FlushSync(context.Background())
		assert.NoError(t, err)
		res := reqres.Response.GetBeginBlock()
		assert.NotNil(t, res)
		resp <- c.Error()
	}()

	select {
	case <-time.After(time.Second):
		require.Fail(t, "No response arrived")
	case err, ok := <-resp:
		require.True(t, ok, "Must not close channel")
		assert.NoError(t, err, "This should return success")
	}
}

func TestHangingSyncCalls(t *testing.T) {
	app := slowApp{}

	s, c := setupClientServer(t, app)
	t.Cleanup(func() {
		if stopped, err := s.Stop(); err != nil {
			t.Log(err)
		} else {
			<-stopped
		}
	})
	t.Cleanup(func() {
		if stopped, err := c.Stop(); err != nil {
			t.Log(err)
		} else {
			<-stopped
		}
	})

	resp := make(chan error, 1)
	go func() {
		// Start BeginBlock and flush it
		reqres, err := c.BeginBlockAsync(ctx, types.RequestBeginBlock{})
		assert.NoError(t, err)
		flush, err := c.FlushAsync(ctx)
		assert.NoError(t, err)
		// wait 20 ms for all events to travel socket, but
		// no response yet from server
		time.Sleep(20 * time.Millisecond)
		// kill the server, so the connections break
		stopped, err := s.Stop()
		assert.NoError(t, err)
		<-stopped

		// wait for the response from BeginBlock
		reqres.Wait()
		flush.Wait()
		resp <- c.Error()
	}()

	select {
	case <-time.After(time.Second):
		require.Fail(t, "No response arrived")
	case err, ok := <-resp:
		require.True(t, ok, "Must not close channel")
		assert.Error(t, err, "We should get EOF error")
	}
}

func setupClientServer(t *testing.T, app types.Application,
) (ssrv.Service, aceiclient.Client) {
	// some port between 20k and 30k
	port := 20000 + rand.Int31()%10000
	addr := fmt.Sprintf("localhost:%d", port)

	s, err := server.NewServer(addr, "socket", app, nil)
	require.NoError(t, err)
	readyCh, sResCh := s.Start(context.Background())
	require.NoError(t, <-readyCh)
	go func() { require.NoError(t, <-sResCh) }()

	c, err := aceiclient.NewSocketClient(addr, true, nil)
	require.NoError(t, err)
	readyCh, cResCh := c.Start(context.Background())
	require.NoError(t, <-readyCh)
	go func() { require.NoError(t, <-cResCh) }()

	return s, c
}

type slowApp struct {
	types.BaseApplication
}

func (slowApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	time.Sleep(200 * time.Millisecond)
	return types.ResponseBeginBlock{}
}
