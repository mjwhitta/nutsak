package nutsak

import (
	"io"
	"net"
	"strings"
	"time"

	"gitlab.com/mjwhitta/errors"
)

// UDPNUt is a UDP network utility.
type UDPNUt struct {
	*baseNUt

	addr     string
	clients  map[string]struct{}
	conn     *net.UDPConn
	connAddr *net.UDPAddr
	mode     int
}

// NewUDPNUt will return a pointer to a UDP network utility instance
// with the provided seed and mode.
func NewUDPNUt(seed string) (*UDPNUt, error) {
	var e error
	var nut *UDPNUt = &UDPNUt{clients: map[string]struct{}{}}

	// Inherit
	nut.baseNUt = super(seed)

	switch nut.Type() {
	case "udp":
		nut.mode = client
	case "udp-l", "udp-listen":
		nut.mode = server
		nut.thetype = "udp-listen"
	default:
		e = errors.Newf("unknown udp type %s", nut.Type())
		return nil, e
	}

	for k, v := range nut.config {
		if k == "addr" {
			nut.addr = v
			if !strings.Contains(nut.addr, ":") {
				nut.addr = "0.0.0.0:" + nut.addr
			}
		} else {
			e = errors.Newf("unknown %s option %s", nut.Type(), k)
			return nil, e
		}
	}

	if nut.addr == "" {
		return nil, errors.Newf("no %s addr provided", nut.Type())
	}

	return nut, nil
}

func (nut *UDPNUt) connect(addr string) error {
	var a *net.UDPAddr
	var e error

	if a, e = net.ResolveUDPAddr("udp", addr); e != nil {
		return errors.Newf("failed to resolve %s: %w", addr, e)
	}

	for {
		if nut.conn, e = net.DialUDP("udp", nil, a); e != nil {
			logErr(1, errors.Newf("connect failed: %w", e).Error())

			time.Sleep(time.Second)
			continue
		}

		// Don't loop if successful connection
		break
	}

	return nil
}

// Down will stop the network utility. In the case of UDP, it will
// close the connection.
func (nut *UDPNUt) Down() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already down
	if !nut.up {
		return nil
	}

	// Down before closing connection
	nut.up = false

	// Close connection
	if nut.conn != nil {
		e = nut.conn.Close()
	}

	return e
}

// KeepAlive will return whether or not the network utility should be
// left running upon EOF. In the case of UDP, it should always return
// true, if it is also up.
func (nut *UDPNUt) KeepAlive() bool {
	return nut.up
}

func (nut *UDPNUt) listen(addr string) error {
	var a *net.UDPAddr
	var e error

	if a, e = net.ResolveUDPAddr("udp", addr); e != nil {
		return errors.Newf("failed to resolve %s: %w", addr, e)
	}

	if nut.conn, e = net.ListenUDP("udp", a); e != nil {
		return errors.Newf("failed to listen on %s: %w", addr, e)
	}

	return nil
}

// Read will read from the current UDP connection.
func (nut *UDPNUt) Read(p []byte) (int, error) {
	var a *net.UDPAddr
	var e error
	var n int

	if !nut.up {
		logSubInfo(2, "%s read: not up", nut.String())
	}

	if nut.conn == nil {
		logSubInfo(2, "%s read: no connection", nut.String())
	}

	if !nut.up || (nut.conn == nil) {
		return 0, io.EOF
	}

	if n, a, e = nut.conn.ReadFromUDP(p); e != nil {
		logSubInfo(2, "%s read: %d bytes", nut.String(), n)

		if !nut.up || (nut.mode == server) {
			e = nil
		}

		return n, e
	}

	if nut.mode == server {
		nut.connAddr = a

		if _, ok := nut.clients[a.String()]; !ok {
			nut.clients[a.String()] = struct{}{}
			logGood(1, "Connection from %s", a.String())
		}
	}

	logSubInfo(2, "%s read: %d bytes", nut.String(), n)
	return n, nil
}

// Up will start the network utility. In the case of UDP, it will
// either connect or listen, depending on the mode.
func (nut *UDPNUt) Up() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already up
	if nut.up {
		return nil
	}

	// Create connection/listener
	switch nut.mode {
	case client:
		e = nut.connect(nut.addr)
	case server:
		e = nut.listen(nut.addr)
	}

	if e == nil {
		nut.up = true
	}

	return e
}

// Write will write to the current UDP connection.
func (nut *UDPNUt) Write(p []byte) (int, error) {
	var e error
	var n int

	if !nut.up {
		logSubInfo(2, "%s write: not up", nut.String())
	}

	if nut.conn == nil {
		logSubInfo(2, "%s write: no connection", nut.String())
	}

	if !nut.up || (nut.conn == nil) {
		return 0, io.EOF
	}

	if nut.mode == client {
		n, e = nut.conn.Write(p)
		logSubInfo(2, "%s write: %d bytes", nut.String(), n)

		if !nut.up {
			e = nil
		}

		return n, e
	}

	if nut.connAddr == nil {
		logSubInfo(2, "%s write: no client connection", nut.String())
		return len(p), nil
	}

	n, e = nut.conn.WriteToUDP(p, nut.connAddr)
	logSubInfo(2, "%s write: %d bytes", nut.String(), n)

	if !nut.up {
		e = nil
	}

	return n, e
}
