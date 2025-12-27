package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/log"
	sak "github.com/mjwhitta/nutsak"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	var e error
	var lefty sak.NUt
	var righty sak.NUt
	var sig chan os.Signal = make(chan os.Signal, 1)
	var tmp string = "-"

	validate()

	if !flags.quiet {
		if !flags.nsfw {
			sak.Banner()
		} else {
			sak.BannerNSFW()
		}
	}

	// Debug setup
	if flags.debug > 0 {
		sak.Logger = log.NewMessenger()
		//nolint:gosec // G115 - very unlikely
		sak.LogLvl = int(flags.debug)
	}

	// Create first NUt
	if lefty, e = sak.NewNUt(cli.Arg(0)); e != nil {
		panic(e)
	}

	// Create second NUt
	if cli.NArg() == 2 { //nolint:mnd // 2 cli args
		tmp = cli.Arg(1)
	}

	if righty, e = sak.NewNUt(tmp); e != nil {
		panic(e)
	}

	// Setup ^C trap
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Wait for ^C
		<-sig

		// Don't care about errors on ^C
		_ = lefty.Down()
		_ = righty.Down()

		os.Exit(130) //nolint:mnd // 130 is typical ^C exit status
	}()

	// Pair NUts to create two-way tunnel
	if e = sak.Pair(lefty, righty); e != nil {
		panic(e)
	}
}
