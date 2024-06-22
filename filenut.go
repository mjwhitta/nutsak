package nutsak

import (
	"io"
	"os"

	"github.com/mjwhitta/errors"
)

// FileNUt is a file network utility.
type FileNUt struct {
	*baseNUt

	addr string
	file *os.File
	mode string
}

// NewFileNUt will return a pointer to a file network utility
// instance with the provided seed.
func NewFileNUt(seed string) (*FileNUt, error) {
	var e error
	var nut *FileNUt = &FileNUt{}

	// Inherit
	nut.baseNUt = super(seed)

	switch nut.Type() {
	case "file":
		nut.mode = "read"
	default:
		return nil, errors.Newf("unknown file type %s", nut.Type())
	}

	for k, v := range nut.config {
		if k == "addr" {
			nut.addr = v
		} else if k == "mode" {
			switch v {
			case "":
			case "append", "read", "write":
				nut.mode = v
			default:
				e = errors.Newf("unknown %s mode %s", nut.Type(), v)
				return nil, e
			}
		} else {
			e = errors.Newf("unknown %s option %s", nut.Type(), k)
			return nil, e
		}
	}

	if nut.addr == "" {
		return nil, errors.Newf("no %s name provided", nut.Type())
	}

	return nut, nil
}

// Down will stop the network utility. In the case of file, it will
// close the file.
func (nut *FileNUt) Down() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already down
	if !nut.up {
		return nil
	}

	nut.up = false

	// Close file
	if nut.file != nil {
		e = nut.file.Close()
	}

	return e
}

// KeepAlive will return whether or not the network utility should be
// left running upon EOF. In the case of a file, it is dependent
// upon mode.
func (nut *FileNUt) KeepAlive() bool {
	switch nut.mode {
	case "append", "write":
		return nut.up
	}

	return false
}

// Read will read from the file.
func (nut *FileNUt) Read(p []byte) (int, error) {
	var e error
	var n int

	if !nut.up {
		logSubInfo(2, "%s read: not up", nut.String())
	}

	if nut.file == nil {
		logSubInfo(2, "%s read: file not open", nut.String())
	}

	if !nut.up || (nut.file == nil) {
		return 0, io.EOF
	}

	switch nut.mode {
	case "append", "write":
		// Nothing to read
		return 0, nil
	}

	n, e = nut.file.Read(p)
	logSubInfo(2, "%s read: %d bytes", nut.String(), n)

	if !nut.up {
		e = nil
	}

	return n, e
}

// Up will start the network utility. In the case of file, it will
// open the file with the specified mode.
func (nut *FileNUt) Up() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already up
	if nut.up {
		return nil
	}

	// Open file
	switch nut.mode {
	case "append":
		nut.file, e = os.OpenFile(
			nut.addr,
			os.O_APPEND|os.O_CREATE|os.O_RDWR,
			0o666,
		)
		_, _ = nut.file.Seek(0, 2)
	case "read":
		nut.file, e = os.Open(nut.addr)
	case "write":
		nut.file, e = os.Create(nut.addr)
	}

	if e == nil {
		nut.up = true
		logGood(1, "opened %s to %s", nut.addr, nut.mode)
	}

	return e
}

// Write will write to the file.
func (nut *FileNUt) Write(p []byte) (int, error) {
	var e error
	var n int

	if !nut.up {
		logSubInfo(2, "%s write: not up", nut.String())
	}

	if nut.file == nil {
		logSubInfo(2, "%s write: file not open", nut.String())
	}

	if !nut.up || (nut.file == nil) {
		return 0, io.EOF
	}

	switch nut.mode {
	case "read":
		// Send to /dev/null
		return len(p), nil
	}

	n, e = nut.file.Write(p)
	logSubInfo(2, "%s write: %d bytes", nut.String(), n)

	if !nut.up {
		e = nil
	}

	return n, e
}
