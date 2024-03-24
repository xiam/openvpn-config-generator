package generator

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/xiam/openvpn-config-generator/lib/generator"
)

const (
	defaultPort  = 1194
	defaultProto = "udp"

	defaultNetwork     = "10.9.0.0"
	defaultNetworkMask = "255.255.0.0"

	defaultDNS1 = "8.8.8.8"
	defaultDNS2 = "8.8.4.4"
)

func env(name string, defaultValue interface{}) interface{} {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

func NewServerConfig() (*generator.Config, error) {
	config := generator.New()

	config.MustSet("port", env("PORT", defaultPort))
	config.MustSet("proto", env("PROTO", defaultProto))
	config.MustSet("dev", "tun")

	config.MustSet("topology", "subnet")

	config.MustSet("server", env("NETWORK", defaultNetwork), env("NETWORK_MASK", defaultNetworkMask))
	config.MustSet("route", env("NETWORK", defaultNetwork), env("NETWORK_MASK", defaultNetworkMask))

	config.MustSet("ifconfig-pool-persist", "ipp.txt")
	config.MustSet("client-config-dir", "ccd")

	//config.MustAdd("push", "redirect-gateway def1 bypass-dhcp")

	config.MustEnable("client-to-client")
	config.MustSet("keepalive", 30, 150)

	config.MustAdd("push", "ping 30")
	config.MustAdd("push", "ping-restart 150")

	config.MustSet("cipher", "AES-128-GCM")
	config.MustSet("ncp-ciphers", "AES-256-GCM:AES-256-CBC:AES-128-GCM:AES-128-CBC")

	config.MustSet("user", "nobody")
	config.MustSet("group", "nobody")

	config.MustEnable("persist-key")
	config.MustEnable("persist-tun")

	config.MustSet("verb", "3")

	config.MustSet("sndbuf", 512000)
	config.MustSet("rcvbuf", 512000)

	config.MustAdd("push", "sndbuf 512000")
	config.MustAdd("push", "rcvbuf 512000")

	config.MustAdd("txqueuelen", 1000)

	config.MustAdd("fast-io")

	config.MustAdd("tun-mtu", 1470)

	config.MustSet("remote-cert-eku", "TLS Web Client Authentication")

	return config, nil
}

func NewClientConfig() (*generator.Config, error) {
	config := generator.New()

	config.MustEnable("client")
	config.MustSet("dev", "tun")
	config.MustSet("proto", env("PROTO", defaultProto))
	config.MustSet("remote", "192.168.1.87", 1194)
	config.MustSet("resolv-retry", "infinite")

	config.MustEnable("nobind")
	config.MustEnable("persist-key")
	config.MustEnable("persist-tun")
	config.MustSet("verb", "3")

	return config, nil
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

func WriteConfig(config *generator.Config, file string) error {
	buf, err := config.Compile()
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
