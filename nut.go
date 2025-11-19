package nutsak

import (
	"strings"

	"github.com/mjwhitta/errors"
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

type nutConstruct func(string) (NUt, error)

// Verify interface compliance at compile time
var (
	_ NUt = (*FileNUt)(nil)
	_ NUt = (*StdioNUt)(nil)
	_ NUt = (*TCPNUt)(nil)
	_ NUt = (*TLSNUt)(nil)
	_ NUt = (*UDPNUt)(nil)

	nutLookup map[string]nutConstruct = map[string]nutConstruct{
		"-":          NewStdioNUt,
		"file":       NewFileNUt,
		"stdin":      NewStdioNUt,
		"stdio":      NewStdioNUt,
		"stdout":     NewStdioNUt,
		"tcp":        NewTCPNUt,
		"tcp-l":      NewTCPNUt,
		"tcp-listen": NewTCPNUt,
		"tls":        NewTLSNUt,
		"tls-l":      NewTLSNUt,
		"tls-listen": NewTLSNUt,
		"udp":        NewUDPNUt,
		"udp-l":      NewUDPNUt,
		"udp-listen": NewUDPNUt,
	}
)

// NewNUt will return a new network utility from the provided seed
// string.
func NewNUt(seed string) (NUt, error) {
	var theType string

	theType, _, _ = strings.Cut(seed, ":")

	if _, ok := nutLookup[theType]; !ok {
		return nil, errors.Newf("unsupported NUt: %s", theType)
	}

	return nutLookup[theType](seed)
}
