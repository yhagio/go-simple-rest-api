package config

import (
	"crypto/rsa"
)

const (
	// For simplicity these files are in the same folder as the app binary.
	// You shouldn't do this in production.
	PrivKeyPath = "app.rsa"     // `> openssl genrsa -out app.rsa 1024`
	PubKeyPath  = "app.rsa.pub" // `> openssl rsa -in app.rsa -pubout > app.rsa.pub`
)

var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)