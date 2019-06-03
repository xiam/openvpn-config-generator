package main

import (
	"crypto/x509"
	"fmt"
	"github.com/spf13/cobra"
	ovpncfg "github.com/xiam/openvpn-config-generator"
	"github.com/xiam/openvpn-config-generator/lib/certtool"
	"log"
	"path"
)

var buildKeyCmd = &cobra.Command{
	Use:   "build-key [OPTIONS]",
	Short: "Create and sign a client certificate",
	Run:   buildKeyFn,
}

func buildKeyFn(cmd *cobra.Command, args []string) {
	caCertFile, _ := cmd.Flags().GetString("cert")
	caCertBytes, err := readPemFile(caCertFile)
	if err != nil {
		cmd.Help()
		fmt.Println("")
		log.Fatal("failed to read certificate: ", err)
	}

	_, err = x509.ParseCertificate(caCertBytes)
	if err != nil {
		log.Fatal("failed to parse certificate: ", err)
	}

	caCertKey, _ := cmd.Flags().GetString("key")
	caKeyBytes, err := readPemFile(caCertKey)
	if err != nil {
		cmd.Help()
		fmt.Println("")
		log.Fatal("failed to read private key: ", err)
	}

	_, err = x509.ParsePKCS8PrivateKey(caKeyBytes)
	if err != nil {
		log.Fatal("failed to parse private key: ", err)
	}

	name, _ := cmd.Flags().GetString("name")
	basename := path.Base(name)

	clientCert, clientKey, err := certtool.BuildClientCertificate(caCertBytes, caKeyBytes, name)
	if err != nil {
		log.Fatal("failed to build server certificate: ", err)
	}

	workdir, _ := cmd.Flags().GetString("workdir")

	certFile := path.Join(workdir, fmt.Sprintf("%s.crt", basename))
	if err := ovpncfg.WriteCert(clientCert, certFile); err != nil {
		log.Fatal("failed to write certificate: ", err)
	}

	keyFile := path.Join(workdir, fmt.Sprintf("%s.key", basename))
	if err := ovpncfg.WriteKey(clientKey, keyFile); err != nil {
		log.Fatal("failed to write key: ", err)
	}

	log.Printf(`Your new client certificate was successfully generated.`)
	log.Printf(`certificate: %q`, certFile)
	log.Printf(`private key: %q`, keyFile)
}

func init() {
	buildKeyCmd.Flags().String("name", "client", "Client's common name")
	buildKeyCmd.Flags().String("workdir", ".", "Work directory")
	buildKeyCmd.Flags().StringP("cert", "c", "ca.crt", "CA certificate path")
	buildKeyCmd.Flags().StringP("key", "k", "ca.key", "CA private key path")
}
