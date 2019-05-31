package main

import (
	"encoding/pem"
	"github.com/spf13/cobra"
	ovpncfg "github.com/xiam/openvpn-config-generator"
	"io/ioutil"
	"log"
)

var clientConfigCmd = &cobra.Command{
	Use:   "client-config [OPTIONS]",
	Short: "Create a client.ovpn file for OpenVPN clients",
	Run:   clientConfigFn,
}

func clientConfigFn(cmd *cobra.Command, args []string) {
	caCert, _ := cmd.Flags().GetString("ca-cert")
	cert, _ := cmd.Flags().GetString("cert")
	key, _ := cmd.Flags().GetString("key")
	tlsKey, _ := cmd.Flags().GetString("tls-crypt")
	output, _ := cmd.Flags().GetString("output")

	remote, _ := cmd.Flags().GetString("remote")
	if remote == "" {
		log.Fatal("missing required --remote parameter")
	}

	checkFile(cmd, caCert, "missing CA certificate")
	checkFile(cmd, cert, "missing certificate")
	checkFile(cmd, key, "missing private key")
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

	tlsKeyBytes, err := ioutil.ReadFile(tlsKey)
	if err != nil {
		log.Fatal("failed to load TLS Authentication key: ", err)
	}

	config, err := ovpncfg.NewClientConfig()
	if err != nil {
		log.Fatal("failed to create client config")
	}

	config.MustSet("remote", remote)

	config.MustEmbed("ca", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes}))

	config.MustEmbed("cert", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}))
	config.MustEmbed("key", pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}))

	config.MustEmbed("tls-crypt", tlsKeyBytes)

	err = ovpncfg.WriteConfig(config, output)
	if err != nil {
		log.Fatal("could not write config file: ", err)
	}

	log.Printf(`Your new client configuration file was written to: %q`, output)
}

func init() {
	clientConfigCmd.Flags().StringP("ca-cert", "r", "ca.crt", "CA certificate")
	clientConfigCmd.Flags().StringP("cert", "c", "client.crt", "Certificate")
	clientConfigCmd.Flags().StringP("key", "k", "client.key", "Private key")
	clientConfigCmd.Flags().StringP("tls-crypt", "t", "key.tlsauth", "TLS Authentication key")
	clientConfigCmd.Flags().String("remote", "", "Address of the remote OpenVPN server")
	clientConfigCmd.Flags().StringP("output", "o", "client.ovpn", "Output file")
}
