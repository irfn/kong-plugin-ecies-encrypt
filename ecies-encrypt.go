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

func (conf Config) PrivateKey() *ecies.PrivateKey {
	b64pk, err := ioutil.ReadFile(conf.PrivateKeyFile)
	if err != nil {
		log.Fatalf("no config %s", err)
	}

	decodedKey, err := base64.StdEncoding.DecodeString(string(b64pk))
	if err != nil {
		log.Fatalf("Bad key configured: %s, %s", b64pk, err.Error())
	}

	return ecies.NewPrivateKeyFromBytes([]byte(decodedKey))
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	pk := conf.PrivateKey()
	path, err := kong.Request.GetPath()
	if err != nil {
		log.Printf("Error accessing request path info: %s", err.Error())
	}
	if path == "/pubkey" {
		kong.Response.Exit(200, base64.RawStdEncoding.EncodeToString(pk.PublicKey.Bytes(false)), map[string][]string{})
	}
}

func (conf Config) Request(kong *pdk.PDK) {
	pk := conf.PrivateKey()
	requestBody, err := kong.Request.GetRawBody()
	if err != nil {
		log.Printf("Error reading request body: %s", err.Error())
	}
	decrptedRequest, err := ecies.Decrypt(pk, requestBody)
	kong.ServiceRequest.SetRawBody(string(decrptedRequest))
}

func (conf Config) Response(kong *pdk.PDK) {
	path, err := kong.Request.GetPath()
	if err != nil {
		log.Printf("Error accessing request path info: %s", err.Error())
	}
	if path == "/pubkey" {
		kong.Response.ExitStatus(200)
	}

	key, err := kong.Request.GetHeader("x-pub-key")
	decodedKey, err := base64.StdEncoding.DecodeString(string(key))
	clientPubKey, err := ecies.NewPublicKeyFromBytes(decodedKey)
	if err != nil {
		log.Printf("Error accessing request client pub key: %s", err.Error())
	}

	body, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		log.Printf("Error reading response body: %s", err.Error())
	}

	encryptedBody, err := ecies.Encrypt(clientPubKey, []byte(body))
	if err != nil {
		log.Printf("Error encrypting response body: %s", err.Error())
	}

	upstreamStatus, _ := kong.ServiceResponse.GetStatus()
	upstreamHeaders, _ := kong.ServiceResponse.GetHeaders(-1)

	kong.Response.Exit(upstreamStatus, base64.RawStdEncoding.EncodeToString(encryptedBody), upstreamHeaders)
}
