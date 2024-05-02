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

	// create persistent storage
	aof, err := NewAof("database.aof")
	if err != nil {
		log.Fatalln(err)
	}
	defer aof.Close()

	// accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("failed to accept connection:", err.Error())
		}

		// log.Println("Accepted connection from", conn.RemoteAddr())	// debugging info

		// read message from client
		respReader := NewRespReader(conn)
		val, err := respReader.Read()
		if err != nil {
			log.Fatalln("error reading from client:", err.Error())
		}

		respWriter := NewRespWriter(conn)

		if val.typ != ARRAYS {
			respWriter.Write(Value{typ: SIMPLE_ERRORS, err: "invalid request, expected array"})
			conn.Close()
			continue
		}

		if len(val.arr) == 0 {
			respWriter.Write(Value{typ: SIMPLE_ERRORS, err: "invalid request, expected array length > 0"})
			conn.Close()
			continue
		}

		command := strings.ToUpper(val.arr[0].bulk)
		args := val.arr[1:]

		handler, ok := Handlers[command]
		if !ok {
			respWriter.Write(Value{typ: SIMPLE_ERRORS, err: "invalid command: " + command})
			conn.Close()
			continue
		}

		res := handler(args)

		// ignore request and send back "PONG"
		// write message to client
		err = NewRespWriter(conn).Write(res)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if command == "SET" && res.typ != SIMPLE_ERRORS {
			aof.Write(val)
		}

		// log.Println("Closed connection with", conn.RemoteAddr())	// debugging info
		conn.Close()
	}
}
