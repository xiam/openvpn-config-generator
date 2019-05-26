package certtool

import (
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"encoding/asn1"

	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func generateKey(lenght int) (crypto.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, lenght)
}

func buildCert(tpl *x509.Certificate, parent *x509.Certificate, parentKey crypto.PrivateKey) ([]byte, []byte, error) {
	priv, err := generateKey(3072)
	if err != nil {
		return nil, nil, err
	}
	pub := priv.(crypto.Signer).Public()

	if parentKey == nil {
		parentKey = priv
	}

	if tpl.SubjectKeyId == nil {
		spkiASN1, err := x509.MarshalPKIXPublicKey(pub)
		if err != nil {
			return nil, nil, err
		}

		var spki struct {
			Algorithm        pkix.AlgorithmIdentifier
			SubjectPublicKey asn1.BitString
		}
		_, err = asn1.Unmarshal(spkiASN1, &spki)
		if err != nil {
			return nil, nil, err
		}

		skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
		tpl.SubjectKeyId = skid[:]
	}

	if parent.SubjectKeyId != nil {
		tpl.AuthorityKeyId = parent.SubjectKeyId
	}

	cert, err := x509.CreateCertificate(rand.Reader, tpl, parent, pub, parentKey)
	if err != nil {
		return nil, nil, err
	}

	der, err := x509.MarshalPKCS8PrivateKey(priv.(*rsa.PrivateKey))
	if err != nil {
		return nil, nil, err
	}

	return cert, der, nil
}

func BuildCA() (cert []byte, key []byte, err error) {
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"ACME"},
			CommonName:   "ACME",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		MaxPathLenZero:        true,
		BasicConstraintsValid: true,
	}

	return buildCert(tpl, tpl, nil)
}

func BuildServerCertificate(caCert []byte, caKey []byte, commonName string) (cert []byte, key []byte, err error) {
	ca, err := x509.ParseCertificate(caCert)
	if err != nil {
		return nil, nil, err
	}

	pkcsPrivKey, err := x509.ParsePKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, err
	}

	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"ACME"},
			CommonName:   commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: false,
	}

	return buildCert(tpl, ca, pkcsPrivKey)
}

func BuildClientCertificate(caCert []byte, caKey []byte, commonName string) (cert []byte, key []byte, err error) {
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Organization: []string{"ACME"},
			CommonName:   commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: false,
	}

	ca, err := x509.ParseCertificate(caCert)
	if err != nil {
		return nil, nil, err
	}

	pkcsPrivKey, err := x509.ParsePKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, err
	}

	return buildCert(tpl, ca, pkcsPrivKey)
}
