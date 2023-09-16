# Go easy dkim signer

## How to use
### Generate keys and add txt-record to your domain
```shell
openssl genrsa -out private 2048
```
```shell
openssl rsa -in private -pubout -out public
```
```shell
sed '1d;$d' public | tr -d '\n' > spublic
```
Then make up a selector, i.e. `myselector` (could be any other string)<br>
Add txt-record to DNS of your domain:
- Key: `myselector._domainkey`
- Value: `v=DKIM1; k=rsa; p=<public key from spublic here>`

### Use the lib in your project
1. Add the lib to your go module: `go get github.com/metaer/go-easy-dkim-signer`
2. Use `easydkim.Sign` method like this:

#### Example 1
```go
package main

import (
	"log"
	"net/smtp"

	"github.com/metaer/go-easy-dkim-signer/easydkim"
)

func main() {
	var err error
	domain := "example.com"
	from := "example@example.com"
	rcpt := "rcpt@example.com"
	selector := "myselector"
	privateFileKeyPath := "private"

	message := []byte("Subject: test subject\r\n" +
		"To: " + rcpt + "\r\n" +
		"From: " + from + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		"Message body")

	message, err = easydkim.Sign(message, privateFileKeyPath, selector, domain)
	if err != nil {
		log.Fatal(err)
	}

	err = smtp.SendMail("localhost:1525", nil, from, []string{rcpt}, message)
	if err != nil {
		log.Fatal(err)
	}
}

```

#### Example 2
Install `go get gopkg.in/gomail.v2` for object-oriented way to create an email message
```go
package main

import (
	"bytes"
	"log"
	"net/smtp"

	"github.com/metaer/go-easy-dkim-signer/easydkim"
	"gopkg.in/gomail.v2"
)

func main() {
	var err error
	domain := "example.com"
	from := "example@example.com"
	rcpt := "rcpt@example.com"
	selector := "myselector"
	privateFileKeyPath := "private"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", rcpt)
	m.SetHeader("Subject", "My subject")
	m.SetBody("text/html", "My body")

	var buffer bytes.Buffer
	_, err = m.WriteTo(&buffer)
	if err != nil {
		log.Fatal(err)
	}
	var signedMessage []byte
	signedMessage, err = easydkim.Sign(buffer.Bytes(), privateFileKeyPath, selector, domain)
	if err != nil {
		log.Fatal(err)
	}
	err = smtp.SendMail("localhost:1525", nil, from, []string{rcpt}, signedMessage)
	if err != nil {
		log.Fatal(err)
	}
}
```

### How to test locally
1. Run `openssl genrsa -out private 2048`
2. Run `docker-compose up -d`
3. Run `example1` or `example2`
4. Open `http://localhost:1580`
5. Check dkim signature in message source