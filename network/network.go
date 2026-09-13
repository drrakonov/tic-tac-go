package network

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"net"
)

type AppState int

const (
	Menu AppState = iota
	Hosting
	Joining
	Playing
)

type Payload struct {
	Type string 
	Row  int
	Col  int
	Text string
}

func generateTLSConfig() (*tls.Config, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}

func HostGame(port string, pin string) (net.Conn, error) {
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		return nil, err
	}

	listener, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		return nil, err
	}
	
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	// SECURITY: Wait for the joiner to send the PIN
	p, err := ReceivePayload(conn)
	if err != nil || p.Type != "AUTH" || p.Text != pin {
		conn.Close()
		listener.Close()
		return nil, errors.New("unauthorized connection attempt")
	}

	// Lock the room so nobody else can connect!
	listener.Close()

	return conn, nil
}

func JoinGame(address string, pin string) (net.Conn, error) {
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}

	// SECURITY: Send the PIN immediately upon connecting
	err = SendPayload(conn, Payload{Type: "AUTH", Text: pin})
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func SendPayload(conn net.Conn, p Payload) error {
	encoder := json.NewEncoder(conn)
	return encoder.Encode(p)
}

func ReceivePayload(conn net.Conn) (Payload, error) {
	var p Payload
	// SECURITY: LimitReader caps the network read to 2048 bytes to prevent OOM DOS attacks
	decoder := json.NewDecoder(io.LimitReader(conn, 2048))
	err := decoder.Decode(&p)
	return p, err
}
