package certtool

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"os"
	"time"
)

func env(name string, defaultValue string) string {
	if s := os.Getenv(name); s != "" {
		return s
	}
	return defaultValue
}

func pkixNameFromEnv() pkix.Name {
	return pkix.Name{
		Organization: []string{env("KEY_ORG", "ACME Corporation")},
		CommonName:   env("KEY_CN", "ACME Certificate"),
		Country:      []string{env("KEY_COUNTRY", "Unknown Country")},
		Locality:     []string{env("KEY_LOCALITY", "Unknown Locality")},
	}
}

func randomSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	return serialNumber, nil
}

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

// BuildCA creates a self-signed CA certificate.
func BuildCA() (cert []byte, key []byte, err error) {
	now := time.Now()

	tpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkixNameFromEnv(),
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		IsCA:                  true,
		MaxPathLenZero:        true,
		BasicConstraintsValid: true,
	}

	return buildCert(tpl, tpl, nil)
}

// BuildServerCertificate creates a certificate that can be used
// for server authentication.
func BuildServerCertificate(caCert []byte, caKey []byte, commonName string) (cert []byte, key []byte, err error) {
	ca, err := x509.ParseCertificate(caCert)
	if err != nil {
		return nil, nil, err
	}

	pkcsPrivKey, err := x509.ParsePKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := randomSerialNumber()
	if err != nil {
		return nil, nil, err
	}

	subject := pkixNameFromEnv()
	subject.CommonName = commonName

	now := time.Now()
	tpl := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: false,
	}

	return buildCert(tpl, ca, pkcsPrivKey)
}

// BuildClientCertificate creates a certificate that can be used
// for client authentication.
func BuildClientCertificate(caCert []byte, caKey []byte, commonName string) (cert []byte, key []byte, err error) {
	ca, err := x509.ParseCertificate(caCert)
	if err != nil {
		return nil, nil, err
	}

	pkcsPrivKey, err := x509.ParsePKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := randomSerialNumber()
	if err != nil {
		return nil, nil, err
	}

	subject := pkixNameFromEnv()
	subject.CommonName = commonName

	now := time.Now()
	tpl := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: false,
	}

	return buildCert(tpl, ca, pkcsPrivKey)
}
