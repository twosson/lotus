package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/build"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	"github.com/filecoin-project/lotus/lib/rpcenc"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"net"
	"net/http"
	"strings"
	"time"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start lotus vm worker",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen",
			Usage: "host address and port the worker api will listen on",
			Value: "0.0.0.0:5685",
		},
		&cli.StringFlag{
			Name:  "maddr",
			Usage: "set commit manager address",
			Value: "0.0.0.0:5685",
		},
		&cli.StringFlag{
			Name:  "proxy",
			Usage: "proxy address and port the worker api will listen on",
			Value: "218.17.24.91:5685",
		},
		&cli.StringFlag{
			Name:  "gpumodel",
			Usage: "set gpu model",
			Value: "RTX 3090",
		},
		&cli.StringFlag{
			Name:  "cpumodel",
			Usage: "set cpu model",
			Value: "Inter",
		},
		&cli.StringFlag{
			Name:  "hostname",
			Usage: "set hostname",
			Value: "filkeep-f010491-001",
		},
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting lotus vm worker")

		ctx := lcli.ReqContext(cctx)
		var nodeApi api.StorageMiner
		var version api.Version
		var closer func()
		var err error

		for {
			nodeApi, closer, err = lcli.GetStorageMinerAPI(cctx, lcli.StorageMinerUseHttp)
			if err == nil {
				version, err = nodeApi.Version(ctx)
				if err == nil {
					log.Infof("Connected miner api for version: %s\n", version.Version)
					break
				}
			}
			fmt.Printf("\r\x1b[0KConnecting to miner API... (%s)", err)
			time.Sleep(time.Second)
			continue
		}

		defer closer()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if version.APIVersion != build.MinerAPIVersion {
			log.Warnf("miner API version doesn't match: expected: %s", api.Version{APIVersion: build.MinerAPIVersion})
		}

		// Check params
		//act, err := nodeApi.ActorAddress(ctx)
		//if err != nil {
		//	return err
		//}
		//
		//ssize, err := nodeApi.ActorSectorSize(ctx, act)
		//if err != nil {
		//	return err
		//}
		//
		//if err := paramfetch.GetParams(ctx, build.ParametersJSON(), uint64(ssize)); err != nil {
		//	return xerrors.Errorf("get params: %w", err)
		//}

		mux := mux.NewRouter()

		proxy := cctx.String("proxy")
		address := cctx.String("listen")
		addressSlice := strings.Split(address, ":")
		gpuModel := cctx.String("gpumodel")
		cpuModel := cctx.String("cpumodel")
		hostname := cctx.String("hostname")
		maddr := cctx.String("maddr")
		log.Infof("lotus vm worker proxy address: %s", proxy)

		commitApi, commitCloser, err := NewCommitApiRpc(ctx, maddr)
		if err != nil {
			if commitCloser != nil {
				commitCloser()
			}
		}
		if commitCloser != nil {
			defer commitCloser()
		}

		w, err := newWorker(nodeApi, commitApi, address, gpuModel, cpuModel, hostname)
		if err != nil {
			return err
		}

		readerHandler, readerServerOpt := rpcenc.ReaderParamDecoder()
		rpcServer := jsonrpc.NewServer(readerServerOpt)
		rpcServer.Register("Filecoin", apistruct.PermissionedWorkerAPI(w))
		rpcServer.Register("Filkeep", w)

		mux.Handle("/rpc/v0", rpcServer)
		mux.Handle("/rpc/streams/v0/push/{uuid}", readerHandler)
		mux.PathPrefix("/").Handler(http.DefaultServeMux)

		authHeader := &auth.Handler{
			Verify: nodeApi.AuthVerify,
			Next:   mux.ServeHTTP,
		}

		srv := &http.Server{
			Handler: authHeader,
			BaseContext: func(listener net.Listener) context.Context {
				return context.Background()
			},
		}

		go func() {
			<-ctx.Done()
			log.Warn("Shutting down...")
			if err := srv.Shutdown(context.TODO()); err != nil {
				log.Errorf("shutting down RPC server failed: %s", err)
			}
			log.Warn("Graceful shutdown successful")
		}()

		nl, err := net.Listen("tcp", "0.0.0.0:"+addressSlice[1])
		if err != nil {
			return err
		}

		minerSession, err := nodeApi.Session(ctx)
		if err != nil {
			return xerrors.Errorf("getting miner session: %w", err)
		}

		waitQuietCh := func() chan struct{} {
			out := make(chan struct{})
			go func() {
				_ = w.WaitQuiet(ctx)
				close(out)
			}()
			return out
		}

		go func() {
			heartbeats := time.NewTicker(stores.HeartbeatInterval)
			defer heartbeats.Stop()

			var readyCh chan struct{}
			for {

				if readyCh == nil {
					log.Info("Making sure no local tasks are running")
					readyCh = waitQuietCh()
				}

				for {
					curSession, err := nodeApi.Session(ctx)
					if err != nil {
						log.Errorf("heartbeat: checking remote session failed: %+v", err)
					} else {
						if curSession != minerSession {
							minerSession = curSession
							break
						}
					}

					select {
					case <-readyCh:
						if err := nodeApi.WorkerConnect(ctx, "http://"+proxy+"/rpc/v0"); err != nil {
							log.Errorf("Registering worker failed: %+v", err)
							cancel()
							return
						}

						log.Info("Worker registered successfully, waiting for tasks")

						readyCh = nil
					case <-heartbeats.C:
					case <-ctx.Done():
						return // graceful shutdown
					}
				}

				log.Errorf("LOTUS-MINER CONNECTION LOST")

			}
		}()

		return srv.Serve(nl)

	},
}
