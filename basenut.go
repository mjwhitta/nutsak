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
	thetype string
	up      bool
}

func super(seed string) *baseNUt {
	var nut *baseNUt
	var tmp []string = strings.SplitN(seed, ":", 2)

	nut = &baseNUt{
		config:  map[string]string{"addr": ""},
		lock:    &sync.RWMutex{},
		thetype: strings.ToLower(tmp[0]),
	}

	if len(tmp) == 2 {
		// First config option should be address
		tmp = strings.SplitN(tmp[1], ",", 2)
		nut.config["addr"] = tmp[0]

		// Parse any remaining config options
		if len(tmp) == 2 {
			for _, cfg := range strings.Split(tmp[1], ",") {
				tmp = strings.SplitN(cfg, "=", 2)
				tmp[0] = strings.ToLower(tmp[0])

				if len(tmp) == 1 {
					nut.config[tmp[0]] = ""
				} else {
					nut.config[tmp[0]] = tmp[1]
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
		nut.prIn.Close()
	}
	if nut.prOut != nil {
		nut.prOut.Close()
	}
	if nut.pwIn != nil {
		nut.pwIn.Close()
	}
	if nut.pwOut != nil {
		nut.pwOut.Close()
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

// String will return a string representation of the network utility.
func (nut *baseNUt) String() string {
	var out string = nut.thetype + ":" + nut.config["addr"] + ","

	for k, v := range nut.config {
		if k == "addr" {
			continue
		}

		if v != "" {
			out += k + "=" + v
		} else {
			out += k
		}
		out += ","
	}

	return strings.TrimSuffix(out, ",")
}

// Type will return the type of the network utility.
func (nut *baseNUt) Type() string {
	return nut.thetype
}

// Up is a default that is not implemented. Each NUt will implement.
func (nut *baseNUt) Up() error {
	// Create pipes
	nut.prIn, nut.pwIn = io.Pipe()
	nut.prOut, nut.pwOut = io.Pipe()

	return nil
}
