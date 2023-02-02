/*
	A "hello world" plugin in Go,
	which reads a request header and sets a response header.
*/

package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	ecies "github.com/ecies/go/v2"
)

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.2"
var Priority = 1

type Config struct {
	Message        string
	PrivateKeyFile string `json:"private_key_file"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Response(kong *pdk.PDK) {
	b64pk, err := ioutil.ReadFile(conf.PrivateKeyFile)
	if err != nil {
		log.Fatalf("no config %s", err)
	}
	log.Printf(string(b64pk))
	decodedKey, err := base64.StdEncoding.DecodeString(string(b64pk))
	if err != nil {
		log.Fatalf("Bad key configured: %s, %s", b64pk, err.Error())
	}

	pk := ecies.NewPrivateKeyFromBytes([]byte(decodedKey))
	body, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		log.Printf("Error reading response body: %s", err.Error())
	}

	encryptedBody, err := ecies.Encrypt(pk.PublicKey, []byte(body))
	if err != nil {
		log.Printf("Error encrypting response body: %s", err.Error())
	}

	upstreamStatus, _ := kong.ServiceResponse.GetStatus()
	upstreamHeaders, _ := kong.ServiceResponse.GetHeaders(-1)
	kong.Response.Exit(upstreamStatus, base64.RawStdEncoding.EncodeToString(encryptedBody), upstreamHeaders)
}
