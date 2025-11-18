// Package nutsak (Network Utility Swiss-Army Knife) will provide a
// means of tunneling traffic between two endpoints. Here is sample
// code to create a Pair:
//
//	var a sak.NUt
//	var b sak.NUt
//	var e error
//
//	// Create first NUt
//	if a, e = sak.NewNUt("tcp-listen:4444,fork"); e != nil {
//	    panic(e)
//	}
//
//	// Create second NUt
//	if b, e = sak.NewNUt("stdout"); e != nil {
//	    panic(e)
//	}
//
//	// Pair NUts to create two-way tunnel
//	if e = sak.Pair(a, b); e != nil {
//	    panic(e)
//	}
//
// This will create a TCP listener on port 4444 that forks each new
// connection. Any received data will be written to STDOUT. The
// Network Utilities (NUts) are created from seeds. Below are the
// supported SEED TYPES along with their documentation:
//
// FILE:addr[,mode=(append|read|write)]
//
// This seed takes an address that is an absolute or relative
// filename. This seed is used to read or write a file on disk. The
// default mode is read.
//
// STDIO:
//
// Aliases: -, STDIN, STDOUT
//
// This seed takes no address or options. It can be used to read from
// stdin or write to stdout.
//
// TCP:addr
//
// This seed takes an address of the form [IP:]PORT. The IP is
// optional and defaults to 0.0.0.0 or [::]. This seed is used to make
// an outgoing TCP connection.
//
// TCP-LISTEN:addr[,echo,fork]
//
// Aliases: TCP-L
//
// This seed takes an address of the form [IP:]PORT. The IP is
// optional and defaults to 0.0.0.0 or [::]. This seed is used to
// listen on the provided TCP address. The echo option causes the TCP
// listener to echo the response back to the client. The fork option
// causes the TCP listener to accept multiple connections in parallel.
//
// TLS:addr[,ca=PATH,cert=PATH,key=PATH,verify]
//
// This seed takes an address of the form [IP:]PORT. The IP is
// optional and defaults to 0.0.0.0 or [::]. This seed is used to make
// an outgoing TLS connection. The ca, cert, and key options take a
// filepath (DER or PEM formatted). The verify option determines if
// the server-side CA should be verified. The cert and key options
// must be used together. If verify is specified, a ca must also be
// specified.
//
// TLS-LISTEN:addr[,ca=PATH,cert=PATH,echo,fork,key=PATH,verify]
//
// Aliases: TLS-L
//
// This seed takes an address of the form [IP:]PORT. The IP is
// optional and defaults to 0.0.0.0 or [::]. This seed is used to
// listen on the provided TCP address. The ca, cert, and key options
// take a filepath (DER or PEM formatted). The echo option causes the
// TLS listener to echo the response back to the client. The fork
// option causes the TLS listener to accept multiple connections in
// parallel. The verify option determines if the client-side
// certificate should be verified. The cert and key options are
// mandatory. If verify is specified, a ca must also be specified.
//
// UDP:addr
//
// This seed takes an address of the form [IP:]PORT. The IP is
// optional and defaults to 0.0.0.0 or [::]. This seed is used to make
// an outgoing UDP connection.
//
// UDP-LISTEN:addr[,echo]
//
// Aliases: UDP-L
//
// This seed takes an address of the form [IP:]PORT. The IP is
// optional and defaults to 0.0.0.0 or [::]. This seed is used to
// listen on the provided UDP address. The echo option causes the UDP
// listener to echo the response back to the client.
package nutsak

import (
	"time"
)

// Pair will connect two NUts together using Stream().
func Pair(a NUt, b NUt) error {
	//nolint:mnd // 2 goroutines
	var wait chan struct{} = make(chan struct{}, 2)

	// Ensure they are up
	if e := a.Up(); e != nil {
		return e //nolint:wrapcheck // this error is from the same pkg
	}

	if e := b.Up(); e != nil {
		return e //nolint:wrapcheck // this error is from the same pkg
	}

	// Stream a to b
	go func() {
		stream(a, b)
		time.Sleep(time.Millisecond)

		_ = b.Down()

		wait <- struct{}{}
	}()

	// Stream b to a
	go func() {
		stream(b, a)
		time.Sleep(time.Millisecond)

		_ = a.Down()

		wait <- struct{}{}
	}()

	<-wait
	<-wait

	return nil
}

// Stream will stream data from a to b using io.Copy().
func Stream(a NUt, b NUt) error {
	// Ensure they are up
	if e := a.Up(); e != nil {
		return e //nolint:wrapcheck // this error is from the same pkg
	}

	if e := b.Up(); e != nil {
		return e //nolint:wrapcheck // this error is from the same pkg
	}

	stream(a, b)

	return nil
}
