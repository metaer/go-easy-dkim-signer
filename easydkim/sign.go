package easydkim

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/emersion/go-msgauth/dkim"
)

func Sign(data []byte, dkimPrivateKeyFilePath string, selector string, domain string) ([]byte, error) {
	msg := bytes.NewReader(data)
	privateKeyBytes, err := os.ReadFile(dkimPrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

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

	var b bytes.Buffer
	if err := dkim.Sign(&b, msg, options); err != nil {
		return nil, fmt.Errorf("dkim signing error: %v", err)
	}

	return b.Bytes(), nil
}
