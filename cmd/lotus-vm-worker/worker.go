package main

import (
	"context"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/extern/sector-storage/sealtasks"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"github.com/filecoin-project/specs-storage/storage"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
	"io"
	"sync/atomic"
	"time"
)

type worker struct {
	job       storiface.CallID
	c1o       storage.Commit1Out
	nodeApi   api.StorageMiner
	commitApi CommitApi
	gpuModel  string
	cpuModel  string
	hostname  string
	listen    string
	session   uuid.UUID
	disabled  int64
	running   bool
}

func newWorker(nodeApi api.StorageMiner, commitApi CommitApi, listen string, gpuModel string, cpuModel string, hostname string) (*worker, error) {
	return &worker{
		nodeApi:   nodeApi,
		commitApi: commitApi,
		listen:    listen,
		gpuModel:  gpuModel,
		cpuModel:  cpuModel,
		hostname:  hostname,
		session:   uuid.New(),
		running:   false,
	}, nil
}

func (w *worker) SealCommit2(ctx context.Context, sector storage.SectorRef, c1o storage.Commit1Out) (storiface.CallID, error) {
	log.Infof("Received sector %d for miner %d", sector.ID.Number, sector.ID.Miner)
	if w.running {
		return storiface.UndefCall, xerrors.Errorf("%s 扇区正在执行任务，请等待 %s 扇区执行完成，在调度任务。", w.job.Sector.Number, w.job.Sector.Number)
	}
	w.running = true
	w.job = storiface.CallID{
		Sector: sector.ID,
		ID:     uuid.New(),
	}
	w.c1o = c1o

	cr := CommitRequest{
		CallID:        w.job.ID.String(),
		MinerNumber:   uint64(sector.ID.Miner),
		SectorNumber:  uint64(sector.ID.Number),
		WorkerName:    w.hostname,
		WorkerAddress: w.listen,
		GPUModel:      w.gpuModel,
		CPUModel:      w.cpuModel,
	}

	if err := w.commitApi.Create(ctx, cr); err != nil {
		w.running = false
		w.job = storiface.UndefCall
		w.c1o = nil
		w.running = false
		return storiface.UndefCall, xerrors.Errorf("%s 扇区执行任务时异常， 请联系 Filkeep 客服. ")
	}
	log.Infof("Create ffi sector %d form miner %d", sector.ID.Number, sector.ID.Miner)
	return w.job, nil
}

func (w *worker) GetCommit1(ctx context.Context, callID storiface.CallID) (storage.Commit1Out, error) {
	if w.c1o == nil {
		return nil, xerrors.Errorf("not found commit1 out")
	}
	return w.c1o, nil
}

func (w *worker) ReturnCommit2(ctx context.Context, callID storiface.CallID, proof storage.Proof, callError *storiface.CallError) error {
	if callID.ID.String() != w.job.ID.String() {
		return xerrors.Errorf("return commit2 proof error.")
	}

	w.c1o = nil
	w.job = storiface.UndefCall
	w.running = false

	cr := CommitRequest{
		CallID:        callID.ID.String(),
		MinerNumber:   uint64(callID.Sector.Miner),
		SectorNumber:  uint64(callID.Sector.Number),
		WorkerName:    w.hostname,
		WorkerAddress: w.listen,
		GPUModel:      w.gpuModel,
		CPUModel:      w.cpuModel,
	}

	if err := w.nodeApi.ReturnSealCommit2(ctx, callID, proof, callError); err != nil {
		if callError.Code > 0 {
			err = w.commitApi.Failed(ctx, cr)
		} else {
			err = w.commitApi.Successful(ctx, cr)
		}
		return err
	}
	return nil
}

func (w *worker) Version(context.Context) (build.Version, error) {
	return build.WorkerAPIVersion, nil
}

func (w *worker) TaskTypes(context.Context) (map[sealtasks.TaskType]struct{}, error) {
	tasks := make(map[sealtasks.TaskType]struct{})
	tasks[sealtasks.TTCommit2] = struct{}{}
	return tasks, nil
}

func (w *worker) Paths(context.Context) ([]stores.StoragePath, error) {
	paths := make([]stores.StoragePath, 0)
	return paths, nil
}

func (w *worker) Info(context.Context) (storiface.WorkerInfo, error) {
	gpus := make([]string, 0)
	gpus = append(gpus, "Filkeep "+w.gpuModel)

	info := storiface.WorkerInfo{
		Hostname: w.hostname,
		Resources: storiface.WorkerResources{
			MemPhysical: 214748364800,
			MemSwap:     0,
			MemReserved: 0,
			CPUs:        10,
			GPUs:        gpus,
		},
	}
	return info, nil
}

func (w *worker) AddPiece(ctx context.Context, sector storage.SectorRef, pieceSizes []abi.UnpaddedPieceSize, newPieceSize abi.UnpaddedPieceSize, pieceData storage.Data) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) SealPreCommit1(ctx context.Context, sector storage.SectorRef, ticket abi.SealRandomness, pieces []abi.PieceInfo) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) SealPreCommit2(ctx context.Context, sector storage.SectorRef, pc1o storage.PreCommit1Out) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) SealCommit1(ctx context.Context, sector storage.SectorRef, ticket abi.SealRandomness, seed abi.InteractiveSealRandomness, pieces []abi.PieceInfo, cids storage.SectorCids) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) FinalizeSector(ctx context.Context, sector storage.SectorRef, keepUnsealed []storage.Range) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) ReleaseUnsealed(ctx context.Context, sector storage.SectorRef, safeToFree []storage.Range) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) MoveStorage(ctx context.Context, sector storage.SectorRef, types storiface.SectorFileType) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) UnsealPiece(context.Context, storage.SectorRef, storiface.UnpaddedByteIndex, abi.UnpaddedPieceSize, abi.SealRandomness, cid.Cid) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) ReadPiece(context.Context, io.Writer, storage.SectorRef, storiface.UnpaddedByteIndex, abi.UnpaddedPieceSize) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) Fetch(context.Context, storage.SectorRef, storiface.SectorFileType, storiface.PathType, storiface.AcquireMode) (storiface.CallID, error) {
	return storiface.UndefCall, nil
}

func (w *worker) TaskDisable(ctx context.Context, tt sealtasks.TaskType) error {
	return nil
}

func (w *worker) TaskEnable(ctx context.Context, tt sealtasks.TaskType) error {
	return nil
}

func (w *worker) Remove(ctx context.Context, sector abi.SectorID) error {
	return nil
}

func (w *worker) StorageAddLocal(ctx context.Context, path string) error {
	return nil
}

func (w *worker) SetEnabled(ctx context.Context, enabled bool) error {
	disabled := int64(1)
	if enabled {
		disabled = 0
	}
	atomic.StoreInt64(&w.disabled, disabled)
	return nil
}

func (w *worker) Enabled(ctx context.Context) (bool, error) {
	return atomic.LoadInt64(&w.disabled) == 0, nil
}

func (w *worker) WaitQuiet(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 1)
	for _ = range ticker.C {
		if w.running == false {
			break
		}
	}
	return nil
}

func (w *worker) ProcessSession(context.Context) (uuid.UUID, error) {
	return w.session, nil
}

func (w *worker) Session(context.Context) (uuid.UUID, error) {
	if atomic.LoadInt64(&w.disabled) == 1 {
		return uuid.UUID{}, xerrors.Errorf("worker disabled")
	}
	return w.session, nil
}

var _ api.WorkerAPI = (*worker)(nil)
