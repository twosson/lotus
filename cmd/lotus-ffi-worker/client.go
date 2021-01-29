package main

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"github.com/filecoin-project/specs-storage/storage"
	"net/http"
)

type VMWorker struct {
	Internal struct {
		ReturnCommit2 func(ctx context.Context, callID storiface.CallID, proof storage.Proof, callError *storiface.CallError) error
		GetCommit1    func(ctx context.Context, callID storiface.CallID) (storage.Commit1Out, error)
	}
}

func NewVMWorkerRPC(ctx context.Context, endpoint string) (VMWorkerAPI, jsonrpc.ClientCloser, error) {
	var res VMWorker
	closer, err := jsonrpc.NewMergeClient(ctx, endpoint, "Filkeep", []interface{}{&res.Internal}, http.Header{})
	return &res, closer, err
}

func (v *VMWorker) GetCommit1(ctx context.Context, callID storiface.CallID) (storage.Commit1Out, error) {
	return v.Internal.GetCommit1(ctx, callID)
}

func (v *VMWorker) ReturnCommit2(ctx context.Context, callID storiface.CallID, proof storage.Proof, callError *storiface.CallError) error {
	return v.Internal.ReturnCommit2(ctx, callID, proof, callError)
}

type VMWorkerAPI interface {
	ReturnCommit2(ctx context.Context, callID storiface.CallID, proof storage.Proof, callError *storiface.CallError) error
	GetCommit1(ctx context.Context, callID storiface.CallID) (storage.Commit1Out, error)
}

var _ VMWorkerAPI = (*VMWorker)(nil)
