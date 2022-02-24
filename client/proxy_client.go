package aceiclient

import (
	"context"

	"encoding/hex"
	"github.com/daotl/go-log/v2"
	ssrv "github.com/daotl/guts/service/suture"
	gsync "github.com/daotl/guts/sync"

	"github.com/daotl/go-acei/types"
	"github.com/daotl/go-acei/types/local"
)

type proxyClient struct {
	*ssrv.BaseService
	*localClient

	client Client
	mtx    *gsync.Mutex
	Callback
}

var _ LocalClient = (*proxyClient)(nil)

func NewProxyClient(logger log.StandardLogger, mtx *gsync.Mutex, client Client) (*proxyClient, error) {
	if mtx == nil {
		mtx = new(gsync.Mutex)
	}
	cli := &proxyClient{
		mtx: mtx,
	}
	cli.client = client
	var err error
	cli.BaseService, err = ssrv.NewBaseService(cli.run, logger)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (app *proxyClient) Info(ctx context.Context, req *types.RequestInfo) (*types.ResponseInfo, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.InfoSync(ctx, *req)
}

func (app *proxyClient) Query(ctx context.Context, req *types.RequestQuery,
) (*types.ResponseQuery, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.QuerySync(ctx, *req)
}

func (app *proxyClient) InitLedger(ctx context.Context, req *types.RequestInitLedger,
) (*types.ResponseInitLedger, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.InitLedgerSync(ctx, *req)
}

func (app *proxyClient) CheckTx(ctx context.Context, req *local.RequestNativeCheckTx,
) (*local.ResponseNativeCheckTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	res, err := app.client.CheckTxSync(ctx, types.RequestCheckTx{Tx: req.Tx.Bytes, Type: req.Type})
	if err != nil {
		return nil, err
	}
	hexByte, err := hex.DecodeString(res.Sender)
	if err != nil {
		return nil, err
	}
	return &local.ResponseNativeCheckTx{
		Code:         res.Code,
		Data:         res.Data,
		Log:          res.Log,
		Info:         res.Info,
		GasWanted:    res.GasWanted,
		GasUsed:      res.GasUsed,
		Events:       res.Events,
		Codespace:    res.Codespace,
		Sender:       hexByte,
		Priority:     res.Priority,
		MempoolError: nil,
	}, nil
}

func (app *proxyClient) BeginBlock(ctx context.Context, req *local.RequestNativeBeginBlock,
) (*types.ResponseBeginBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	blockHeaderExt := req.GetHeader()
	return app.client.BeginBlockSync(ctx, types.RequestBeginBlock{
		Hash: req.GetHash(),
		Header: types.Header{
			Creator:          blockHeaderExt.Creator,
			Timestamp:        uint64(blockHeaderExt.Time),
			PreviousHashes:   blockHeaderExt.PrevHashes,
			Height:           uint64(blockHeaderExt.Height),
			TransactionsRoot: blockHeaderExt.TxRoot,
			TransactionCount: blockHeaderExt.TxCount,
			Extra:            blockHeaderExt.Extra,
			Signature:        blockHeaderExt.Sig,
		},
		Extra: blockHeaderExt.Extra,
	})
}

func (app *proxyClient) DeliverTx(ctx context.Context, req *local.RequestNativeDeliverTx,
) (*types.ResponseDeliverTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.DeliverTxSync(ctx, types.RequestDeliverTx{Tx: req.Tx.Bytes})
}

func (app *proxyClient) EndBlock(ctx context.Context, req *types.RequestEndBlock,
) (*local.ResponseNativeEndBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	res, err := app.client.EndBlockSync(ctx, *req)
	if err != nil {
		return nil, err
	}
	return &local.ResponseNativeEndBlock{
		Events: res.Events,
		Extra:  res.Extra,
	}, nil
}

func (app *proxyClient) Commit(ctx context.Context) (*types.ResponseCommit, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.CommitSync(ctx)
}

func (app *proxyClient) ListSnapshots(ctx context.Context, req *types.RequestListSnapshots,
) (*types.ResponseListSnapshots, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.ListSnapshotsSync(ctx, *req)
}

func (app *proxyClient) OfferSnapshot(ctx context.Context, req *types.RequestOfferSnapshot,
) (*types.ResponseOfferSnapshot, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.OfferSnapshotSync(ctx, *req)
}

func (app *proxyClient) LoadSnapshotChunk(ctx context.Context, req *types.RequestLoadSnapshotChunk,
) (*types.ResponseLoadSnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.LoadSnapshotChunkSync(ctx, *req)
}

func (app *proxyClient) ApplySnapshotChunk(ctx context.Context, req *types.RequestApplySnapshotChunk,
) (*types.ResponseApplySnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	return app.client.ApplySnapshotChunkSync(ctx, *req)
}
