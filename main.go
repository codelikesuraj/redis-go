package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

const (
	ADDR          = "localhost:6379"
	SIMPLE_STRING = '+'
	SIMPLE_ERRORS = '-'
	INTEGERS      = ':'
	BULK_STRINGS  = '$'
	ARRAYS        = '*'
)

type Value struct {
	typ  string
	str  string
	err  string
	num  int
	bulk string
	arr  []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return []byte{}, 0, nil
		}

		line = append(line, b)
		n++

		if b == '\n' && line[len(line)-2] == '\r' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (int, error) {
	line, _, err := r.readLine()
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func (r *Resp) Read() (Value, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch b {
	case SIMPLE_STRING:
		return r.readSimpleString()
	case SIMPLE_ERRORS:
		return r.readError()
	case INTEGERS:
		return r.readInt()
	case BULK_STRINGS:
		return r.readBulkString()
	case ARRAYS:
		fmt.Println("This is an array")
	default:
		fmt.Printf("Unknown type %c\n", b)
	}
	return Value{}, nil
}

func (r *Resp) readSimpleString() (Value, error) {
	fmt.Println("This is a simple string")

	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("The simple string is:", string(l))

	return Value{
		typ: "simple string",
		str: string(l),
	}, nil
}

func (r *Resp) readBulkString() (Value, error) {
	fmt.Println("This is a bulk string")

	n, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("The bulk string is:", string(l), "with size", n, "bytes")

	return Value{
		typ:  "bulk string",
		bulk: string(l),
	}, nil
}

func (r *Resp) readInt() (Value, error) {
	fmt.Println("This is an integer")

	i, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("The integer is:", i)

	return Value{
		typ: "integer",
		num: i,
	}, nil
}

func (r *Resp) readError() (Value, error) {
	fmt.Println("This is an error")

	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("The error is:", string(l))

	return Value{
		typ: "error",
		str: string(l),
	}, nil
}

func main() {
	input := "-this is an error\r\n"
	resp := NewResp(strings.NewReader(input))
	_, err := resp.Read()
	if err != nil {
		log.Fatalln("error reading from client:", err.Error())
	}
	// create a tcp listener/server
	// listener, err := net.Listen("tcp", ADDR)
	// if err != nil {
	// 	log.Fatalln("failed to listen:", err.Error())
	// }
	// defer listener.Close()

	// log.Println("Started listening for connections on", listener.Addr())

	// // accept incoming connections
	// conn, err := listener.Accept()
	// if err != nil {
	// 	log.Fatalln("failed to accept connection:", err.Error())
	// }
	// defer conn.Close()

	// log.Println("Accepted connection from", conn.RemoteAddr())

	// for {
	// 	// read message from client
	// 	resp := NewResp(conn)
	// 	_, err = resp.Read()
	// 	if err != nil {
	// 		log.Fatalln("error reading from client:", err.Error())
	// 	}

	// 	// ignore request and send back "PONG"
	// 	conn.Write([]byte("+PONG\r\n"))
	// }
}
