package ovpncfg

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

func env(name string, defaultValue interface{}) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fmt.Sprintf("%v", defaultValue)
}

func writeConfig(config *generator.Config, file string) error {
	buf, err := config.Compile()
	if err != nil {
		return err
	}

	return writeFile(buf, file)
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

	config.MustAdd("push", "redirect-gateway def1 bypass-dhcp")

	config.MustAdd("push", fmt.Sprintf("dhcp-option DNS %s", env("DNS1", defaultDNS1)))
	config.MustAdd("push", fmt.Sprintf("dhcp-option DNS %s", env("DNS2", defaultDNS2)))

	config.MustEnable("client-to-client")
	config.MustSet("keepalive", 10, 120)

	config.MustAdd("push", "ping 15")
	config.MustAdd("push", "ping-restart 60")

	config.MustSet("cipher", "AES-256-GCM")
	config.MustSet("ncp-ciphers", "AES-256-GCM:AES-256-CBC:AES-128-GCM:AES-128-CBC:BF-CBC")

	config.MustEnable("comp-lzo")

	config.MustSet("user", "nobody")
	config.MustSet("group", "nobody")

	config.MustEnable("persist-key")
	config.MustEnable("persist-tun")

	config.MustSet("verb", 5)

	config.MustSet("sndbuf", 0)
	config.MustSet("rcvbuf", 0)

	config.MustAdd("push", "sndbuf 0")
	config.MustAdd("push", "rcvbuf 0")

	config.MustSet("fragment", 0)
	config.MustSet("mssfix", 0)

	config.MustSet("remote-cert-eku", "TLS Web Client Authentication")

	return config, nil
}

func NewClientConfig() (*generator.Config, error) {
	config := generator.New()

	config.MustEnable("client")
	config.MustSet("dev", "tun")
	config.MustSet("proto", env("PROTO", defaultProto))
	config.MustSet("keysize", 256)
	config.MustSet("remote", "192.168.1.87", 1194)
	config.MustSet("resolv-retry", "infinite")

	config.MustSet("cipher", "AES-256-CBC")
	config.MustEnable("nobind")
	config.MustSet("link-mtu", 1550)
	config.MustEnable("persist-key")
	config.MustEnable("persist-tun")
	config.MustEnable("comp-lzo")
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

func writeCert(cert []byte, file string) error {
	return writeFile(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}), file)
}

func writeKey(cert []byte, file string) error {
	return writeFile(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: cert,
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
