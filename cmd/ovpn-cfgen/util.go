package main

import (
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

func checkFile(cmd *cobra.Command, file string, message string) {
	_, err := os.Stat(file)
	if err != nil {
		cmd.Help()
		fmt.Println("")
		log.Fatal(fmt.Sprintf("%s: could not stat file %v", message, err))
	}
}

func readPemFile(file string) ([]byte, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	pemBody, _ := pem.Decode(buf)
	return pemBody.Bytes, nil
}
