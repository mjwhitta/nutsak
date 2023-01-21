package nutsak

import (
	"io"
	"os"

	"github.com/mjwhitta/errors"
)

// StdioNUt is a stdio network utility.
type StdioNUt struct {
	*baseNUt
}

// NewStdioNUt will return a pointer to a stdio network utility
// instance with the provided seed.
func NewStdioNUt(seed string) (*StdioNUt, error) {
	var e error
	var nut *StdioNUt = &StdioNUt{super(seed)}

	switch nut.Type() {
	case "-":
		nut.thetype = "stdio"
	case "stdio", "stdin", "stdout":
	default:
		e = errors.Newf("unknown stdio type %s", nut.Type())
		return nil, e
	}

	for k := range nut.config {
		if k == "addr" {
			continue
		} else {
			e = errors.Newf("unknown %s option %s", nut.Type(), k)
			return nil, e
		}
	}

	if nut.config["addr"] != "" {
		e = errors.Newf("%s does not need address", nut.Type())
		return nil, e
	}

	return nut, nil
}

// Down will stop the network utility. In the case of stdio, it will
// do nothing.
func (nut *StdioNUt) Down() error {
	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already down
	if !nut.up {
		return nil
	}

	nut.up = false
	nut.baseNUt.Down()

	return nil
}

// KeepAlive will return whether or not the network utility should be
// left running upon EOF. In the case of stdio, it should always
// return true, if it is also up.
func (nut *StdioNUt) KeepAlive() bool {
	return nut.up
}

// Read will read from stdin.
func (nut *StdioNUt) Read(p []byte) (int, error) {
	var e error
	var n int

	if !nut.up {
		logSubInfo(2, "%s read: not up", nut.String())
		return 0, io.EOF
	}

	n, e = nut.prIn.Read(p)
	logSubInfo(2, "%s read: %d bytes", nut.String(), n)

	if !nut.up {
		e = nil
	}

	return n, e
}

// Up will start the network utility. In the case of stdio, it will do
// nothing.
func (nut *StdioNUt) Up() error {
	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already up
	if nut.up {
		return nil
	}

	nut.baseNUt.Up()

	go io.Copy(nut.pwIn, os.Stdin)
	go io.Copy(os.Stdout, nut.prOut)

	nut.up = true
	logGood(1, "opened stdio")

	return nil
}

// Write will write to stdout.
func (nut *StdioNUt) Write(p []byte) (int, error) {
	var e error
	var n int

	if !nut.up {
		logSubInfo(2, "%s write: not up", nut.String())
		return 0, io.EOF
	}

	n, e = nut.pwOut.Write(p)
	logSubInfo(2, "%s write: %d bytes", nut.String(), n)

	if !nut.up {
		e = nil
	}

	return n, e
}
