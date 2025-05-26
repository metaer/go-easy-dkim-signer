package easydkim

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/emersion/go-msgauth/dkim"
)

func Sign(data []byte, dkimPrivateKeyFilePath string, selector string, domain string) ([]byte, error) {
	privateKeyBytes, err := os.ReadFile(dkimPrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

	return signWithBytesPrivateKey(data, privateKeyBytes, selector, domain)
}

func SignWithStringPrivateKey(data []byte, privateKeyString string, selector string, domain string) ([]byte, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %v", err)
	}

	return signWithBytesPrivateKey(data, privateKeyBytes, selector, domain)
}

func signWithBytesPrivateKey(data []byte, privateKeyBytes []byte, selector string, domain string) ([]byte, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	result, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey := result.(*rsa.PrivateKey)

	options := &dkim.SignOptions{
		Domain:   domain,
		Selector: selector,
		Signer:   privateKey,
	}

	msg := bytes.NewReader(data)

	var b bytes.Buffer
	if err := dkim.Sign(&b, msg, options); err != nil {
		return nil, fmt.Errorf("dkim signing error: %v", err)
	}

	return b.Bytes(), nil
}
