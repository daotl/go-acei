package aceiclient_test

import (
	"context"
	"fmt"
	"github.com/daotl/go-acei/types/local"
	modeltest "github.com/daotl/go-doubl/test"
	"math/rand"
	"testing"
	"time"

	"github.com/daotl/go-log/v2"
	ssrv "github.com/daotl/guts/service/suture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aceiclient "github.com/daotl/go-acei/client"
	"github.com/daotl/go-acei/server"
	"github.com/daotl/go-acei/types"
)

func TestProxyClientBeginBlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := baseApp{}
	logger := log.TestingLogger()

	_, c := setupGrpcClientServer(ctx, t, logger, app)

	resp := make(chan error, 1)
	go func() {
		bx, err := modeltest.Util.ExtendBlockHeader(&modeltest.TestBlockHeader)
		assert.NoError(t, err)
		res, err := c.BeginBlock(ctx, &local.RequestNativeBeginBlock{Header: bx, Extra: nil})
		assert.NoError(t, err)
		err = c.Flush(ctx)
		assert.NoError(t, err)
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

func TestProxyClientCheckTx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := baseApp{}
	logger := log.TestingLogger()

	_, c := setupGrpcClientServer(ctx, t, logger, app)

	resp := make(chan error, 1)
	go func() {
		txx := modeltest.GenRandomTransactionExt()
		res, err := c.CheckTx(ctx, &local.RequestNativeCheckTx{
			Tx:   txx,
			Type: 0,
		})
		assert.NoError(t, err)
		err = c.Flush(ctx)
		assert.NoError(t, err)
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

func setupGrpcClientServer(
	ctx context.Context,
	t *testing.T,
	logger log.StandardLogger,
	app types.Application,
) (ssrv.Service, aceiclient.LocalClient) {
	t.Helper()

	// some port between 20k and 30k
	port := 20000 + rand.Int31()%10000
	addr := fmt.Sprintf("localhost:%d", port)

	s, err := server.NewServer(logger, addr, "socket", app)
	require.NoError(t, err)
	readyCh, sResCh := s.Start(ctx)
	require.NoError(t, <-readyCh)
	go func() { require.NoError(t, <-sResCh) }()
	t.Cleanup(func() {
		if stopped, err := s.Stop(); err == nil {
			<-stopped
		}
	})

	c, err := aceiclient.NewSocketClient(logger, addr, true)
	require.NoError(t, err)
	readyCh, cResCh := c.Start(ctx)
	require.NoError(t, <-readyCh)
	go func() { require.NoError(t, <-cResCh) }()
	t.Cleanup(func() {
		if stopped, err := c.Stop(); err == nil {
			<-stopped
		}
	})

	proxyC, err := aceiclient.NewProxyClient(logger, nil, c)
	require.NoError(t, err)

	return s, proxyC
}

type baseApp struct {
	types.BaseApplication
}

func (baseApp) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	time.Sleep(200 * time.Millisecond)
	return types.ResponseBeginBlock{}
}
