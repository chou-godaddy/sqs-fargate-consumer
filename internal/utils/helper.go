package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	gosecrets "github.com/gdcorp-domains/fulfillment-gosecrets"
)

// RetrieveCert retriees cert and key from secret manager
func RetrieveCert(certName string, keyName string) (*tls.Certificate, error) {
	secretRetriever := gosecrets.NewSecretRetriever()

	certBytes, certErr := secretRetriever.Get(context.Background(), gosecrets.SecretConfig{
		AWS: &gosecrets.AWSSecretConfig{
			Name:   certName,
			Region: "us-west-2",
		},
	})
	if certErr != nil {
		return nil, certErr
	}
	keyBytes, keyErr := secretRetriever.Get(context.Background(), gosecrets.SecretConfig{
		AWS: &gosecrets.AWSSecretConfig{
			Name:   keyName,
			Region: "us-west-2",
		},
	})
	if keyErr != nil {
		return nil, keyErr
	}
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func StartHTTPSServer() {
	// Retrieve cert and key from secret manager
	registrarSvcCert, err := RetrieveCert("registrar.dev.client.int.godaddy.com.crt", "registrar.dev.client.int.godaddy.com.key")
	if err != nil {
		fmt.Print("error retrieving cert and key")
		panic(err)
	}

	go func() {
		fmt.Println("Creating https server")
		s := &http.Server{
			Addr:    ":443",
			Handler: nil,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{*registrarSvcCert},
			},
		}

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		fmt.Println("Https server is listening on 443 with TLS")

		if err := s.ListenAndServeTLS("", ""); err != nil {
			fmt.Print("error listening")
			panic(err)
		}

	}()
}
