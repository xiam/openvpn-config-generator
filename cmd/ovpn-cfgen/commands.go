package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ovpn-cfgen [COMMAND] [OPTIONS]",
	Short: "ovpn-cfgen is a configuration file generator for OpenVPN",
	Run:   rootCmdFn,
}

func rootCmdFn(cmd *cobra.Command, args []string) {
	log.Printf("help")
}
