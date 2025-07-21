package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func RunServer() error {
	listener, err := net.Listen("tcp", "127.0.0.1:502")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Listening on 127.0.0.1:502")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
		} else {
			fmt.Printf("Connection read error: %v\n", err)
		}
		return
	}
	fmt.Printf("Received: %x\n", buf[:n])
}
