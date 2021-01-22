package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/lotus/build"
	sectorstorage "github.com/filecoin-project/lotus/extern/sector-storage"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"github.com/filecoin-project/specs-storage/storage"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"sync/atomic"
)

type worker struct {
	*sectorstorage.LocalWorker

	ls       stores.LocalStorage
	disabled int64
}

func (w *worker) SealCommit2(ctx context.Context, sector storage.SectorRef, c1o storage.Commit1Out) (storiface.CallID, error) {
	fmt.Print(sector.ID)
	fmt.Print(sector.ProofType)
	return storiface.UndefCall, nil
}

func (w *worker) Version(context.Context) (build.Version, error) {
	return build.WorkerAPIVersion, nil
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
	w.LocalWorker.WaitQuiet() // uses WaitGroup under the hood so no ctx :/
	return nil
}

func (w *worker) ProcessSession(ctx context.Context) (uuid.UUID, error) {
	return w.LocalWorker.Session(ctx)
}

func (w *worker) Session(ctx context.Context) (uuid.UUID, error) {
	if atomic.LoadInt64(&w.disabled) == 1 {
		return uuid.UUID{}, xerrors.Errorf("worker disabled")
	}
	return w.LocalWorker.Session(ctx)
}
