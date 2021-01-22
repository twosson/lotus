package main

import (
	"golang.org/x/xerrors"
	"net"
	"os"
	"strings"
	"time"
)

func extractRoutableIP(timeout time.Duration) (string, error) {
	minerMultiAddrKey := "MINER_API_INFO"
	env, ok := os.LookupEnv(minerMultiAddrKey)
	if !ok {
		return "", xerrors.New("MINER_API_INFO environment variable required to extract IP")
	}
	minerAddr := strings.Split(env, "/")
	conn, err := net.DialTimeout("tcp", minerAddr[2]+":"+minerAddr[4], timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	return strings.Split(localAddr.IP.String(), ":")[0], nil
}
