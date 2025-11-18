//nolint:godoclint // These are tests
package nutsak_test

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	sak "github.com/mjwhitta/nutsak"
	assert "github.com/stretchr/testify/require"
)

func compare(t *testing.T, fn string) {
	t.Helper()

	var b []byte
	var e error
	var expected string = "5c07cdb2b518bb539198b0ffd28ecd23"
	var hash [sha512.Size]byte

	b, e = os.ReadFile(filepath.Clean(fn))
	assert.NoError(t, e)

	hash = sha512.Sum512(b)
	assert.Equal(t, expected, hex.EncodeToString(hash[:])[0:32])
}

func sharedNetworkTests(t *testing.T, fn string, seeds ...string) {
	t.Helper()

	t.Run(
		"NoResolve",
		func(t *testing.T) {
			var a sak.NUt
			var e error

			// Create NUt
			a, e = sak.NewNUt(seeds[0])
			assert.NoError(t, e)

			e = a.Up()
			assert.Error(t, e)

			a, e = sak.NewNUt(seeds[1])
			assert.NoError(t, e)

			e = a.Up()
			assert.Error(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt(seeds[2])
			assert.Error(t, e)
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
			var grErrs []chan error = []chan error{
				make(chan error, 1),
				make(chan error, 1),
			}

			defer func() {
				for i := range grErrs {
					for e := range grErrs[i] {
						assert.NoError(t, e)
					}
				}
			}()

			// Create NUts
			a, e = sak.NewNUt(seeds[3])
			assert.NoError(t, e)

			b, e = sak.NewNUt(seeds[4])
			assert.NoError(t, e)

			c, e = sak.NewNUt(seeds[5])
			assert.NoError(t, e)

			d, e = sak.NewNUt(seeds[6])
			assert.NoError(t, e)

			// Pair NUts
			go func() {
				grErrs[0] <- sak.Pair(c, d)

				close(grErrs[0])
			}()
			go func() {
				grErrs[1] <- sak.Pair(a, b)

				close(grErrs[1])
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.NoError(t, e)

			e = b.Down()
			assert.NoError(t, e)

			e = c.Down()
			assert.NoError(t, e)

			e = d.Down()
			assert.NoError(t, e)

			compare(t, fn)
		},
	)
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
			assert.Error(t, e)
		},
	)

	t.Run(
		"UnknownMode",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("file:testdata/in,mode=asdf")
			assert.Error(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("file:testdata/in,asdf")
			assert.Error(t, e)
		},
	)

	t.Run(
		"SuccessAppend",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var e error
			var grErrs chan error = make(chan error, 1)

			defer func() {
				for e := range grErrs {
					assert.NoError(t, e)
				}
			}()

			// Ensure file doesn't already exist
			_ = os.Remove("testdata/out_file_a")

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.NoError(t, e)

			b, e = sak.NewNUt("file:testdata/out_file_a,mode=append")
			assert.NoError(t, e)

			// Pair NUts
			go func() {
				grErrs <- sak.Pair(a, b)

				close(grErrs)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.NoError(t, e)

			e = b.Down()
			assert.NoError(t, e)

			compare(t, "testdata/out_file_a")
		},
	)

	t.Run(
		"SuccessWrite",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var e error
			var grErrs chan error = make(chan error, 1)

			defer func() {
				for e := range grErrs {
					assert.NoError(t, e)
				}
			}()

			// Create NUts
			a, e = sak.NewNUt("file:testdata/in")
			assert.NoError(t, e)

			b, e = sak.NewNUt("file:testdata/out_file_w,mode=write")
			assert.NoError(t, e)

			// Pair NUts
			go func() {
				grErrs <- sak.Pair(a, b)

				close(grErrs)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Down()
			assert.NoError(t, e)

			e = b.Down()
			assert.NoError(t, e)

			compare(t, "testdata/out_file_w")
		},
	)
}

func TestStdioNUt(t *testing.T) {
	t.Run(
		"InvalidAddr",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("stdio:asdf")
			assert.Error(t, e)
		},
	)

	t.Run(
		"UnknownOption",
		func(t *testing.T) {
			var e error

			// Create NUt
			_, e = sak.NewNUt("stdio:,asdf")
			assert.Error(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var a sak.NUt
			var b sak.NUt
			var e error
			var grErrs chan error = make(chan error, 1)

			defer func() {
				for e := range grErrs {
					assert.NoError(t, e)
				}
			}()

			// Create NUts
			a, e = sak.NewNUt("-")
			assert.NoError(t, e)

			b, e = sak.NewNUt("-")
			assert.NoError(t, e)

			e = a.Open() // For coverage
			assert.NoError(t, e)

			// Pair NUts
			go func() {
				grErrs <- sak.Pair(a, b)

				close(grErrs)
			}()

			// Wait
			time.Sleep(2 * time.Second)

			// Stop NUts
			e = a.Close() // For coverage
			assert.NoError(t, e)

			e = b.Down()
			assert.NoError(t, e)
		},
	)
}

func TestStream(t *testing.T) {
	var a sak.NUt
	var b sak.NUt
	var e error
	var grErrs chan error = make(chan error, 1)

	defer func() {
		for e := range grErrs {
			assert.NoError(t, e)
		}
	}()

	// Create NUts
	a, e = sak.NewNUt("-")
	assert.NoError(t, e)

	b, e = sak.NewNUt("-")
	assert.NoError(t, e)

	// Stream NUts
	go func() {
		grErrs <- sak.Stream(a, b)

		close(grErrs)
	}()

	// Wait
	time.Sleep(2 * time.Second)

	// Stop NUts
	e = a.Down()
	assert.NoError(t, e)

	e = b.Down()
	assert.NoError(t, e)
}

func TestTCPNUt(t *testing.T) {
	sharedNetworkTests(
		t,
		"testdata/out_tcp",
		"tcp:doesnotexist.asdf.com:4444",
		"tcp-l:doesnotexist.asdf.com:4444",
		"tcp:127.13.37.1:4444,asdf",
		"file:testdata/in",
		"tcp:127.13.37.1:4444",
		"tcp-l:127.13.37.1:4444,fork",
		"file:testdata/out_tcp,mode=write",
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
			assert.Error(t, e)
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
			assert.Error(t, e)
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
			assert.Error(t, e)
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
			assert.Error(t, e)
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
			assert.Error(t, e)

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
			assert.Error(t, e)
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
			assert.Error(t, e)

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
			assert.Error(t, e)
		},
	)

	sharedNetworkTests(
		t,
		"testdata/out_tls",
		"tls:doesnotexist.asdf.com:8443",
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
		"tls:127.13.37.1:8443,asdf",
		"file:testdata/in",
		strings.Join(
			[]string{
				"tls:127.13.37.1:8443",
				"ca=testdata/pki/ca/ca.cert.pem",
				"cert=testdata/pki/certs/user.cert.pem",
				"key=testdata/pki/private/user.key.pem",
			},
			",",
		),
		strings.Join(
			[]string{
				"tls-l:127.13.37.1:8443",
				"ca=testdata/pki/ca/ca.cert.pem",
				"cert=testdata/pki/certs/localhost.cert.pem",
				"fork",
				"key=testdata/pki/private/localhost.key.pem",
			},
			",",
		),
		"file:testdata/out_tls,mode=write",
	)
}

func TestUDPNUt(t *testing.T) {
	sharedNetworkTests(
		t,
		"testdata/out_udp",
		"udp:doesnotexist.asdf.com:4444",
		"udp-l:doesnotexist.asdf.com:4444",
		"udp:127.13.37.1:4444,asdf",
		"file:testdata/in",
		"udp:127.13.37.1:5353",
		"udp-l:127.13.37.1:5353",
		"file:testdata/out_udp,mode=write",
	)
}

func TestUnknownSeeds(t *testing.T) {
	var e error

	_, e = sak.NewNUt("asdf:")
	assert.Error(t, e)

	_, e = sak.NewFileNUt("asdf:")
	assert.Error(t, e)

	_, e = sak.NewStdioNUt("asdf:")
	assert.Error(t, e)

	_, e = sak.NewTCPNUt("asdf:")
	assert.Error(t, e)

	_, e = sak.NewTLSNUt("asdf:")
	assert.Error(t, e)

	_, e = sak.NewUDPNUt("asdf:")
	assert.Error(t, e)
}
