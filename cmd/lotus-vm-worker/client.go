package main

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"net/http"
)

type CommitRequest struct {
	CallID        string `json:"callId"`
	MinerNumber   uint64 `json:"minerNumber"`
	SectorNumber  uint64 `json:"sectorNumber"`
	WorkerName    string `json:"workerName"`
	WorkerAddress string `json:"workerAddress"`
	GPUModel      string `json:"gpuModel"`
	CPUModel      string `json:"cpuModel"`
}

type CommitApi interface {
	Create(context.Context, CommitRequest) error
	Successful(context.Context, CommitRequest) error
	Failed(context.Context, CommitRequest) error
}

type CommitApiStruct struct {
	Internal struct {
		Create     func(context.Context, CommitRequest) error
		Successful func(context.Context, CommitRequest) error
		Failed     func(context.Context, CommitRequest) error
	}
}

func NewCommitApiRpc(ctx context.Context, endpoint string) (CommitApi, jsonrpc.ClientCloser, error) {
	var res CommitApiStruct
	closer, err := jsonrpc.NewMergeClient(ctx, endpoint, "FFI", []interface{}{&res.Internal}, http.Header{})
	return &res, closer, err
}

func (c *CommitApiStruct) Create(ctx context.Context, cr CommitRequest) error {
	return c.Internal.Create(ctx, cr)
}

func (c *CommitApiStruct) Successful(ctx context.Context, cr CommitRequest) error {
	return c.Internal.Successful(ctx, cr)
}

func (c *CommitApiStruct) Failed(ctx context.Context, cr CommitRequest) error {
	return c.Internal.Failed(ctx, cr)
}

var _ CommitApi = (*CommitApiStruct)(nil)
