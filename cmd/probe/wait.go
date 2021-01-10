package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"

	"github.com/taisph/probe/pkg/addrprobe"
)

type waitCmdCli struct {
	Wait WaitCmd `cmd help:"Wait for hosts"`
}

type WaitCmd struct {
	Timeout int      `name:"timeout" help:"Give up after this many seconds" default:"30"`
	Delay   int      `name:"delay" help:"Delay in seconds between connection attempts" default:"5"`
	Hosts   []string `arg name:"host:port" help:"Comma separated list of host and ports to probe" required`
}

func (cmd *WaitCmd) Run(ctx *kong.Context, cfg appCfg) error {
	p := addrprobe.New(addrprobe.Config{Log: cfg.log, Quit: cfg.quit, Delay: time.Duration(cmd.Delay) * time.Second})
	if !p.Run(cmd.Hosts, time.Duration(cmd.Timeout)*time.Second) {
		cfg.log.Error().Msg("One or more probes failed")
		os.Exit(1)
	}

	return nil
}
