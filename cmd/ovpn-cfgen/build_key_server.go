package main

import (
	"fmt"
	"github.com/spf13/cobra"
	ovpncfg "github.com/xiam/openvpn-config-generator"
	"github.com/xiam/openvpn-config-generator/lib/certtool"
	"ioutil"
	"log"
	"os"
	"path"
)

var buildKeyServerCmd = &cobra.Command{
	Use:   "build-key-server [OPTIONS]",
	Short: "Create and sign a server certificate",
	Run:   buildKeyServerFn,
}

func buildKeyServerFn(cmd *cobra.Command, args []string) {
	certFile, _ := cmd.Flags().GetString("ca-cert")
	certKey, _ := cmd.Flags().GetString("ca-key")

	certBody, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatal("failed to read certificate: ", err)
	}
	certPem, _ := pem.Decode(certBody)

	caCert, err := x509.ParseCertificate(certPem.Bytes)

	basename, _ := cmd.Flags().GetString("commonname")
	basename = path.Base(basename)

	caCert, caKey, err := certtool.BuildCA()
	if err != nil {
		log.Fatal("failed to build CA: ", err)
	}

	basedir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to retrieve base dir: ", err)
	}

	certFile := path.Join(basedir, fmt.Sprintf("%s.crt", basename))
	if err := ovpncfg.WriteCert(caCert, certFile); err != nil {
		log.Fatal("failed to write certificate: ", err)
	}

	keyFile := path.Join(basedir, fmt.Sprintf("%s.key", basename))
	if err := ovpncfg.WriteKey(caKey, keyFile); err != nil {
		log.Fatal("failed to write key: ", err)
	}

	log.Printf(`Your new CA certificate was successfully generated.`)
	log.Printf(`certificate: %q`, certFile)
	log.Printf(`private key: %q`, keyFile)
}

func init() {
	buildKeyServerCmd.Flags().StringP("commonname", "cn", "server", "common name of the certificate")
	buildKeyServerCmd.Flags().StringP("ca-cert", "cc", "ca.crt", "CA certificate path")
	buildKeyServerCmd.Flags().StringP("ca-key", "ck", "ca.key", "CA private key path")
}
