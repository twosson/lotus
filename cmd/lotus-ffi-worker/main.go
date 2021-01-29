package main

import (
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/lib/lotuslog"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"os"
)

var log = logging.Logger("main")

func main() {
	build.RunningNodeType = build.NodeWorker
	lotuslog.SetupLogLevels()
	local := []*cli.Command{
		runCmd,
	}

	app := &cli.App{
		Name:     "lotus-commit",
		Usage:    "Remote miner worker for commit",
		Version:  build.UserVersion(),
		Commands: local,
	}

	app.Setup()

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		return
	}

}
