package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
)

func Test4() {
	//var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: "fransebas",
		Auth: []ssh.AuthMethod{
			ssh.Password("Riemann"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()

	// Request the remote side to open port 8080 on all interfaces.
	l, err := conn.Listen("tcp", "0.0.0.0:60")
	if err != nil {
		log.Fatal("unable to register tcp forward: ", err)
	}
	defer l.Close()

	// Serve HTTP with your SSH server acting as a reverse proxy.
	http.Serve(l, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "Hello world!\n")
	}))
}
