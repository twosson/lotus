package main

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	paramfetch "github.com/filecoin-project/go-paramfetch"
	"github.com/filecoin-project/lotus/build"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"net"
	"net/http"
)

const (
	ss2KiB   = 2 << 10
	ss8MiB   = 8 << 20
	ss512MiB = 512 << 20
	ss32GiB  = 32 << 30
	ss64GiB  = 64 << 30
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start lotus vm worker",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "listen",
			Usage:   "host address and port the worker api will listen on",
			Value:   "0.0.0.0:8888",
			EnvVars: []string{"LISTEN"},
		},
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting lotus ffi worker")

		ctx := lcli.ReqContext(cctx)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if err := paramfetch.GetParams(ctx, build.ParametersJSON(), uint64(ss32GiB)); err != nil {
			return xerrors.Errorf("get params: %w", err)
		}

		if err := paramfetch.GetParams(ctx, build.ParametersJSON(), uint64(ss64GiB)); err != nil {
			return xerrors.Errorf("get params: %w", err)
		}

		vm, closer, err := NewVMWorkerRPC(context.Background(), "")
		if err != nil {
			if closer != nil {
				closer()
			}
			return err
		}
		defer closer()

		handler, err := NewFFIWorkerHandler(vm)
		if err != nil {
			return err
		}

		router := mux.NewRouter()
		rpcServer := jsonrpc.NewServer()
		rpcServer.Register("Filkeep", handler)
		router.Handle("/rpc/v0", rpcServer)
		router.PathPrefix("/").Handler(http.DefaultServeMux)

		header := &auth.Handler{
			Next: router.ServeHTTP,
		}

		srv := &http.Server{
			Handler: header,
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

		address := cctx.String("listen")

		nl, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}

		return srv.Serve(nl)
	},
}
