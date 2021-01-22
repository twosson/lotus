package main

import (
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/lib/lotuslog"
	"github.com/filecoin-project/lotus/node/repo"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

var log = logging.Logger("main")

const FlagWorkerRepo = "worker-repo"

func main() {
	build.RunningNodeType = build.NodeWorker
	lotuslog.SetupLogLevels()
	local := []*cli.Command{
		runCmd,
	}

	app := &cli.App{
		Name:    "lotus-commit",
		Usage:   "Remote miner worker for commit",
		Version: build.UserVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagWorkerRepo,
				EnvVars: []string{"LOTUS_WORKER_PATH"},
				Value:   "~/.lotusworker", // TODO: Consider XDG_DATA_HOME
				Usage:   "Specify worker repo path",
			},
		},
		Commands: local,
	}

	app.Setup()
	app.Metadata["repoType"] = repo.Worker

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		return
	}
}
