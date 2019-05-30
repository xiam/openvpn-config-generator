package main

import (
	"crypto/x509"
	"fmt"
	"github.com/spf13/cobra"
	ovpncfg "github.com/xiam/openvpn-config-generator"
	"github.com/xiam/openvpn-config-generator/lib/certtool"
	"log"
	"os"
	"path"
)

var buildKeyCmd = &cobra.Command{
	Use:   "build-key [OPTIONS]",
	Short: "Create and sign a client certificate",
	Run:   buildKeyFn,
}

func buildKeyFn(cmd *cobra.Command, args []string) {
	caCertFile, _ := cmd.Flags().GetString("ca-cert")
	caCertBytes, err := readPemFile(caCertFile)
	if err != nil {
		log.Fatal("failed to read certificate: ", err)
	}

	_, err = x509.ParseCertificate(caCertBytes)
	if err != nil {
		log.Fatal("failed to parse certificate: ", err)
	}

	caCertKey, _ := cmd.Flags().GetString("ca-key")
	caKeyBytes, err := readPemFile(caCertKey)
	if err != nil {
		log.Fatal("failed to read private key: ", err)
	}

	_, err = x509.ParsePKCS8PrivateKey(caKeyBytes)
	if err != nil {
		log.Fatal("failed to parse private key: ", err)
	}

	commonname, _ := cmd.Flags().GetString("commonname")
	basename := path.Base(commonname)

	clientCert, clientKey, err := certtool.BuildClientCertificate(caCertBytes, caKeyBytes, commonname)
	if err != nil {
		log.Fatal("failed to build server certificate: ", err)
	}

	basedir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to retrieve base dir: ", err)
	}

	certFile := path.Join(basedir, fmt.Sprintf("%s.crt", basename))
	if err := ovpncfg.WriteCert(clientCert, certFile); err != nil {
		log.Fatal("failed to write certificate: ", err)
	}

	keyFile := path.Join(basedir, fmt.Sprintf("%s.key", basename))
	if err := ovpncfg.WriteKey(clientKey, keyFile); err != nil {
		log.Fatal("failed to write key: ", err)
	}

	log.Printf(`Your new client certificate was successfully generated.`)
	log.Printf(`certificate: %q`, certFile)
	log.Printf(`private key: %q`, keyFile)
}

func init() {
	buildKeyCmd.Flags().StringP("commonname", "n", "client", "name of the client")
	buildKeyCmd.Flags().StringP("ca-cert", "c", "ca.crt", "CA certificate path")
	buildKeyCmd.Flags().StringP("ca-key", "k", "ca.key", "CA private key path")
}
