package generator

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/xiam/openvpn-config-generator/lib/config"
)

const (
	defaultHost  = "192.168.1.87"
	defaultPort  = 1194
	defaultProto = "udp"

	defaultNetwork     = "10.9.0.0"
	defaultNetworkMask = "255.255.0.0"

	defaultDNS1 = "8.8.8.8"
	defaultDNS2 = "8.8.4.4"

	defaultSendRecvBuffSize = 512000
)

func env(name string, defaultValue interface{}) interface{} {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

func NewServerConfig() (*config.Config, error) {
	c := config.New()

	c.MustSet("port", env("PORT", defaultPort))
	c.MustSet("proto", env("PROTO", defaultProto))
	c.MustSet("dev", "tun")

	c.MustSet("topology", "subnet")

	c.MustSet("server", env("NETWORK", defaultNetwork), env("NETWORK_MASK", defaultNetworkMask))
	c.MustSet("route", env("NETWORK", defaultNetwork), env("NETWORK_MASK", defaultNetworkMask))

	c.MustSet("ifconfig-pool-persist", "ipp.txt")
	c.MustSet("client-config-dir", "ccd")

	c.MustEnable("client-to-client")
	c.MustSet("keepalive", 30, 150)

	c.MustAdd("push", "ping 30")
	c.MustAdd("push", "ping-restart 150")

	c.MustSet("cipher", "AES-128-GCM")
	c.MustSet("ncp-ciphers", "AES-256-GCM:AES-256-CBC:AES-128-GCM:AES-128-CBC")

	c.MustSet("user", "nobody")
	c.MustSet("group", "nobody")

	c.MustEnable("persist-key")
	c.MustEnable("persist-tun")

	c.MustSet("verb", "3")

	c.MustSet("sndbuf", defaultSendRecvBuffSize)
	c.MustSet("rcvbuf", defaultSendRecvBuffSize)

	c.MustAdd("push", fmt.Sprintf("sndbuf %d", defaultSendRecvBuffSize))
	c.MustAdd("push", fmt.Sprintf("rcvbuf %d", defaultSendRecvBuffSize))

	c.MustAdd("txqueuelen", 1000)

	c.MustEnable("fast-io")

	c.MustAdd("tun-mtu", 1470)

	c.MustSet("remote-cert-eku", "TLS Web Client Authentication")

	return c, nil
}

func NewClientConfig() (*config.Config, error) {
	c := config.New()

	c.MustEnable("client")
	c.MustSet("dev", "tun")
	c.MustSet("proto", env("PROTO", defaultProto))
	c.MustSet("remote", env("HOST", defaultHost), env("PORT", defaultPort))
	c.MustSet("resolv-retry", "infinite")

	c.MustEnable("nobind")
	c.MustEnable("persist-key")
	c.MustEnable("persist-tun")
	c.MustSet("verb", "3")

	return c, nil
}

func GenOpenVPNStaticKey() ([]byte, error) {
	key := make([]byte, 256)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	buf := []byte(
		"#\n" +
			"# 2048 bit OpenVPN static key\n" +
			"#\n" +
			"-----BEGIN OpenVPN Static key V1-----\n",
	)

	for i := 0; i < 16; i++ {
		buf = append(buf, []byte(fmt.Sprintf("%x\n", key[i*16:(i+1)*16]))...)
	}

	buf = append(buf, []byte("-----END OpenVPN Static key V1-----")...)

	return buf, nil
}

func WriteConfig(c *config.Config, file string) error {
	buf, err := c.Compile()
	if err != nil {
		return err
	}

	return writeFile(buf, file)
}

func WriteCert(cert []byte, file string) error {
	return writeFile(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}), file)
}

func WriteKey(key []byte, file string) error {
	return writeFile(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: key,
	}), file)
}

func writeFile(buf []byte, file string) error {
	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
