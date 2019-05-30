package main

import (
	"encoding/pem"
	"io/ioutil"
)

func readPemFile(file string) ([]byte, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	pemBody, _ := pem.Decode(buf)
	return pemBody.Bytes, nil
}
