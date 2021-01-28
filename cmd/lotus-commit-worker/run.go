package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	paramfetch "github.com/filecoin-project/go-paramfetch"
	"github.com/filecoin-project/go-statestore"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/build"
	lcli "github.com/filecoin-project/lotus/cli"
	sectorstorage "github.com/filecoin-project/lotus/extern/sector-storage"
	"github.com/filecoin-project/lotus/extern/sector-storage/sealtasks"
	"github.com/filecoin-project/lotus/extern/sector-storage/stores"
	"github.com/filecoin-project/lotus/lib/rpcenc"
	"github.com/filecoin-project/lotus/metrics"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-datastore/namespace"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
	"net"
	"net/http"
	"strings"
	"time"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start lotus worker",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen",
			Usage: "host address and port the worker api will listen on",
			Value: "0.0.0.0:3456",
		},
		&cli.StringFlag{
			Name:  "proxy",
			Usage: "proxy address and port the worker api will listen on",
			Value: "0.0.0.0:3456",
		},
		&cli.StringFlag{
			Name:  "timeout",
			Usage: "used when 'listen' is unspecified. must be a valid duration recognized by golang's time.ParseDuration function",
			Value: "30m",
		},
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting lotus worker")

		// Connect to storage-miner
		ctx := lcli.ReqContext(cctx)

		var nodeApi api.StorageMiner
		var v api.Version
		var closer func()
		var err error

		for {
			nodeApi, closer, err = lcli.GetStorageMinerAPI(cctx, lcli.StorageMinerUseHttp)
			if err == nil {
				v, err = nodeApi.Version(ctx)
				if err == nil {
					log.Infof("Connected miner api for version: %s\n", v.Version)
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

		if v.APIVersion != build.MinerAPIVersion {
			log.Warnf("miner API version doesn't match: expected: %s", api.Version{APIVersion: build.MinerAPIVersion})
		}

		// Check params
		act, err := nodeApi.ActorAddress(ctx)
		if err != nil {
			return err
		}

		ssize, err := nodeApi.ActorSectorSize(ctx, act)
		if err != nil {
			return err
		}

		if err := paramfetch.GetParams(ctx, build.ParametersJSON(), uint64(ssize)); err != nil {
			return xerrors.Errorf("get params: %w", err)
		}

		// Open repo
		repoPath := cctx.String(FlagWorkerRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}
		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if !ok {
			if err := r.Init(repo.Worker); err != nil {
				return err
			}
			lr, err := r.Lock(repo.Worker)
			if err != nil {
				return err
			}

			// init datastore for r.Exists
			_, err = lr.Datastore("/metadata")
			if err != nil {
				return err
			}

			if err := lr.Close(); err != nil {
				return xerrors.Errorf("close repo: %w", err)
			}
		}

		lr, err := r.Lock(repo.Worker)
		if err != nil {
			return err
		}
		defer func() {
			if err := lr.Close(); err != nil {
				log.Error("closing repo", err)
			}
		}()
		ds, err := lr.Datastore("/metadata")
		if err != nil {
			return err
		}

		log.Info("Opening local storage; connecting to master")

		const unspecifiedAddress = "0.0.0.0"
		address := cctx.String("listen")
		addressSlice := strings.Split(address, ":")
		if ip := net.ParseIP(addressSlice[0]); ip != nil {
			if ip.String() == unspecifiedAddress {
				timeout, err := time.ParseDuration(cctx.String("timeout"))
				if err != nil {
					return err
				}
				rip, err := extractRoutableIP(timeout)
				if err != nil {
					return err
				}
				address = rip + ":" + addressSlice[1]
			}
		}

		proxy := cctx.String("proxy")

		// Create / expose the worker
		wsts := statestore.New(namespace.Wrap(ds, modules.WorkerCallsPrefix))
		var taskTypes []sealtasks.TaskType
		taskTypes = append(taskTypes, sealtasks.TTFetch, sealtasks.TTCommit1, sealtasks.TTCommit2, sealtasks.TTFinalize)
		localWorker := sectorstorage.NewLocalWorker(sectorstorage.WorkerConfig{
			TaskTypes: taskTypes,
			NoSwap:    true,
		}, nil, nil, nodeApi, nodeApi, wsts)

		workerApi := &worker{
			LocalWorker: localWorker,
			ls:          lr,
		}

		mux := mux.NewRouter()
		log.Info("Setting up control endpoint at " + proxy)

		readerHandler, readerServerOpt := rpcenc.ReaderParamDecoder()
		rpcServer := jsonrpc.NewServer(readerServerOpt)
		rpcServer.Register("Filecoin", apistruct.PermissionedWorkerAPI(metrics.MetricedWorkerAPI(workerApi)))

		mux.Handle("/rpc/v0", rpcServer)
		mux.Handle("/rpc/streams/v0/push/{uuid}", readerHandler)
		mux.PathPrefix("/").Handler(http.DefaultServeMux)

		ah := &auth.Handler{
			Verify: nodeApi.AuthVerify,
			Next:   mux.ServeHTTP,
		}

		srv := &http.Server{
			Handler: ah,
			BaseContext: func(listener net.Listener) context.Context {
				ctx, _ := tag.New(context.Background(), tag.Upsert(metrics.APIInterface, "lotus-worker"))
				return ctx
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

		nl, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}

		{
			a, err := net.ResolveTCPAddr("tcp", address)
			if err != nil {
				return xerrors.Errorf("parsing address: %w", err)
			}

			ma, err := manet.FromNetAddr(a)
			if err != nil {
				return xerrors.Errorf("creating api multiaddress: %w", err)
			}

			if err := lr.SetAPIEndpoint(ma); err != nil {
				return xerrors.Errorf("setting api endpoint: %w", err)
			}

			ainfo, err := lcli.GetAPIInfo(cctx, repo.StorageMiner)
			if err != nil {
				return xerrors.Errorf("could not get miner API info: %w", err)
			}

			// TODO: ideally this would be a token with some permissions dropped
			if err := lr.SetAPIToken(ainfo.Token); err != nil {
				return xerrors.Errorf("setting api token: %w", err)
			}
		}

		minerSession, err := nodeApi.Session(ctx)
		if err != nil {
			return xerrors.Errorf("getting miner session: %w", err)
		}

		waitQuietCh := func() chan struct{} {
			out := make(chan struct{})
			go func() {
				workerApi.LocalWorker.WaitQuiet()
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
						return
					}
				}

				log.Errorf("LOTUS-MINER CONNECTION LOST")
			}
		}()

		return srv.Serve(nl)
	},
}
