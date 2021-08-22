package main

import (
	"os"
	"time"

	"github.com/alecthomas/kong"

	"github.com/taisph/probe/pkg/addrprobe"
)

type waitCmdCli struct {
	Wait waitCmd `cmd help:"Wait for hosts"`
}

type waitCmd struct {
	Timeout      int      `name:"timeout" help:"Give up after this many seconds." default:"30"`
	Delay        int      `name:"delay" help:"Delay in seconds between connection attempts." default:"5"`
	NetAddresses []string `arg name:"net-address" help:"Comma separated list of network addresses to probe in the format host:port, tcp:host:port or unix:/path/to/socket."`
}

// Run the command
func (cmd *waitCmd) Run(ctx *kong.Context, cfg appCfg) error {
	p := addrprobe.New(addrprobe.Config{Log: cfg.log, Quit: cfg.quit, Delay: time.Duration(cmd.Delay) * time.Second})
	if !p.Run(cmd.NetAddresses, time.Duration(cmd.Timeout)*time.Second) {
		cfg.log.Error().Msg("One or more probes failed")
		os.Exit(1)
	}

	return nil
}
