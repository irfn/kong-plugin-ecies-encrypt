package main

import (
	"log"

	"io/ioutil"

	b64 "encoding/base64"

	ecies "github.com/ecies/go/v2"
)

func main() {
	k, err := ecies.GenerateKey()
	if err != nil {
		panic(err)
	}

	log.Println("key pair has been generated")

	// write the whole body at once
	pKey := b64.StdEncoding.EncodeToString(k.Bytes())

	err = ioutil.WriteFile("ecies.pk.key", []byte(pKey), 0644)
	if err != nil {
		panic(err)
	}
}
