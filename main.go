package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	addr          = ":8080"
	serverTLSCert = "server.crt"
	serverTLSKey  = "server.key"
	clientTLSCert = "client.crt"
	clientTLSKey  = "client.key"
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "Address to listen on")
	flag.StringVar(&serverTLSCert, "tls-server-cert", serverTLSCert, "Path for server TLS cert")
	flag.StringVar(&serverTLSKey, "tls-server-key", serverTLSKey, "Path for server TLS key")
	flag.StringVar(&clientTLSCert, "tls-client-cert", clientTLSCert, "Path for client TLS cert")
	flag.StringVar(&clientTLSKey, "tls-client-key", clientTLSKey, "Path for client TLS key")
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			runClient()
			fmt.Println("")
			log.Println("client: trying again 5 seconds")
			time.Sleep(time.Second * 5)
		}
	}()
	wg.Wait()
}

func runServer() {
	cert, err := tls.LoadX509KeyPair(serverTLSCert, serverTLSKey)
	if err != nil {
		log.Fatalf("server: unable to load keys: %s", err)
	}

	tlsconf := &tls.Config{Certificates: []tls.Certificate{cert}}

	lis, err := tls.Listen("tcp", addr, tlsconf)
	if err != nil {
		log.Fatalf("server: unable to start listener: %s", err)
	}
	log.Printf("server listening on: %s", addr)

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("server: unable to access connections: %s", err)
		}
		defer conn.Close()
		log.Printf("server: got request from %s", conn.RemoteAddr())

		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("server: keys verified")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				pk, err := x509.MarshalPKIXPublicKey(v.PublicKey)
				if err != nil {
					log.Fatalf("server: unable to marshall public key %s", err)
				}
				log.Printf("server: public key: %s", pk)
			}
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}
		log.Printf("server: conn: echo %q\n", string(buf[:n]))
		n, err = conn.Write(buf[:n])

		n, err = conn.Write(buf[:n])
		log.Printf("server: conn: wrote %d bytes", n)

		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}

func runClient() {
	cert, err := tls.LoadX509KeyPair(clientTLSCert, clientTLSKey)
	if err != nil {
		log.Fatalf("client: error loading keys: %s", err)
	}
	tlsconf := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", addr, tlsconf)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting")
}
