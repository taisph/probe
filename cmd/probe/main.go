package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/taisph/go_common/pkg/commonlog"
	"github.com/taisph/go_common/pkg/commonprofiler"
	"github.com/taisph/go_common/pkg/commonversion"
)

const appName = "probe"
const appDescription = "Probe waits for a connection to a list of hosts and ports."

type appCliGlobals struct {
	Version kong.VersionFlag `name:"version" help:"Print version information and quit."`
}

type appCli struct {
	appCliGlobals
	commonlog.LogCli
	commonprofiler.ProfilerCli

	waitCmdCli
}

type appCfg struct {
	log  zerolog.Logger
	quit chan bool
}

func main() {
	cli := appCli{}

	ctx := kong.Parse(&cli,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.Vars{"app": appName, "version": commonversion.Print(appName)})

	log := commonlog.New(&cli.LogCliConfig)
	log.Info().Str("version", commonversion.Info()).Str("log_level", log.GetLevel().String()).
		Msgf("Starting %s %s", ctx.Model.Name, ctx.Command())

	prof := commonprofiler.New(&cli.ProfilerCliConfig)
	if err := prof.Begin(); err != nil {
		log.Fatal().Err(err).Msg("Error starting profilers")
	}

	cfg := appCfg{log: log, quit: make(chan bool)}

	var done = make(chan bool)

	go func() {
		err := ctx.Run(ctx, cfg)
		ctx.FatalIfErrorf(err)
		done <- true
	}()

	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-stop:
		log.Info().Msgf("Caught signal: %+v", sig)
		cfg.quit <- true
		<-done
	case <-done:
	}

	if err := prof.End(); err != nil {
		log.Error().Err(err).Msg("Error stopping profilers")
	}

	log.Info().Msg("Exiting")
}
