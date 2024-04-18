package main

import (
	"log"
	"net"
	"strings"
)

const ADDR = "localhost:6379"

func main() {
	// create a tcp listener/server
	listener, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatalln("failed to listen:", err.Error())
	}
	defer listener.Close()

	log.Println("Started listening for connections on", listener.Addr())

	// accept incoming connections
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln("failed to accept connection:", err.Error())
	}
	defer conn.Close()

	log.Println("Accepted connection from", conn.RemoteAddr())

	for {
		// read message from client
		resp := NewResp(conn)
		_, err = resp.Read()
		if err != nil {
			log.Fatalln("error reading from client:", err.Error())
		}

		// ignore request and send back "PONG"
		conn.Write([]byte("+PONG\r\n"))
	}

	// testOnlyNestedArrays()
}

func testOnlyNestedArrays() {
	// test nested array
	input := "*9\r\n"
	input += "$4\r\n"
	input += "this\r\n"
	input += "$2\r\n"
	input += "is\r\n"
	input += "+an\r\n"
	input += "+array\r\n"
	input += "*5\r\n"
	input += "+with\r\n"
	input += "+an\r\n"
	input += "+inside\r\n"
	input += "+value\r\n"
	input += "*3\r\n"
	input += "+also with an \r\n"
	input += "+inside number\r\n"
	input += ":4\r\n"
	input += "$3\r\n"
	input += "and\r\n"
	input += "$1\r\n"
	input += "a\r\n"
	input += "+simple string outside\r\n"
	input += "$4\r\n"
	input += "value\r\n"

	resp := NewResp(strings.NewReader(input))
	_, err := resp.Read()
	if err != nil {
		log.Fatalln("error reading from client:", err.Error())
	}
}
