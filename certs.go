package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/snwfdhmp/errlog"
	"math/big"
	"os"
	"time"
)

const PEM_CERT = "./localhost+1.pem"
const PEM_KEY = "localhost+1-key.pem"

func (s *server) Certificates() {
	if !(fileExists(PEM_CERT) && fileExists(PEM_KEY)) {
		s.Log().Info("Self-signed Cert not found. Creating...")
		if err := generateSelfSignedCerts(s.Config().Domain, s); errlog.Debug(err) {
			s.Log().Error(err)
			panic(err)
		}
		s.Log().Info("Created self-signed certificates")
	} else {
		s.Log().Info("Self-signed certificates found.")
	}
}

// https://gist.github.com/samuel/8b500ddd3f6118d052b5e6bc16bc4c09
func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func generateSelfSignedCerts(domain string, s *server) error {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if errlog.Debug(err) {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Zendesk Product Security Challenge March 2020"},
			CommonName:   "localhost",
		},
		DNSNames:              []string{"localhost", "127.0.0.1"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if errlog.Debug(err) {
		return err
	}

	out := &bytes.Buffer{}
	if err := pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); errlog.Debug(err) {
		return err
	}
	f, err := os.Create(PEM_CERT)
	if errlog.Debug(err) {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(out.String()); errlog.Debug(err) {
		return err
	}

	g, err := os.Create(PEM_KEY)
	if errlog.Debug(err) {
		return err
	}

	defer g.Close()
	out.Reset()
	if err := pem.Encode(out, pemBlockForKey(priv)); errlog.Debug(err) {
		return err
	}
	if _, err := g.WriteString(out.String()); errlog.Debug(err) {
		return err
	}

	return nil
}

// https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
