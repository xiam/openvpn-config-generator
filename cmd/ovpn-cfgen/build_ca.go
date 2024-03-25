package main

import (
	"fmt"
	"log"
	"path"

	"github.com/spf13/cobra"
	generator "github.com/xiam/openvpn-config-generator"
	"github.com/xiam/openvpn-config-generator/lib/certtool"
)

var buildCACmd = &cobra.Command{
	Use:   "build-ca [OPTIONS]",
	Short: "Create a self-signed CA Certificate",
	Run:   buildCAFn,
}

func buildCAFn(cmd *cobra.Command, args []string) {
	basename, _ := cmd.Flags().GetString("basename")
	workdir, _ := cmd.Flags().GetString("workdir")

	basename = path.Base(basename)

	caCert, caKey, err := certtool.BuildCA()
	if err != nil {
		log.Fatal("failed to build CA: ", err)
	}

	certFile := path.Join(workdir, fmt.Sprintf("%s.crt", basename))
	if err := generator.WriteCert(caCert, certFile); err != nil {
		log.Fatal("failed to write certificate: ", err)
	}

	keyFile := path.Join(workdir, fmt.Sprintf("%s.key", basename))
	if err := generator.WriteKey(caKey, keyFile); err != nil {
		log.Fatal("failed to write key: ", err)
	}

	log.Printf(`Your new CA certificate was successfully generated.`)
	log.Printf(`certificate: %q`, certFile)
	log.Printf(`private key: %q`, keyFile)
}

func init() {
	buildCACmd.Flags().String("basename", "ca", "Base name of the CA files (e.g.: {$basename}.{crt,key}).")
	buildCACmd.Flags().String("workdir", ".", "Work directory")
}
