package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/mjwhitta/cli"
	"gitlab.com/mjwhitta/log"
	sak "gitlab.com/mjwhitta/nutsak"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r.(error).Error())
			}
			log.ErrX(Exception, r.(error).Error())
		}
	}()

	var e error
	var lefty sak.NUt
	var righty sak.NUt
	var sig = make(chan os.Signal, 1)
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
		sak.LogLvl = int(flags.debug)
	}

	// Create first NUt
	if lefty, e = sak.NewNUt(cli.Arg(0)); e != nil {
		panic(e)
	}

	// Create second NUt
	if cli.NArg() == 2 {
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

		lefty.Down()
		righty.Down()
		os.Exit(130)
	}()

	// Pair NUts to create two-way tunnel
	if e = sak.Pair(lefty, righty); e != nil {
		panic(e)
	}
}
