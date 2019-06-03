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

var buildKeyServerCmd = &cobra.Command{
	Use:   "build-key-server [OPTIONS]",
	Short: "Create and sign a server certificate",
	Run:   buildKeyServerFn,
}

func buildKeyServerFn(cmd *cobra.Command, args []string) {
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

	serverCert, serverKey, err := certtool.BuildServerCertificate(caCertBytes, caKeyBytes, name)
	if err != nil {
		log.Fatal("failed to build server certificate: ", err)
	}

	workdir, _ := cmd.Flags().GetString("workdir")

	certFile := path.Join(workdir, fmt.Sprintf("%s.crt", basename))
	if err := ovpncfg.WriteCert(serverCert, certFile); err != nil {
		log.Fatal("failed to write certificate: ", err)
	}

	keyFile := path.Join(workdir, fmt.Sprintf("%s.key", basename))
	if err := ovpncfg.WriteKey(serverKey, keyFile); err != nil {
		log.Fatal("failed to write key: ", err)
	}

	log.Printf(`Your new server certificate was successfully generated.`)
	log.Printf(`certificate: %q`, certFile)
	log.Printf(`private key: %q`, keyFile)
}

func init() {
	buildKeyServerCmd.Flags().String("name", "server", "Server's common name")
	buildKeyServerCmd.Flags().String("workdir", ".", "Work directory")
	buildKeyServerCmd.Flags().StringP("cert", "c", "ca.crt", "CA certificate path")
	buildKeyServerCmd.Flags().StringP("key", "k", "ca.key", "CA private key path")
}
