package main

import (
	"os"
	"strings"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	sak "github.com/mjwhitta/nutsak"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	debug   cli.Counter
	nocolor bool
	nsfw    bool
	quiet   bool
	verbose bool
	version bool
}

// SEEDTYPES should be initialized at compile time.
var SEEDTYPES string

func init() {
	// Configure cli package
	cli.Align = true
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = hl.Sprintf("%s [OPTIONS] <seed> [seed]", os.Args[0])
	cli.BugEmail = "nutsak.bugs@whitta.dev"
	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		hl.Sprintf("%d: Invalid option\n", InvalidOption),
		hl.Sprintf("%d: Missing option\n", MissingOption),
		hl.Sprintf("%d: Invalid argument\n", InvalidArgument),
		hl.Sprintf("%d: Missing argument\n", MissingArgument),
		hl.Sprintf("%d: Extra argument\n", ExtraArgument),
		hl.Sprintf("%d: Exception", Exception),
	)
	cli.Info(
		")( Sak is a networking utility, similar to socat. It's",
		"written in Go using OO the Network Utility Swiss-Army Knife",
		"(NUtSAK) library.\n\n",
		"If only one seed is provided, the other is assumed to be",
		"STDIO.",
	)
	cli.Section(
		"SEED SPECIFICATIONS",
		"With the seed command line arguments, the user gives sak",
		"instructions and the necessary information for establishing",
		"the byte streams.\n\n",
		"A seed specification is usually of the form",
		"proto:addr[,opt]. A seed is generally case-insensitive.",
		"Each supported protocol is detailed below. The address and",
		"the optional comma-separated options are specific to each",
		"protocol.\n\n",
		"Zero or more address options may be given with each",
		"address. They influence the address in some ways. Options",
		"consist of a keyword and an optional value, separated by an",
		"\"=\".",
	)
	cli.Section(
		"SEED TYPES",
		"This section describes the available seed types with their",
		"keywords, options, and semantics.\n",
		strings.ReplaceAll(SEEDTYPES, "\\n", "\n"),
	)
	cli.SeeAlso = []string{"nc", "ncat", "socat"}
	cli.Title = "NUtSAK"

	// Parse cli flags
	cli.Flag(
		&flags.debug,
		"d",
		"debug",
		"Show additional levels of debug messages.",
	)
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(&flags.nsfw, "nsfw", false, "Show NSFW banner.", true)
	cli.Flag(&flags.quiet, "q", "quiet", false, "Do not show banner.")
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Flag(&flags.nsfw, "xxx", false, "Show NSFW banner.", true)
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit, if version was requested
	if flags.version {
		hl.Printf("sak version %s\n", sak.Version)
		os.Exit(Good)
	}

	// Validate cli flags
	if cli.NArg() < 1 {
		cli.Usage(MissingArgument)
	} else if cli.NArg() > 2 {
		cli.Usage(ExtraArgument)
	}
}
