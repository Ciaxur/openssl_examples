package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Helper function that parses a certificate from a raw PEM bytes array.
func ParseCertificateFromPEMByptes(pemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemBytes))
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate from PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

func main() {
	// Flags for mTLS.
	serverCrt := flag.String("cert", "", "Path to the server's certificate.")
	serverKey := flag.String("key", "", "Path to the server's private key.")
	caCrtPath := flag.String("ca", "", "Path to the CA's certificate.")
	host := flag.String("host", "localhost", "Server's hostname endpoint.")
	port := flag.Uint("port", 3000, "Server's port number endpoint.")
	flag.Parse()

	// Check required flags where passed in.
	if *serverCrt == "" || *serverKey == "" || *caCrtPath == "" {
		log.Fatal("required arguments: server and CA cert and key are required")
	}

	// Create the CA pool which will be used for verifying the client with.
	caCrtContent, err := ioutil.ReadFile(*caCrtPath)
	if err != nil {
		panic(fmt.Errorf("failed to read the content of CA %s", *caCrtPath))
	}

	// Parse the certificate(s) from PEM bytes.
	caCrt, err := ParseCertificateFromPEMByptes(caCrtContent)
	if err != nil {
		panic(fmt.Errorf("failed to parse certifacte from PEM bytes: %v", err))
	}

	// Construct a list of parsed trusted CAs.
	trustedCerts := []x509.Certificate{}
	trustedCerts = append(trustedCerts, *caCrt)

	caCertPool := x509.NewCertPool()
	for _, trustedCert := range trustedCerts {
		caCertPool.AddCert(&trustedCert)
	}

	// Create server.
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", *host, *port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			ServerName: "localhost",
			ClientCAs:  caCertPool,
			// Require and verify the client's cert against it being signed by the CA.
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS12,
		},
	}

	router := mux.NewRouter()

	// Add server root endpoints.
	http.Handle("/", router)

	// Basic ping endpoint.
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Client request: [Handshake=%v][ServerName=%v][RemoteAddr=%v]\n", r.TLS.HandshakeComplete, r.TLS.ServerName, r.RemoteAddr)
		w.Write([]byte("pong"))
	}).Methods("GET")

	log.Printf("Listening on %s:%d.\n", *host, *port)
	if err := server.ListenAndServeTLS(*serverCrt, *serverKey); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
