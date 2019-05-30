package main

import (
	"encoding/pem"
	"github.com/spf13/cobra"
	ovpncfg "github.com/xiam/openvpn-config-generator"
	"io/ioutil"
	"log"
)

var serverConfigCmd = &cobra.Command{
	Use:   "server-config [OPTIONS]",
	Short: "Create a server.conf file for OpenVPN",
	Run:   serverConfigFn,
}

func serverConfigFn(cmd *cobra.Command, args []string) {
	caCert, _ := cmd.Flags().GetString("ca-cert")
	cert, _ := cmd.Flags().GetString("cert")
	key, _ := cmd.Flags().GetString("key")
	dhKey, _ := cmd.Flags().GetString("dh")
	tlsKey, _ := cmd.Flags().GetString("tls-crypt")
	output, _ := cmd.Flags().GetString("output")

	checkFile(cmd, caCert, "missing CA certificate")
	checkFile(cmd, cert, "missing certificate")
	checkFile(cmd, key, "missing private key")
	checkFile(cmd, dhKey, "missing Diffie-Hellman key")
	checkFile(cmd, dhKey, "missing TLS Authentication Key")

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

	config, err := ovpncfg.NewServerConfig()
	if err != nil {
		log.Fatal("failed to create server config")
	}

	config.MustEmbed("ca", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes}))

	config.MustEmbed("cert", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}))
	config.MustEmbed("key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}))

	config.MustEmbed("dh", dhKeyBytes)

	config.MustEmbed("tls-crypt", tlsKeyBytes)

	err = ovpncfg.WriteConfig(config, output)
	if err != nil {
		log.Fatal("could not write config file: ", err)
	}
}

func init() {
	serverConfigCmd.Flags().StringP("ca-cert", "r", "ca.crt", "CA certificate")
	serverConfigCmd.Flags().StringP("cert", "c", "server.crt", "Certificate")
	serverConfigCmd.Flags().StringP("key", "k", "server.key", "Private key")
	serverConfigCmd.Flags().StringP("dh", "d", "dh.pem", "Diffie-Helman key exchange file")
	serverConfigCmd.Flags().StringP("tls-crypt", "t", "key.tlsauth", "TLS Authentication key")
	serverConfigCmd.Flags().StringP("output", "o", "server.conf", "Output file")
}
