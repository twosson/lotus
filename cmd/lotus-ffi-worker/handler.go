package main

import (
	"context"
	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"github.com/filecoin-project/specs-storage/storage"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

type FFIWorkerHandler struct {
	running bool
}

func NewFFIWorkerHandler() (*FFIWorkerHandler, error) {
	return &FFIWorkerHandler{}, nil
}

func (h *FFIWorkerHandler) SealCommit2(ctx context.Context, req FFIRequest) (FFIResponse, error) {
	rsp := FFIResponse{}
	if h.running {
		return rsp, xerrors.Errorf("Working...")
	}

	h.running = true

	sector := storage.SectorRef{
		ID: abi.SectorID{
			Miner:  abi.ActorID(req.MinerNumber),
			Number: abi.SectorNumber(req.SectorNumber),
		},
		ProofType: req.ProofType,
	}

	vm, closer, err := NewVMWorkerRPC(context.Background(), req.WorkerAddress)
	if err != nil {
		if closer != nil {
			closer()
		}
		return rsp, err
	}

	// 获取phase1Out
	go func() {
		defer closer()
		id, err := uuid.Parse(req.CallID)
		if err != nil {
			errCall := &storiface.CallError{
				Code:    1,
				Message: "parse uuid error",
			}
			_ = vm.ReturnCommit2(ctx, storiface.UndefCall, nil, errCall)
		}
		callId := storiface.CallID{
			Sector: abi.SectorID{
				Miner:  abi.ActorID(req.MinerNumber),
				Number: abi.SectorNumber(req.SectorNumber),
			},
			ID: id,
		}
		c1o, err := vm.GetCommit1(ctx, callId)
		if err != nil {
			errCall := &storiface.CallError{
				Code:    1,
				Message: "get commit1 out error",
			}
			_ = vm.ReturnCommit2(ctx, callId, nil, errCall)
		}

		result, err := ffi.SealCommitPhase2(c1o, sector.ID.Number, sector.ID.Miner)
		h.running = false
		if err != nil {
			errCall := &storiface.CallError{
				Code:    1,
				Message: "ffi commit2 error",
			}
			_ = vm.ReturnCommit2(ctx, callId, nil, errCall)
		} else {
			errCall := &storiface.CallError{
				Code:    0,
				Message: "ffi commit2 success",
			}
			_ = vm.ReturnCommit2(ctx, callId, result, errCall)
		}
	}()
	return rsp, nil
}

var _ FFIWorkerApi = (*FFIWorkerHandler)(nil)
