package nutsak

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"

	"gitlab.com/mjwhitta/errors"
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

func readCert(fn string) (*x509.Certificate, error) {
	var b []byte
	var e error

	if b, e = hex.DecodeString(fn); e != nil {
		if b, e = os.ReadFile(fn); e != nil {
			return nil, errors.Newf("failed to read %s: %w", fn, e)
		}
	}

	return decodeCert(b)
}

func readKey(fn string) (*rsa.PrivateKey, error) {
	var b []byte
	var e error

	if b, e = hex.DecodeString(fn); e != nil {
		if b, e = os.ReadFile(fn); e != nil {
			return nil, errors.Newf("failed to read %s: %w", fn, e)
		}
	}

	return decodeKey(b)
}
