package main

import (
	"context"
	"github.com/filecoin-project/go-state-types/abi"
)

type FFIResponse struct {
	State      uint64 `json:"state"`
	Algorithm  string `json:"algorithm"`
	InDataTime int64  `json:"inDataTime"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
}

type FFIRequest struct {
	CallID        string                  `json:"callId"`
	MinerNumber   uint64                  `json:"minerNumber"`
	SectorNumber  uint64                  `json:"sectorNumber"`
	WorkerName    string                  `json:"workerName"`
	WorkerAddress string                  `json:"workerAddress"`
	GPUModel      string                  `json:"gpuModel"`
	CPUModel      string                  `json:"cpuModel"`
	ProofType     abi.RegisteredSealProof `json:"proofType"`
}

type FFIWorkerApi interface {
	SealCommit2(ctx context.Context, req FFIRequest) (FFIResponse, error)
}
