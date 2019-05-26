package ovpncfg

import (
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiam/openvpn-config-generator/lib/certtool"
)

var dhParameters = []byte(`-----BEGIN DH PARAMETERS-----
MIGHAoGBAP+ReUM54jfDGAspfG9K0V1VCxc+bZQlNHBg//bgPbIaEQRGRaVbXlLX
1c5PiSwbagn7wt4odyig3om4PfzPlPRsI0jSy8Ml1BmkDLDi045FT4LOCzoBzFeJ
7Qfr2mQG8gRqXAJhsic7uJbpHyXGVnime8+hWVtUOUB501+ktaxjAgEC
-----END DH PARAMETERS-----`)

func TestServerConfig(t *testing.T) {
	config, err := NewServerConfig()
	assert.NoError(t, err)

	config.Embed("ca", []byte(`-----BEGIN CERTIFICATE-----
MIIEDTCCAnWgAwIBAgIBATANBgkqhkiG9w0BAQsFADAeMQ0wCwYDVQQKEwRBQ01F
MQ0wCwYDVQQDEwRBQ01FMB4XDTE5MDUyNjE3MDkxNFoXDTI5MDUyNjE3MDkxNFow
HjENMAsGA1UEChMEQUNNRTENMAsGA1UEAxMEQUNNRTCCAaIwDQYJKoZIhvcNAQEB
BQADggGPADCCAYoCggGBAKGaRb0DOYml4Fnr5YVfuHj6VQkyOwiTNXzVvkIge9FZ
U6OGvBtNrVAzRuW8ZP3Roe8Oq1D1yLxU6KFUBWx1NPWkxfuYEEKJvMCRofpOJdiV
nNKcM0IW2wlb0RD7EewocgUV0iNnxv4R6XH3WCJzaTWJmDX9eiQOdPw+RKwe+GhY
DnPOIkYnbxyGbyOd6KFaGVA1V3ktCZRm3qn301Oq6KzXASKgXlPc50objMAyC01t
IGYDyTlbLlYpCI8WXz4s/PzloZslCSy8AR5Jfs+tzSagp7SaRqN0DMH0x8/UHCx0
Tr3ML8zZBUu/ZHXGpzv3S4UfneUo2vGYLsT9zqiJdDp5rki/ATBYUZO71VxbA/AR
YHj8dh4z5fWjNQZrQzC5hcgKS2ipAk6IrnYk2klr+nhme4se2bA1vGmyFjX20Lse
Z0SpeBkwnOGQwfNUircixFfu3x3luIxM48RGTvM9g1eGdr2KJcZJdKnPQMpFYVKT
jcw61fgNT06kuolHiwM3vwIDAQABo1YwVDASBgNVHRMBAf8ECDAGAQH/AgEAMB0G
A1UdDgQWBBQxPEkIOyNFrEOaxCgJvg5UIeiNqjAfBgNVHSMEGDAWgBQxPEkIOyNF
rEOaxCgJvg5UIeiNqjANBgkqhkiG9w0BAQsFAAOCAYEARnbIs8/42gWpR5EcAeKR
Ux2p6K7avGg6sYPkcOVfbKAXL1th9ReFyLiJd2K8cSAwY2WB1oIlaJpwS4XVL0WG
0zmAMtnYHkuKNsjBv2OYpOGMEfblYP11tX0+A8JH8Bw/Cx9GzZj/hmVipDaIygpV
r6LvumwcgPMQky42yq4qFEUNqkyB/uSE0MBXRE1QO8PfbKbyD+4DT4Za6X6UUWIS
udwpCNRCBSOsEEPUopme/2QwAP2uyIufbtuc9cddxTudSlLpuIepmMD2umb2vnsp
Dn40iM40C83QYVMpm5hmF5/RUEUchw6baJUqJBaEoVNmHp282DfQ8/QdME1u918A
7IFH2Sdby2mLPRd3MXpTtzHLnQFN1XZV1/0tQLM5EAms51/eS1j+en32vJFgSRtr
usFFyXn/9I/PXNwtbL75M65g02LqtSv7omq7sqVuToOdROmLFqlB3JYhiQQ+vZOn
DGIm0cBbNr3/z3XWsfHVyr0rkTJEbzwenCxdnVrY193q
-----END CERTIFICATE-----`))

	config.Embed("cert", []byte(`-----BEGIN CERTIFICATE-----
MIIEJDCCAoygAwIBAgIBAjANBgkqhkiG9w0BAQsFADAeMQ0wCwYDVQQKEwRBQ01F
MQ0wCwYDVQQDEwRBQ01FMB4XDTE5MDUyNjE3MDkxNVoXDTI5MDUyNjE3MDkxNVow
JDENMAsGA1UEChMEQUNNRTETMBEGA1UEAxMKc2VydmVyLnRsZDCCAaIwDQYJKoZI
hvcNAQEBBQADggGPADCCAYoCggGBAOImIzOB3zp/bMmPof4xDFzHO/6ioGN3uNuo
x49K3XTzFbGCYojkqw8EKgsVhmWuTwmQ2CctzWh+u26L3W/nKQB3p9yvMKO6bMJC
XSiNCJ7nXwaf+v44FGUcSnSRRdY7F2N6fr4EdFYkeTJGowQp8N4P3NotUglNV6ir
P5/+pjJq9aUDIDARhmZkuKD24mUE1XhejDp0OHuNk0oJM5Y8BJUfqrGrGMdbIjM6
7p5ZiH1C99NSdaLlY8CO7kK/Sb/uNk3O2+F0ZSlopf00HMfbBA6kNFFoR+Hf9S7K
Hn0MFHVHB01GvJKTe3P8b7Qfllq4V71w8prkkVVw8KzplYxpdoUpmIqI7O5vE4Aq
ng7WWpsh5tcsEysub2vYpC9d3ciICYbmPPZMGIRxAvn10QyDuHObz3YcxknkxVH7
fgI3Uhb+M9bIu6fipS1fys3R83JI1oMXmE2EqWxN9hMgBqusXUGZJ8vsXfl8SfNE
QxJdj1CRQrDwIbYy/gM20IY2k1FO0QIDAQABo2cwZTAOBgNVHQ8BAf8EBAMCBaAw
EwYDVR0lBAwwCgYIKwYBBQUHAwEwHQYDVR0OBBYEFIBUFCyy9TP4fsnDAmwmz8Qg
Kqx6MB8GA1UdIwQYMBaAFDE8SQg7I0WsQ5rEKAm+DlQh6I2qMA0GCSqGSIb3DQEB
CwUAA4IBgQCJOsFJffxfbxP4P8ycL25KpQd+0SDVNJSRrtHp13p0HRnazJfdXpUp
B9XE/eEGbUwxTw3D26Vq9dBofpYX6+ZIEhAvJJ8HJlJIg1oFBw21NkidrFSV9enB
RLlInrjEmIgq3bwjapHcEuuy1JpdV1cLrqfI6Ofd7l6yjeE31z84Y6t9MA+Fz9XN
hDcxGAA704ogGZnMBze4pAlNN0wWTNbguBlQWQDyDbTXFsWD83R+J9ivS/P8+VVv
AQeRzcmvhEbkOg6J8XVDbH495DO2fnGTgyxLYHbb9k07hNukkkziMl2KA+fPf2jd
LZbx0rDUHvVYo4dUE7wSgBEPSizZZq7Mx8Umz+L0OdBv65TibFy0lItf1m5As0GL
j7u+kFUiGTq0Xeq9TFrgpDKGJD98Cnc14OUvFFdFxxmrR5BY934aw2r7M/S3RbOB
6K9nO7qGjTr4exoTlnWMkpFJPOaUxhQRg0Pc5Tpq54IGtO5nrEvQFKCmsHAs4GaD
p/dfzSKVG0Q=
-----END CERTIFICATE-----`))

	config.Embed("key", []byte(`-----BEGIN PRIVATE KEY-----
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
6RXvi12M2Sr66MYNDD/yKP+d
-----END PRIVATE KEY-----`))

	config.Embed("dh", dhParameters)

	config.Embed("tls-crypt", []byte(`-----BEGIN OpenVPN Static key V1-----
e5e4d6af39289d53
171ecc237a8f996a
97743d146661405e
c724d5913c550a0c
30a48e52dfbeceb6
e2e7bd4a8357df78
4609fe35bbe99c32
bdf974952ade8fb9
71c204aaf4f256ba
eeda7aed4822ff98
fd66da2efa9bf8c5
e70996353e0f96a9
c94c9f9afb17637b
283da25cc99b37bf
6f7e15b38aedc3e8
e6adb40fca5c5463
-----END OpenVPN Static key V1-----`))

	buf, err := config.Compile()
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	err = writeConfig(config, "server.conf")
	assert.NoError(t, err)
}

func TestClientConfig(t *testing.T) {
	config, err := NewClientConfig()
	assert.NoError(t, err)

	config.Embed("ca", []byte(`-----BEGIN CERTIFICATE-----
MIIEJDCCAoygAwIBAgIBAjANBgkqhkiG9w0BAQsFADAeMQ0wCwYDVQQKEwRBQ01F
MQ0wCwYDVQQDEwRBQ01FMB4XDTE5MDUyNjE3MDkxNVoXDTI5MDUyNjE3MDkxNVow
JDENMAsGA1UEChMEQUNNRTETMBEGA1UEAxMKc2VydmVyLnRsZDCCAaIwDQYJKoZI
hvcNAQEBBQADggGPADCCAYoCggGBAOImIzOB3zp/bMmPof4xDFzHO/6ioGN3uNuo
x49K3XTzFbGCYojkqw8EKgsVhmWuTwmQ2CctzWh+u26L3W/nKQB3p9yvMKO6bMJC
XSiNCJ7nXwaf+v44FGUcSnSRRdY7F2N6fr4EdFYkeTJGowQp8N4P3NotUglNV6ir
P5/+pjJq9aUDIDARhmZkuKD24mUE1XhejDp0OHuNk0oJM5Y8BJUfqrGrGMdbIjM6
7p5ZiH1C99NSdaLlY8CO7kK/Sb/uNk3O2+F0ZSlopf00HMfbBA6kNFFoR+Hf9S7K
Hn0MFHVHB01GvJKTe3P8b7Qfllq4V71w8prkkVVw8KzplYxpdoUpmIqI7O5vE4Aq
ng7WWpsh5tcsEysub2vYpC9d3ciICYbmPPZMGIRxAvn10QyDuHObz3YcxknkxVH7
fgI3Uhb+M9bIu6fipS1fys3R83JI1oMXmE2EqWxN9hMgBqusXUGZJ8vsXfl8SfNE
QxJdj1CRQrDwIbYy/gM20IY2k1FO0QIDAQABo2cwZTAOBgNVHQ8BAf8EBAMCBaAw
EwYDVR0lBAwwCgYIKwYBBQUHAwEwHQYDVR0OBBYEFIBUFCyy9TP4fsnDAmwmz8Qg
Kqx6MB8GA1UdIwQYMBaAFDE8SQg7I0WsQ5rEKAm+DlQh6I2qMA0GCSqGSIb3DQEB
CwUAA4IBgQCJOsFJffxfbxP4P8ycL25KpQd+0SDVNJSRrtHp13p0HRnazJfdXpUp
B9XE/eEGbUwxTw3D26Vq9dBofpYX6+ZIEhAvJJ8HJlJIg1oFBw21NkidrFSV9enB
RLlInrjEmIgq3bwjapHcEuuy1JpdV1cLrqfI6Ofd7l6yjeE31z84Y6t9MA+Fz9XN
hDcxGAA704ogGZnMBze4pAlNN0wWTNbguBlQWQDyDbTXFsWD83R+J9ivS/P8+VVv
AQeRzcmvhEbkOg6J8XVDbH495DO2fnGTgyxLYHbb9k07hNukkkziMl2KA+fPf2jd
LZbx0rDUHvVYo4dUE7wSgBEPSizZZq7Mx8Umz+L0OdBv65TibFy0lItf1m5As0GL
j7u+kFUiGTq0Xeq9TFrgpDKGJD98Cnc14OUvFFdFxxmrR5BY934aw2r7M/S3RbOB
6K9nO7qGjTr4exoTlnWMkpFJPOaUxhQRg0Pc5Tpq54IGtO5nrEvQFKCmsHAs4GaD
p/dfzSKVG0Q=
-----END CERTIFICATE-----`))

	config.Embed("cert", []byte(`-----BEGIN CERTIFICATE-----
MIIEKDCCApCgAwIBAgIBAzANBgkqhkiG9w0BAQsFADAeMQ0wCwYDVQQKEwRBQ01F
MQ0wCwYDVQQDEwRBQ01FMB4XDTE5MDUyNjE3MDkxNVoXDTI5MDUyNjE3MDkxNVow
KDENMAsGA1UEChMEQUNNRTEXMBUGA1UEAxMObWFjYm9vay1haXItMTEwggGiMA0G
CSqGSIb3DQEBAQUAA4IBjwAwggGKAoIBgQDZQ7Csztso0OgdySkKGAzVamoZtCY8
snuLJrVrMAondxjzTSoe85RsrtZjj7EKYgF0Etz64C3ibnoPie/rEGEVxbWKKgf7
+5qIQndQEh+TI0rvKpCmjxO0EuhrQQDFL0YKIkTCJxIVSQSsU59vtJog/zNCENup
zLJTX6Tl/SbDlUWy2+D7lyjXD0WvcIpnAn6UkFMWrlqFhxkLR6dmI96KkVncQTUf
LZ8RQMW9/pdX0Kjo48j/HIHSTIwOUHoGsaiJdKpD2YUMM/GmKZYJ9oTwF+zGW0pF
SyuuyTzD+slS7J3X9UXcVqJ+a6jm8PfC7YnjGDl2KiahilqbVaP7HSAVSlvhrI7T
MvRlDHa/xbOvHTKOmQSIgK5rfMim98KA/d7w7BLL3PrIR870p2X8JHk2udEzr7i6
srxMPq67Br+mGmeAQHjA+ArqyWnEwSlnW4c1AUVese+syoQ+ULIJkstbBOlU/9Kh
fuXAwlSF5YNsRl3rgqVA1XrNvdOYlg1IWzMCAwEAAaNnMGUwDgYDVR0PAQH/BAQD
AgeAMBMGA1UdJQQMMAoGCCsGAQUFBwMCMB0GA1UdDgQWBBRtgmfCIY++fYdK8fFT
jY8hFlJJhTAfBgNVHSMEGDAWgBQxPEkIOyNFrEOaxCgJvg5UIeiNqjANBgkqhkiG
9w0BAQsFAAOCAYEAHTNKrwNW8kbgSA23S6nFG58dOIchXWGZT2TBlBAHU9/PxMlw
J4l0ThlfKrY29v/gqnh4HuaSbBtQNbh/ZbOgIphhDuRXp/LPorQEcoX/kQiZxK/C
e2Wiz1YXPqktvyotJmUuqNWKwbwQuivLmzWlreUFjsAp0CoGZWm+/8UwF2/OrVS3
tG9eMcqVkKo6JOZQtJ5XS8rpRSpxzS6WBIGexVNOuJPLGYPKN9aIRN5aOQeDrRCu
ahQ+N5ukrw+FTBXduWrxAnUQCZRL+S9pZ3rxRN4KZsXZ0yLtFlj1AACURz7guhyN
dcwezhD7t8jcRDb0hLBhxyzDeea6vEbbJwKewxFdC/+msHDsMnaeFXzbPYXKjuQ8
JP7Me7kOR6a2lia9kvxOswqJzGzQJcQBPaIrlBCJsSJdbNt38aYjsZ88kP39YsTH
vgpnzorDGrglUTsSV6kgqgQhSNurS6yRWWWQFy+npCkUk44/ZRfIYhn8/C33Ndf0
TTo/GlRxg1owOA8e
-----END CERTIFICATE-----`))

	config.Embed("key", []byte(`-----BEGIN PRIVATE KEY-----
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
6RXvi12M2Sr66MYNDD/yKP+d
-----END PRIVATE KEY-----`))

	config.Embed("tls-crypt", []byte(`-----BEGIN OpenVPN Static key V1-----
e5e4d6af39289d53
171ecc237a8f996a
97743d146661405e
c724d5913c550a0c
30a48e52dfbeceb6
e2e7bd4a8357df78
4609fe35bbe99c32
bdf974952ade8fb9
71c204aaf4f256ba
eeda7aed4822ff98
fd66da2efa9bf8c5
e70996353e0f96a9
c94c9f9afb17637b
283da25cc99b37bf
6f7e15b38aedc3e8
e6adb40fca5c5463
-----END OpenVPN Static key V1-----`))

	buf, err := config.Compile()
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	err = writeConfig(config, "client.conf")
	assert.NoError(t, err)
}

func TestGenOpenVPNKey(t *testing.T) {
	buf, err := GenOpenVPNStaticKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
}

func TestGenCertificates(t *testing.T) {
	caCert, caKey, err := certtool.BuildCA()
	assert.NoError(t, err)

	serverCert, serverKey, err := certtool.BuildServerCertificate(caCert, caKey, "server.tld")
	assert.NoError(t, err)

	clientCert, clientKey, err := certtool.BuildClientCertificate(caCert, caKey, "client.local")
	assert.NoError(t, err)

	writeCert(caCert, "ca.crt")
	writeCert(serverCert, "server.crt")
	writeCert(clientCert, "client.crt")

	writeKey(serverKey, "server.key")
	writeKey(clientKey, "client.key")
}

func TestAutomaticConfig(t *testing.T) {
	caCert, caKey, err := certtool.BuildCA()
	assert.NoError(t, err)

	serverCert, serverKey, err := certtool.BuildServerCertificate(caCert, caKey, "server.tld")
	assert.NoError(t, err)

	clientCert, clientKey, err := certtool.BuildClientCertificate(caCert, caKey, "client.local")
	assert.NoError(t, err)

	vpnStaticKey, err := GenOpenVPNStaticKey()
	assert.NoError(t, err)

	{
		config, err := NewServerConfig()
		assert.NoError(t, err)

		config.Embed("ca", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCert}))

		config.Embed("cert", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCert}))
		config.Embed("key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: serverKey}))

		config.Embed("dh", dhParameters)

		config.Embed("tls-crypt", vpnStaticKey)

		err = writeConfig(config, "server-auto.conf")
		assert.NoError(t, err)
	}

	{
		config, err := NewClientConfig()
		assert.NoError(t, err)

		config.Embed("ca", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCert}))

		config.Embed("cert", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCert}))
		config.Embed("key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: clientKey}))

		config.Embed("tls-crypt", vpnStaticKey)

		err = writeConfig(config, "client-auto.conf")
		assert.NoError(t, err)
	}
}
