package certtool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCA(t *testing.T) {
	cert, key, err := BuildCA()

	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, key)
}

func TestBuildKeyServer(t *testing.T) {
	caCert, caKey, err := BuildCA()
	assert.NoError(t, err)

	cert, key, err := BuildServerCertificate(caCert, caKey, "server.tld")
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, key)
}

func TestBuildKey(t *testing.T) {
	caCert, caKey, err := BuildCA()
	assert.NoError(t, err)

	cert, key, err := BuildClientCertificate(caCert, caKey, "client.local")
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, key)
}
