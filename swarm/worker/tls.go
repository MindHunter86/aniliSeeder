package worker

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

func (m *Worker) getTLSCertificate() (_ tls.Certificate, e error) {
	// TODO
	// if gCli.String("swarm-master-custom-ca") != ""
	// get key, get pub, return

	var cbytes, kbytes []byte
	if cbytes, kbytes, e = m.createPublicPrivatePair(); e != nil {
		return
	}

	return tls.X509KeyPair(cbytes, kbytes)
}

func (*Worker) createPublicPrivatePair() (_, _ []byte, e error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		return
	}

	certBytes, e := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if e != nil {
		return
	}

	var cbuf = new(bytes.Buffer)
	pem.Encode(cbuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var p []byte
	if p, e = x509.MarshalECPrivateKey(priv); e != nil {
		return
	}

	var pbuf = new(bytes.Buffer)
	pem.Encode(pbuf, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: p,
	})

	if gCli.Bool("http-debug") {
		fmt.Println("\n" + cbuf.String())
		fmt.Println("\n" + pbuf.String())
	}

	return cbuf.Bytes(), pbuf.Bytes(), nil
}
