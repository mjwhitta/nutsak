package nutsak

import "github.com/mjwhitta/log"

const (
	client = iota
	server
)

// Logger will be used to log information deemed relevant to the user.
var Logger *log.Messenger

// LogLvl will be used to determine the amount of log messages
// displayed.
var LogLvl int

// Version is the package version.
const Version string = "1.0.11"
