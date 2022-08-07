package nutsak

import (
	"strings"

	"gitlab.com/mjwhitta/errors"
)

// NUt is a network utility capable of reading and writing to a
// network connection.
type NUt interface {
	// BaseNUt
	Close() error
	IsUp() bool
	Open() error
	String() string
	Type() string

	// Needs implemented
	Down() error
	KeepAlive() bool
	Read(p []byte) (n int, e error)
	Up() error
	Write(p []byte) (n int, e error)
}

// NewNUt will return a new network utility from the provided seed
// string.
func NewNUt(seed string) (NUt, error) {
	var tmp []string = strings.SplitN(seed, ":", 2)

	switch tmp[0] {
	case "file":
		return NewFileNUt(seed)
	case "-", "stdin", "stdio", "stdout":
		return NewStdioNUt(seed)
	case "tcp", "tcp-l", "tcp-listen":
		return NewTCPNUt(seed)
	case "tls", "tls-l", "tls-listen":
		return NewTLSNUt(seed)
	case "udp", "udp-l", "udp-listen":
		return NewUDPNUt(seed)
	}

	return nil, errors.Newf("unsupported NUt: %s", tmp[0])
}
