package main

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	generator "github.com/xiam/openvpn-config-generator"
)

var serverConfigCmd = &cobra.Command{
	Use:   "server-config [OPTIONS]",
	Short: "Create a server.conf file for OpenVPN",
	Run:   serverConfigFn,
}

func serverConfigFn(cmd *cobra.Command, args []string) {
	caCert, _ := cmd.Flags().GetString("ca")
	cert, _ := cmd.Flags().GetString("cert")
	key, _ := cmd.Flags().GetString("key")
	dhKey, _ := cmd.Flags().GetString("dh")
	tlsKey, _ := cmd.Flags().GetString("tls-crypt")
	output, _ := cmd.Flags().GetString("output")

	network, _ := cmd.Flags().GetString("network")
	netmask, _ := cmd.Flags().GetString("netmask")

	dns1, _ := cmd.Flags().GetString("dns1")
	dns2, _ := cmd.Flags().GetString("dns2")

	checkFile(cmd, caCert, "missing CA certificate")
	checkFile(cmd, cert, "missing certificate")
	checkFile(cmd, key, "missing private key")
	checkFile(cmd, dhKey, "missing Diffie-Hellman key")
	checkFile(cmd, tlsKey, "missing TLS Authentication Key")

	caCertBytes, err := readPemFile(caCert)
	if err != nil {
		log.Fatal("failed to parse CA certificate: ", err)
	}

	certBytes, err := readPemFile(cert)
	if err != nil {
		log.Fatal("failed to parse server certificate: ", err)
	}

	keyBytes, err := readPemFile(key)
	if err != nil {
		log.Fatal("failed to parse server private key: ", err)
	}

	dhKeyBytes, err := ioutil.ReadFile(dhKey)
	if err != nil {
		log.Fatal("failed to load Diffie-Hellman exchange key: ", err)
	}

	tlsKeyBytes, err := ioutil.ReadFile(tlsKey)
	if err != nil {
		log.Fatal("failed to load TLS Authentication key: ", err)
	}

	config, err := generator.NewServerConfig()
	if err != nil {
		log.Fatal("failed to create server config")
	}

	serverValue := []interface{}{network, netmask}
	config.MustSet("server", serverValue...)
	config.MustSet("route", serverValue...)

	config.MustAdd("push", fmt.Sprintf("dhcp-option DNS %s", dns1))
	config.MustAdd("push", fmt.Sprintf("dhcp-option DNS %s", dns2))

	config.MustEmbed("ca", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes}))

	config.MustEmbed("cert", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}))
	config.MustEmbed("key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}))

	config.MustEmbed("dh", dhKeyBytes)

	config.MustEmbed("tls-crypt", tlsKeyBytes)

	err = generator.WriteConfig(config, output)
	if err != nil {
		log.Fatal("could not write config file: ", err)
	}

	log.Printf(`Your new server configuration file was written to: %q`, output)
}

func init() {
	serverConfigCmd.Flags().StringP("ca", "r", "ca.crt", "CA certificate")
	serverConfigCmd.Flags().StringP("cert", "c", "server.crt", "Certificate")
	serverConfigCmd.Flags().StringP("key", "k", "server.key", "Private key")
	serverConfigCmd.Flags().StringP("dh", "d", "dh.pem", "Diffie-Helman key exchange file")
	serverConfigCmd.Flags().StringP("tls-crypt", "t", "key.tlsauth", "TLS Authentication key")
	serverConfigCmd.Flags().String("network", "10.9.0.0", "Network")
	serverConfigCmd.Flags().String("netmask", "255.255.0.0", "Netmask")
	serverConfigCmd.Flags().String("dns1", "8.8.8.8", "DNS1")
	serverConfigCmd.Flags().String("dns2", "8.8.4.4", "DNS2")
	serverConfigCmd.Flags().StringP("output", "o", "server.conf", "Output file")
}
