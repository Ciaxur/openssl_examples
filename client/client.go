package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// Flags for mTLS.
	clientCrt := flag.String("cert", "", "Path to the client's certificate.")
	clientKey := flag.String("key", "", "Path to the client's private key.")
	caCrt := flag.String("ca", "", "Path to the CA's certificate.")
	host := flag.String("host", "localhost", "Server's hostname endpoint.")
	port := flag.Uint("port", 3000, "Server's port number endpoint.")
	flag.Parse()
	srvEndpoint := fmt.Sprintf("%s:%d", *host, *port)

	// Check required flags where passed in.
	if *clientCrt == "" || *clientKey == "" || *caCrt == "" {
		log.Fatal("required arguments: client and CA cert and key are required")
	}

	// Load client's key pair.
	cert, err := tls.LoadX509KeyPair(*clientCrt, *clientKey)
	if err != nil {
		log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", *clientCrt, *clientKey)
	}

	// Load the CA that authorized the server's certs.
	log.Printf("CA Cert: %s\n", *caCrt)
	caCrtContent, err := ioutil.ReadFile(*caCrt)
	if err != nil {
		log.Fatalf("Could not read the contents of CA cert %s\n", *caCrt)
	}

	// Create a CA certificate pool, in order for the certificate to be
	// validated.
	caCrtPool := x509.NewCertPool()
	caCrtPool.AppendCertsFromPEM(caCrtContent)

	// Create the TLS for the client.
	t := http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCrtPool,
		},
	}

	// Create the client with the tls transport and invoke a request to the
	// server.
	client := http.Client{
		Transport: &t,
		Timeout:   5 * time.Second,
	}

	// Send request.
	endpoint := fmt.Sprintf("https://%s/ping", srvEndpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBufferString(""))
	if err != nil {
		log.Fatal("failed to construct request")
		return
	}

	log.Println("Invoking a request to endpoint:", endpoint)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to invoke GET request with server: ", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status from the server %d with content: %s\n", resp.StatusCode, body)

}
