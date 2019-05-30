package main

import (
	"fmt"
	"github.com/spf13/cobra"
	ovpncfg "github.com/xiam/openvpn-config-generator"
	"github.com/xiam/openvpn-config-generator/lib/certtool"
	"log"
	"os"
	"path"
)

var buildCACmd = &cobra.Command{
	Use:   "build-ca [OPTIONS]",
	Short: "Create a self-signed CA Certificate",
	Run:   buildCAFn,
}

func buildCAFn(cmd *cobra.Command, args []string) {
	basename, _ := cmd.Flags().GetString("basename")

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
	buildCACmd.Flags().StringP("basename", "b", "ca", "base name of the CA files, e.g.: {$basename}.{crt,key}")
}
