package nutsak

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/mjwhitta/errors"
)

// TCPNUt is a TCP network utility.
type TCPNUt struct {
	*baseNUt

	addr       string
	conn       *net.TCPConn
	connecting bool
	fork       bool
	list       *net.TCPListener
	mode       int
}

// NewTCPNUt will return a pointer to a TCP network utility instance
// with the provided seed.
func NewTCPNUt(seed string) (*TCPNUt, error) {
	var e error
	var nut *TCPNUt = &TCPNUt{}

	// Inherit
	nut.baseNUt = super(seed)

	switch nut.Type() {
	case "tcp":
		nut.mode = client
	case "tcp-l", "tcl-listen":
		nut.mode = server
		nut.thetype = "tcp-listen"
	default:
		return nil, errors.Newf("unknown tcp type %s", nut.Type())
	}

	// Parse options
	for k, v := range nut.config {
		if k == "addr" {
			nut.addr = v
			if !strings.Contains(nut.addr, ":") {
				nut.addr = "0.0.0.0:" + nut.addr
			}
		} else if (k == "fork") && (nut.mode == server) {
			nut.fork = true
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

func (nut *TCPNUt) connect(addr string) error {
	var a *net.TCPAddr
	var e error

	if a, e = net.ResolveTCPAddr("tcp", addr); e != nil {
		return errors.Newf("failed to resolve %s: %w", addr, e)
	}

	go func() {
		var up = make(chan struct{}, 2)
		var wait = make(chan struct{}, 2)

		for nut.up {
			if nut.conn, e = net.DialTCP("tcp", nil, a); e != nil {
				if nut.up {
					e = errors.Newf("connect failed: %w", e)
					logErr(1, e.Error())
				}

				time.Sleep(time.Second)
				continue
			}

			go func() {
				up <- struct{}{}
				io.Copy(nut.pwIn, nut.conn)
				wait <- struct{}{}
			}()

			go func() {
				up <- struct{}{}
				io.Copy(nut.conn, nut.prOut)
				wait <- struct{}{}
			}()

			// Wait for up
			<-up
			<-up
			time.Sleep(time.Millisecond)

			// Officially up and running
			nut.connecting = false

			// Block
			<-wait
			<-wait
		}
	}()

	return nil
}

// Down will stop the network utility. In the case of TCP, it will
// close the connection or listener, depending on the mode.
func (nut *TCPNUt) Down() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already down
	if !nut.up {
		return nil
	}

	// Down before closing connection/listener and pipes
	nut.connecting = false
	nut.up = false

	// Close connection/listener
	switch nut.mode {
	case client:
		if nut.conn != nil {
			e = nut.conn.Close()
		}
	case server:
		if nut.list != nil {
			e = nut.list.Close()
		}
	}

	// Close pipes to stop io.Copy()
	nut.baseNUt.Down()

	return e
}

// KeepAlive will return whether or not the network utility should be
// left running upon EOF. In the case of TCP, it is dependent upon
// mode.
func (nut *TCPNUt) KeepAlive() bool {
	switch nut.mode {
	case server:
		return nut.up
	}

	return false
}

func (nut *TCPNUt) listen(addr string) error {
	var a *net.TCPAddr
	var c *net.TCPConn
	var e error

	if a, e = net.ResolveTCPAddr("tcp", addr); e != nil {
		return errors.Newf("failed to resolve %s: %w", addr, e)
	}

	if nut.list, e = net.ListenTCP("tcp", a); e != nil {
		return errors.Newf("failed to listen on %s: %w", addr, e)
	}

	go func() {
		var up = make(chan struct{}, 2)
		var wait chan struct{}

		if nut.fork {
			wait = make(chan struct{}, 1)
		} else {
			wait = make(chan struct{}, 2)
		}

		for {
			if c, e = nut.list.AcceptTCP(); e != nil {
				if nut.up {
					e = errors.Newf("connection failed: %w", e)
					logErr(1, e.Error())
				}

				time.Sleep(time.Millisecond)
				continue
			}

			logGood(1, "Connection from %s", c.RemoteAddr().String())

			go func() {
				up <- struct{}{}
				io.Copy(nut.pwIn, c)
				wait <- struct{}{}
			}()

			go func() {
				up <- struct{}{}
				io.Copy(c, nut.prOut)
				wait <- struct{}{}
			}()

			// Wait for up
			<-up
			<-up
			time.Sleep(time.Millisecond)

			// Officially up and running
			nut.connecting = false

			// Block
			if !nut.fork {
				<-wait
			}
			<-wait
		}
	}()

	return nil
}

// Read will read from the current TCP connection.
func (nut *TCPNUt) Read(p []byte) (int, error) {
	var e error
	var n int

	if nut.connecting {
		logSubInfo(2, "%s read: still connecting", nut.String())
	}

	for nut.connecting {
		time.Sleep(time.Millisecond)
	}

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

// Up will start the network utility. In the case of TCP, it will
// either connect or listen, depending on the mode.
func (nut *TCPNUt) Up() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already up
	if nut.up {
		return nil
	}

	// Up after pipes created
	nut.baseNUt.Up()
	nut.connecting = true
	nut.up = true

	// Create connection/listener
	switch nut.mode {
	case client:
		e = nut.connect(nut.addr)
	case server:
		e = nut.listen(nut.addr)
	}

	if e != nil {
		nut.up = false
	}

	return e
}

// Write will write to the current TCP connection.
func (nut *TCPNUt) Write(p []byte) (int, error) {
	var e error
	var n int

	if nut.connecting {
		logSubInfo(2, "%s write: still connecting", nut.String())
	}

	for nut.connecting {
		time.Sleep(time.Millisecond)
	}

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
