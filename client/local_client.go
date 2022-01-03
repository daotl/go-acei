package aceiclient

import (
	"context"

	"github.com/daotl/go-log/v2"
	ssrv "github.com/daotl/guts/service/suture"
	gsync "github.com/daotl/guts/sync"

	"github.com/daotl/go-acei/types"
	"github.com/daotl/go-acei/types/local"
)

// NOTE: use defer to unlock mutex because Application might panic (e.g., in
// case of malicious tx or query). It only makes sense for publicly exposed
// methods like CheckTx (/broadcast_tx_* RPC endpoint) or Query (/abci_query
// RPC endpoint), but defers are used everywhere for the sake of consistency.
type localClient struct {
	*ssrv.BaseService

	mtx *gsync.Mutex
	local.Application
	Callback
}

var _ LocalClient = (*localClient)(nil)

// NewLocalClient creates a local client, which will be directly calling the
// methods of the given app.
//
// Both Async and Sync methods ignore the given context.Context parameter.
func NewLocalClient(logger log.StandardLogger, mtx *gsync.Mutex, app local.Application,
) (*localClient, error) {
	if mtx == nil {
		mtx = new(gsync.Mutex)
	}
	cli := &localClient{
		mtx:         mtx,
		Application: app,
	}
	var err error
	cli.BaseService, err = ssrv.NewBaseService(cli.run, logger)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Just a placeholder: localClient won't actually be used as a Suture service
func (cli *localClient) run(ctx context.Context, ready func(error)) error {
	ready(nil)
	<-ctx.Done()
	return nil
}

// TODO: change types.Application to include Error()?
func (app *localClient) Error() error {
	return nil
}

func (app *localClient) Flush(ctx context.Context) {
	return
}

func (app *localClient) Echo(ctx context.Context, msg string) *types.ResponseEcho {
	return &types.ResponseEcho{Message: msg}
}

func (app *localClient) Info(ctx context.Context, req *types.RequestInfo) *types.ResponseInfo {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.Info(req)
}

func (app *localClient) DeliverTx(ctx context.Context, req *local.RequestNativeDeliverTx,
) *types.ResponseDeliverTx {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.DeliverTx(req)
}

func (app *localClient) CheckTx(ctx context.Context, req *local.RequestNativeCheckTx,
) *local.ResponseNativeCheckTx {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.CheckTx(req)
}

func (app *localClient) Query(ctx context.Context, req *types.RequestQuery,
) *types.ResponseQuery {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.Query(req)
}

func (app *localClient) Commit(ctx context.Context) *types.ResponseCommit {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.Commit()
}

func (app *localClient) InitLedger(ctx context.Context, req *types.RequestInitLedger,
) *types.ResponseInitLedger {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.InitLedger(req)
}

func (app *localClient) BeginBlock(ctx context.Context, req *local.RequestNativeBeginBlock,
) *types.ResponseBeginBlock {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.BeginBlock(req)
}

func (app *localClient) EndBlock(ctx context.Context, req *types.RequestEndBlock,
) *local.ResponseNativeEndBlock {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.EndBlock(req)
}

func (app *localClient) ListSnapshots(ctx context.Context, req *types.RequestListSnapshots,
) *types.ResponseListSnapshots {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.ListSnapshots(req)
}

func (app *localClient) OfferSnapshot(ctx context.Context, req *types.RequestOfferSnapshot,
) *types.ResponseOfferSnapshot {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.OfferSnapshot(req)
}

func (app *localClient) LoadSnapshotChunk(ctx context.Context, req *types.RequestLoadSnapshotChunk,
) *types.ResponseLoadSnapshotChunk {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.LoadSnapshotChunk(req)
}

func (app *localClient) ApplySnapshotChunk(ctx context.Context, req *types.RequestApplySnapshotChunk,
) *types.ResponseApplySnapshotChunk {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.Application.ApplySnapshotChunk(req)
}
