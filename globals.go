package nutsak

import "github.com/mjwhitta/log"

// Version is the package version.
const Version string = "1.1.9"

//nolint:grouper // This is an iota block
const (
	modeClient = iota
	modeServer
)

var (
	// Logger will be used to log information deemed relevant to the
	// user.
	Logger *log.Messenger

	// LogLvl will be used to determine the amount of log messages
	// displayed.
	LogLvl int
)
