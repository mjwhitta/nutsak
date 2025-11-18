package nutsak

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mjwhitta/errors"
)

func decodeCert(b []byte) (*x509.Certificate, error) {
	var block *pem.Block
	var cert *x509.Certificate
	var e error

	if block, _ = pem.Decode(b); block != nil {
		b = block.Bytes
	}

	if cert, e = x509.ParseCertificate(b); e != nil {
		return nil, errors.Newf("failed to parse cert: %w", e)
	}

	return cert, nil
}

func decodeKey(b []byte) (*rsa.PrivateKey, error) {
	var block *pem.Block
	var e error
	var key *rsa.PrivateKey

	if block, _ = pem.Decode(b); block != nil {
		b = block.Bytes
	}

	if key, e = x509.ParsePKCS1PrivateKey(b); e != nil {
		return nil, errors.Newf("failed to parse key: %w", e)
	}

	return key, nil
}

//nolint:unparam // It might change later
func logErr(lvl int, msg string, args ...any) {
	if (Logger == nil) || (LogLvl < lvl) {
		return
	}

	_ = Logger.Errf(msg, args...)
}

//nolint:unparam // It might change later
func logGood(lvl int, msg string, args ...any) {
	if (Logger == nil) || (LogLvl < lvl) {
		return
	}

	_ = Logger.Goodf(msg, args...)
}

//nolint:unparam // It might change later
func logSubInfo(lvl int, msg string, args ...any) {
	if (Logger == nil) || (LogLvl < lvl) {
		return
	}

	_ = Logger.SubInfof(msg, args...)
}

func readCert(fn string) (*x509.Certificate, error) {
	var b []byte
	var e error

	if b, e = hex.DecodeString(fn); e != nil {
		if b, e = os.ReadFile(filepath.Clean(fn)); e != nil {
			return nil, errors.Newf("failed to read %s: %w", fn, e)
		}
	}

	return decodeCert(b)
}

func readKey(fn string) (*rsa.PrivateKey, error) {
	var b []byte
	var e error

	if b, e = hex.DecodeString(fn); e != nil {
		if b, e = os.ReadFile(filepath.Clean(fn)); e != nil {
			return nil, errors.Newf("failed to read %s: %w", fn, e)
		}
	}

	return decodeKey(b)
}

func stream(a NUt, b NUt) {
	var e error

	// Let things settle
	for !a.IsUp() || !b.IsUp() {
		time.Sleep(time.Millisecond)
	}

	time.Sleep(time.Millisecond)

	for {
		if _, e = io.Copy(b, a); !a.KeepAlive() {
			return
		}

		if e != nil {
			e = errors.Newf("failed to connect %s to %s: %w", a, b, e)
			logErr(1, "%s", e.Error())
		}

		time.Sleep(time.Second)
	}
}
