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
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("failed to accept connection:", err.Error())
		}
		defer conn.Close()

		// log.Println("Accepted connection from", conn.RemoteAddr())	// debugging info

		// read message from client
		respReader := NewRespReader(conn)
		val, err := respReader.Read()
		if err != nil {
			log.Fatalln("error reading from client:", err.Error())
		}

		var respVal Value
		if val.typ == SIMPLE_STRING && strings.ToLower(val.str) == "ping" {
			respVal.typ = val.typ
			respVal.str = "PONG"
		} else if val.typ == BULK_STRINGS && strings.ToLower(val.bulk) == "ping" {
			respVal.typ = val.typ
			respVal.str = "PONG"
		} else if val.typ == ARRAYS && val.arr[0].typ == BULK_STRINGS && strings.ToLower(val.arr[0].bulk) == "ping" {
			respVal.typ = SIMPLE_STRING
			respVal.str = "PONG"
		} else {
			respVal.typ = SIMPLE_STRING
			respVal.str = "IT WORKS"
		}

		// ignore request and send back "PONG"
		// write message to client
		err = NewRespWriter(conn).Write(respVal)
		if err != nil {
			log.Fatalln(err.Error())
		}

		// log.Println("Closed connection with", conn.RemoteAddr())	// debugging info
	}
}
