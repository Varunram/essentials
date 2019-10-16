package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
)

// ref: https://ericchiang.github.io/post/go-tls/

// helper function to create a cert template with a serial number and other required fields
func CertTemplate() (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Println("failed to generate serial number: ", err)
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Yhat, Inc."}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func createCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certpem []byte, err error) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		log.Println(err)
		return
	}
	// parse the resulting certificate so we can use it again
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		log.Println(err)
		return
	}
	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certpem = pem.EncodeToMemory(&b)
	return
}

func GenCert() (string, string, tls.Certificate, error) {
	// generate a new key-pair
	var keypem []byte
	var tlscert tls.Certificate

	certTemplate, err := CertTemplate()
	if err != nil {
		log.Println(err)
		return "", "", tlscert, err
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
		return "", "", tlscert, err
	}

	certTemplate.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	certTemplate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	certTemplate.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	cert, certpem, err := createCert(certTemplate, certTemplate, &key.PublicKey, key)
	if err != nil {
		log.Println(err)
		return "", "", tlscert, err
	}

	log.Printf("%s\n", certpem)
	log.Printf("%#x\n", cert.Signature)

	// PEM encode the private key
	keypem = pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	// Create a TLS cert using the private key and certificate
	tlscert, err = tls.X509KeyPair(certpem, keypem)
	if err != nil {
		log.Println(err)
		return "", "", tlscert, err
	}

	certString := fmt.Sprintf("%s", certpem)
	return string(keypem), certString, tlscert, nil
}
