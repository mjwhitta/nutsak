package nutsak_test

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
	"testing"
	"time"

	sak "github.com/mjwhitta/nutsak"
	assert "github.com/stretchr/testify/require"
)

func compare(t *testing.T, fn string) {
	var b []byte
	var e error
	var expected string = "12daa5c0e1cddb3200bad9190aa2c0fa"
	var hash [md5.Size]byte

	b, e = os.ReadFile(fn)
	assert.Nil(t, e)

	hash = md5.Sum(b)
	assert.Equal(t, expected, hex.EncodeToString(hash[:]))
}

func TestBanner(t *testing.T) {
	sak.Banner(true)
	sak.BannerNSFW(true)
}

func TestFileNUt(t *testing.T) {
	t.Run(
		"NoName",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("file:")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownMode",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("file:testdata/in,mode=asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("file:testdata/in,asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownType",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewFileNUt("asdf:")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"SuccessAppend",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var e error

			// Ensure file doesn't already exist
			os.Remove("testdata/out_file_a")

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.Nil(t, e)

			b, e = sak.NewNUt("file:testdata/out_file_a,mode=append")
			assert.Nil(t, e)

			// Pair NUts
			go func() {
				e = sak.Pair(a, b)
				assert.Nil(t, e)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.Nil(t, e)

			e = b.Down()
			assert.Nil(t, e)

			compare(t, "testdata/out_file_a")
		},
	)

	t.Run(
		"SuccessWrite",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var e error

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.Nil(t, e)

			b, e = sak.NewNUt("file:testdata/out_file_w,mode=write")
			assert.Nil(t, e)

			// Pair NUts
			go func() {
				e = sak.Pair(a, b)
				assert.Nil(t, e)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.Nil(t, e)

			e = b.Down()
			assert.Nil(t, e)

			compare(t, "testdata/out_file_w")
		},
	)
}

func TestNewNUt(t *testing.T) {
	var e error

	// Create NUt
	_, e = sak.NewNUt("asdf:")
	assert.NotNil(t, e)
}

func TestStdioNUt(t *testing.T) {
	t.Run(
		"InvalidAddr",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("stdio:asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("stdio:,asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownType",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewStdioNUt("asdf:")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var e error

			// Create NUts
			a, e = sak.NewNUt("-")
			assert.Nil(t, e)

			b, e = sak.NewNUt("-")
			assert.Nil(t, e)

			e = a.Open() // For coverage
			assert.Nil(t, e)

			// Pair NUts
			go func() {
				e = sak.Pair(a, b)
				assert.Nil(t, e)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Close() // For coverage
			assert.Nil(t, e)

			e = b.Down()
			assert.Nil(t, e)
		},
	)
}

func TestStream(t *testing.T) {
	var a sak.NUt
	var b sak.NUt
	var e error

	// Create NUts
	a, e = sak.NewNUt("-")
	assert.Nil(t, e)

	b, e = sak.NewNUt("-")
	assert.Nil(t, e)

	// Stream NUts
	go func() {
		e = sak.Stream(a, b)
		assert.Nil(t, e)
	}()

	// Wait
	time.Sleep(2 * time.Second)

	// Stop NUts
	e = a.Down()
	assert.Nil(t, e)

	e = b.Down()
	assert.Nil(t, e)
}

func TestTCPNUt(t *testing.T) {
	t.Run(
		"NoResolve",
		func(t *testing.T) {
			var a sak.NUt
			var e error

			// Create NUt
			a, e = sak.NewNUt("tcp:doesnotexist.asdf.com:4444")
			assert.Nil(t, e)

			e = a.Up()
			assert.NotNil(t, e)

			a, e = sak.NewNUt("tcp-l:doesnotexist.asdf.com:4444")
			assert.Nil(t, e)

			e = a.Up()
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("tcp:127.13.37.1:4444,asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownType",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewTCPNUt("asdf:")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var c sak.NUt
			var d sak.NUt
			var e error

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.Nil(t, e)

			b, e = sak.NewNUt("tcp:127.13.37.1:4444")
			assert.Nil(t, e)

			c, e = sak.NewNUt("tcp-l:127.13.37.1:4444,fork")
			assert.Nil(t, e)

			d, e = sak.NewNUt("file:testdata/out_tcp,mode=write")
			assert.Nil(t, e)

			// Pair NUts
			go func() {
				e = sak.Pair(c, d)
				assert.Nil(t, e)
			}()
			go func() {
				e = sak.Pair(a, b)
				assert.Nil(t, e)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.Nil(t, e)

			e = b.Down()
			assert.Nil(t, e)

			e = c.Down()
			assert.Nil(t, e)

			e = d.Down()
			assert.Nil(t, e)

			compare(t, "testdata/out_tcp")
		},
	)
}

func TestTLSNUt(t *testing.T) {
	t.Run(
		"InvalidCA",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls:127.13.37.1:8443",
						"ca=/noexist",
						"cert=/noexist",
						"key=/noexist",
					},
					",",
				),
			)
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"InvalidCert",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls:127.13.37.1:8443",
						"cert=/noexist",
						"key=/noexist",
					},
					",",
				),
			)
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"InvalidKey",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls:127.13.37.1:8443",
						"key=/noexist",
					},
					",",
				),
			)
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"MissingCA",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls-l:127.13.37.1:8443",
						"cert=testdata/pki/certs/localhost.cert.pem",
						"fork",
						"key=testdata/pki/private/localhost.key.pem",
						"verify",
					},
					",",
				),
			)
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"MissingCert",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls:127.13.37.1:8443",
						"key=testdata/pki/private/localhost.key.pem",
					},
					",",
				),
			)
			assert.NotNil(t, e)

			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls-l:127.13.37.1:8443",
						"ca=testdata/pki/ca/ca.cert.pem",
						"fork",
						"key=testdata/pki/private/localhost.key.pem",
						"verify",
					},
					",",
				),
			)
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"MissingKey",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls:127.13.37.1:8443",
						"cert=testdata/pki/certs/localhost.cert.pem",
					},
					",",
				),
			)
			assert.NotNil(t, e)

			_, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls-l:127.13.37.1:8443",
						"ca=testdata/pki/ca/ca.cert.pem",
						"cert=testdata/pki/certs/localhost.cert.pem",
						"fork",
						"verify",
					},
					",",
				),
			)
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"NoResolve",
		func(t *testing.T) {
			var a sak.NUt
			var e error

			// Create NUt
			a, e = sak.NewNUt("tls:doesnotexist.asdf.com:8443")
			assert.Nil(t, e)

			e = a.Up()
			assert.NotNil(t, e)

			a, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls-l:doesnotexist.asdf.com:8443",
						"ca=testdata/pki/ca/ca.cert.pem",
						"cert=testdata/pki/certs/localhost.cert.pem",
						"fork",
						"key=testdata/pki/private/localhost.key.pem",
						"verify",
					},
					",",
				),
			)
			assert.Nil(t, e)

			e = a.Up()
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("tls:127.13.37.1:8443,asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownType",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewTLSNUt("asdf:")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var c sak.NUt
			var d sak.NUt
			var e error

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.Nil(t, e)

			b, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls:127.13.37.1:8443",
						"ca=testdata/pki/ca/ca.cert.pem",
						"cert=testdata/pki/certs/user.cert.pem",
						"key=testdata/pki/private/user.key.pem",
						"verify",
					},
					",",
				),
			)
			assert.Nil(t, e)

			c, e = sak.NewNUt(
				strings.Join(
					[]string{
						"tls-l:127.13.37.1:8443",
						"ca=testdata/pki/ca/ca.cert.pem",
						"cert=testdata/pki/certs/localhost.cert.pem",
						"fork",
						"key=testdata/pki/private/localhost.key.pem",
						"verify",
					},
					",",
				),
			)
			assert.Nil(t, e)

			d, e = sak.NewNUt("file:testdata/out_tls,mode=write")
			assert.Nil(t, e)

			// Pair NUts
			go func() {
				e = sak.Pair(c, d)
				assert.Nil(t, e)
			}()
			go func() {
				e = sak.Pair(a, b)
				assert.Nil(t, e)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.Nil(t, e)

			e = b.Down()
			assert.Nil(t, e)

			e = c.Down()
			assert.Nil(t, e)

			e = d.Down()
			assert.Nil(t, e)

			compare(t, "testdata/out_tls")
		},
	)
}

func TestUDPNUt(t *testing.T) {
	t.Run(
		"NoResolve",
		func(t *testing.T) {
			var a sak.NUt
			var e error

			// Create NUt
			a, e = sak.NewNUt("udp:doesnotexist.asdf.com:4444")
			assert.Nil(t, e)

			e = a.Up()
			assert.NotNil(t, e)

			a, e = sak.NewNUt("udp-l:doesnotexist.asdf.com:4444")
			assert.Nil(t, e)

			e = a.Up()
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("udp:127.13.37.1:4444,asdf")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"UnknownType",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewUDPNUt("asdf:")
			assert.NotNil(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var c sak.NUt
			var d sak.NUt
			var e error

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.Nil(t, e)

			b, e = sak.NewNUt("udp:127.13.37.1:5353")
			assert.Nil(t, e)

			c, e = sak.NewNUt("udp-l:127.13.37.1:5353")
			assert.Nil(t, e)

			d, e = sak.NewNUt("file:testdata/out_udp,mode=write")
			assert.Nil(t, e)

			// Pair NUts
			go func() {
				e = sak.Pair(c, d)
				assert.Nil(t, e)
			}()
			go func() {
				e = sak.Pair(a, b)
				assert.Nil(t, e)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.Nil(t, e)

			e = b.Down()
			assert.Nil(t, e)

			e = c.Down()
			assert.Nil(t, e)

			e = d.Down()
			assert.Nil(t, e)

			compare(t, "testdata/out_udp")
		},
	)
}
