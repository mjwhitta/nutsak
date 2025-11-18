package nutsak

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"strings"
	"time"

	"github.com/mjwhitta/errors"
)

// TLSNUt is a TLS network utility.
type TLSNUt struct {
	*baseNUt

	addr       string
	ca         *x509.Certificate
	cert       *x509.Certificate
	conn       *tls.Conn
	connecting bool
	echo       bool
	fork       bool
	key        *rsa.PrivateKey
	list       net.Listener
	mode       int
	tlscfg     *tls.Config
	verify     bool
}

// NewTLSNUt will return a pointer to a TLS network utility instance
// with the provided seed.
func NewTLSNUt(seed string) (*TLSNUt, error) {
	var e error
	var nut *TLSNUt = &TLSNUt{}

	// Inherit
	nut.baseNUt = super(seed)

	switch nut.Type() {
	case "tls":
		nut.mode = modeClient
	case "tls-l", "tcl-listen":
		nut.mode = modeServer
		nut.theType = "tls-listen"
	default:
		e = errors.Newf("unknown tls type %s", nut.Type())
		return nil, e
	}

	for k, v := range nut.config {
		if k == "addr" {
			nut.addr = v

			if !strings.Contains(nut.addr, ":") {
				nut.addr = "0.0.0.0:" + nut.addr
			}
		} else if e = nut.parseOpts(k, v); e != nil {
			return nil, e
		}
	}

	if nut.addr == "" {
		return nil, errors.Newf("no %s addr provided", nut.Type())
	}

	if e = nut.setupTLSConfig(); e != nil {
		return nil, e
	}

	return nut, nil
}

func (nut *TLSNUt) connect(addr string) error {
	if _, e := net.ResolveTCPAddr("tcp", addr); e != nil {
		return errors.Newf("failed to resolve %s: %w", addr, e)
	}

	go func() {
		var e error
		//nolint:mnd // 2 goroutines
		var up chan struct{} = make(chan struct{}, 2)
		//nolint:mnd // 2 goroutines
		var wait chan struct{} = make(chan struct{}, 2)

		for nut.up {
			nut.conn, e = tls.Dial("tcp", addr, nut.tlscfg)
			if e != nil {
				if nut.up {
					e = errors.Newf("connect failed: %w", e)
					logErr(1, "%s", e.Error())
				}

				time.Sleep(time.Second)

				continue
			}

			go func() {
				up <- struct{}{}

				_, _ = io.Copy(nut.pwIn, nut.conn)

				wait <- struct{}{}
			}()

			go func() {
				up <- struct{}{}

				_, _ = io.Copy(nut.conn, nut.prOut)

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

// Down will stop the network utility. In the case of TLS, it will
// close the connection or listener, depending on the mode.
func (nut *TLSNUt) Down() error {
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
	case modeClient:
		if nut.conn != nil {
			if e = nut.conn.Close(); e != nil {
				e = errors.Newf("failed to close connection: %w", e)
			}
		}
	case modeServer:
		if nut.list != nil {
			if e = nut.list.Close(); e != nil {
				e = errors.Newf("failed to stop listener: %w", e)
			}
		}
	}

	// Close pipes to stop io.Copy()
	_ = nut.baseNUt.Down()

	return e
}

// KeepAlive will return whether or not the network utility should be
// left running upon EOF. In the case of TLS, it is dependent upon
// mode.
func (nut *TLSNUt) KeepAlive() bool {
	if nut.mode == modeServer {
		return nut.up
	}

	return false
}

func (nut *TLSNUt) listen(addr string) error {
	var a *net.TCPAddr
	var c net.Conn
	var e error
	var l *net.TCPListener

	if a, e = net.ResolveTCPAddr("tcp", addr); e != nil {
		return errors.Newf("failed to resolve %s: %w", addr, e)
	}

	if l, e = net.ListenTCP("tcp", a); e != nil {
		return errors.Newf("failed to listen on %s: %w", addr, e)
	}

	// Wrap with TLS
	nut.list = tls.NewListener(l, nut.tlscfg)

	go func() {
		//nolint:mnd // 2 goroutines
		var up chan struct{} = make(chan struct{}, 2)
		var wait chan struct{}

		if nut.fork {
			wait = make(chan struct{}, 1)
		} else {
			//nolint:mnd // 2 goroutines
			wait = make(chan struct{}, 2)
		}

		for nut.up {
			if c, e = nut.list.Accept(); e != nil {
				if nut.up {
					e = errors.Newf("connection failed: %w", e)
					logErr(1, "%s", e.Error())
				}

				time.Sleep(time.Millisecond)

				continue
			}

			logGood(1, "Connection from %s", c.RemoteAddr().String())

			go func() {
				up <- struct{}{}

				_, _ = io.Copy(nut.pwIn, c)

				wait <- struct{}{}
			}()

			go func() {
				up <- struct{}{}

				_, _ = io.Copy(c, nut.prOut)

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

func (nut *TLSNUt) parseOpts(k string, v string) error {
	var e error

	switch k {
	case "ca":
		if nut.ca, e = readCert(v); e != nil {
			return e
		}
	case "cert":
		if nut.cert, e = readCert(v); e != nil {
			return e
		}
	case "echo":
		if nut.mode == modeClient {
			return errors.Newf("unknown %s option %s", nut.Type(), k)
		}

		nut.echo = true
	case "fork":
		if nut.mode == modeClient {
			return errors.Newf("unknown %s option %s", nut.Type(), k)
		}

		nut.fork = true
	case "key":
		if nut.key, e = readKey(v); e != nil {
			return e
		}
	case "verify":
		nut.verify = true
	default:
		return errors.Newf("unknown %s option %s", nut.Type(), k)
	}

	return nil
}

// Read will read from the current TLS connection.
//
//nolint:dupl,mnd // TLS is TCP (so yeah), log levels
func (nut *TLSNUt) Read(p []byte) (int, error) {
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

	if n, e = nut.prIn.Read(p); e != nil {
		logSubInfo(2, "%s read: %d bytes", nut.String(), n)

		if !nut.up || (nut.mode == modeServer) {
			e = nil
		}

		return n, e
	}

	logSubInfo(2, "%s read: %d bytes", nut.String(), n)

	if (nut.mode == modeServer) && nut.echo {
		if _, e = nut.Write(p); e != nil {
			return n, e
		}
	}

	return n, nil
}

func (nut *TLSNUt) setupTLSConfig() error {
	var b [][]byte
	var pool *x509.CertPool

	// Create initial config
	nut.tlscfg = &tls.Config{} //nolint:gosec // G402 - not a problem

	switch nut.mode {
	case modeClient:
		nut.tlscfg.InsecureSkipVerify = !nut.verify

		if nut.ca != nil {
			pool = x509.NewCertPool()
			pool.AddCert(nut.ca)

			nut.tlscfg.RootCAs = pool
		}

		if (nut.cert == nil) && (nut.key != nil) {
			return errors.New("no cert provided")
		} else if (nut.cert != nil) && (nut.key == nil) {
			return errors.New("no key provided")
		}
	case modeServer:
		if nut.cert == nil {
			return errors.New("no cert provided")
		}

		if nut.key == nil {
			return errors.New("no key provided")
		}

		// Create cert pool
		if nut.verify {
			// Verify clients, but we'll need a CA to verify against
			if nut.ca == nil {
				return errors.New("no ca provided")
			}

			pool = x509.NewCertPool()
			pool.AddCert(nut.ca)

			nut.tlscfg.ClientAuth = tls.RequireAndVerifyClientCert
			nut.tlscfg.ClientCAs = pool
		}
	}

	// Add the cert to the chain
	if nut.cert != nil {
		b = append(b, nut.cert.Raw)
	}

	// Add the CA to the chain
	if nut.ca != nil {
		b = append(b, nut.ca.Raw)
	}

	// Add TLS cert to config
	if len(b) > 0 {
		nut.tlscfg.Certificates = append(
			nut.tlscfg.Certificates,
			tls.Certificate{
				Certificate: b,
				PrivateKey:  nut.key,
			},
		)
	}

	return nil
}

// Up will start the network utility. In the case of TLS, it will
// either connect or listen, depending on the mode.
func (nut *TLSNUt) Up() error {
	var e error

	nut.lock.Lock()
	defer nut.lock.Unlock()

	// Check if already up
	if nut.up {
		return nil
	}

	// Up after pipes created
	_ = nut.baseNUt.Up()
	nut.connecting = true
	nut.up = true

	// Create connection/listener
	switch nut.mode {
	case modeClient:
		e = nut.connect(nut.addr)
	case modeServer:
		e = nut.listen(nut.addr)
	}

	if e != nil {
		nut.up = false
	}

	return e
}

// Write will write to the current TLS connection.
//
//nolint:mnd // log levels
func (nut *TLSNUt) Write(p []byte) (int, error) {
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
