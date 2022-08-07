package nutsak_test

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
	sak "gitlab.com/mjwhitta/nutsak"
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
	sak.Banner()
	sak.BannerNSFW()
}

func TestFileNUt(t *testing.T) {
	var a sak.NUt
	var b sak.NUt
	var e error

	// Create NUts
	a, e = sak.NewNUt("file:testdata/in")
	assert.Nil(t, e)

	b, e = sak.NewNUt("file:testdata/out_file,mode=write")
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

	compare(t, "testdata/out_file")
}

func TestStdioNUt(t *testing.T) {
	var a sak.NUt
	var b sak.NUt
	var e error

	// Create NUts
	a, e = sak.NewNUt("-")
	assert.Nil(t, e)

	b, e = sak.NewNUt("-")
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
}

func TestTCPNUt(t *testing.T) {
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
}

func TestTLSNUt(t *testing.T) {
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
}

func TestUDPNUt(t *testing.T) {
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
}
