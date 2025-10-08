package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {

	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go handleConnection(serverConn)
	setString := "*3\r\n$3\r\nSET\r\n$4\r\npear\r\n$5\r\ngrape\r\n"

	_, err := clientConn.Write([]byte(setString))
	if err != nil {
		t.Error("Error parsing SET - couldn't write to client")
	}

	// Read client side response
	buffer := make([]byte, 1024)
	n, err := clientConn.Read(buffer)
	if err != nil {
		t.Errorf(" Error parsing SET")
	}

	response := string(buffer[:n])
	if response != "+OK\r\n" {
		t.Errorf(" Error parsing SET: response != +OK\r\n. Response %s", response)
	}

	fmt.Println("store now: ", store)
	fmt.Printf("store pointer (test): %p\n", &store)
	time.Sleep(50 * time.Millisecond)
	getString := "*2\r\n$3\r\nGET\r\n$4\r\npear\r\n"
	_, err = clientConn.Write([]byte(getString))
	if err != nil {
		t.Error("Error parsing GET - couldn't write request to server")
	}

	// Read client side response
	n, err = clientConn.Read(buffer)
	if err != nil {
		t.Errorf(" Error parsing GET")
	}

	response = string(buffer[:n])
	if response != "$5\r\ngrape\r\n" {
		t.Errorf(" Error parsing GET: response != \"$5\\r\\n\\grape\\r\\n\". Response %s", response)
	}

	if err != nil {

		t.Errorf("unexpected error: %v", err)
	}

}
