package nutsak

import (
	"io"
	"strings"
	"sync"
)

// baseNUt is a network utility that stores relevant data for all
// network utility types.
type baseNUt struct {
	config map[string]string
	lock   *sync.RWMutex
	//      /  pwIn:prIn  ->  Read() \
	// src {                          } NUt
	//      \ prOut:pwOut <- Write() /
	prIn    *io.PipeReader
	prOut   *io.PipeReader
	pwIn    *io.PipeWriter
	pwOut   *io.PipeWriter
	theType string
	up      bool
}

func super(seed string) *baseNUt {
	var addr string
	var hasOpts bool
	var nut *baseNUt
	var opts string
	var theType string

	theType, opts, hasOpts = strings.Cut(seed, ":")

	nut = &baseNUt{
		config:  map[string]string{"addr": ""},
		lock:    &sync.RWMutex{},
		theType: strings.ToLower(theType),
	}

	if hasOpts {
		// First config option should be address
		addr, opts, hasOpts = strings.Cut(opts, ",")
		nut.config["addr"] = addr

		// Parse any remaining config options
		if hasOpts {
			for _, cfg := range strings.Split(opts, ",") {
				if k, v, ok := strings.Cut(cfg, "="); ok {
					nut.config[strings.ToLower(k)] = v
				} else {
					nut.config[strings.ToLower(k)] = ""
				}
			}
		}
	}

	return nut
}

// Close is an alias for Down().
func (nut *baseNUt) Close() error {
	return nut.Down()
}

// Down is a default that is not implemented. Each NUt will implement.
func (nut *baseNUt) Down() error {
	// Close pipes
	if nut.prIn != nil {
		_ = nut.prIn.Close()
	}

	if nut.prOut != nil {
		_ = nut.prOut.Close()
	}

	if nut.pwIn != nil {
		_ = nut.pwIn.Close()
	}

	if nut.pwOut != nil {
		_ = nut.pwOut.Close()
	}

	return nil
}

// IsUp will return whether or not the network utility is up and
// running.
func (nut *baseNUt) IsUp() bool {
	return nut.up
}

// Open is an alias for Up().
func (nut *baseNUt) Open() error {
	return nut.Up()
}

// String will return a string representation of the baseNUt.
func (nut *baseNUt) String() string {
	var sb strings.Builder

	sb.WriteString(nut.theType + ":" + nut.config["addr"] + ",")

	for k, v := range nut.config {
		if k == "addr" {
			continue
		}

		if v != "" {
			sb.WriteString(k + "=" + v)
		} else {
			sb.WriteString(k)
		}

		sb.WriteString(",")
	}

	return strings.TrimSuffix(sb.String(), ",")
}

// Type will return the type of the network utility.
func (nut *baseNUt) Type() string {
	return nut.theType
}

// Up is a default that is not implemented. Each NUt will implement.
func (nut *baseNUt) Up() error {
	// Create pipes
	nut.prIn, nut.pwIn = io.Pipe()
	nut.prOut, nut.pwOut = io.Pipe()

	return nil
}
