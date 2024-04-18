package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	SIMPLE_STRING = '+'
	SIMPLE_ERRORS = '-'
	INTEGERS      = ':'
	BULK_STRINGS  = '$'
	ARRAYS        = '*'
)

type Resp struct {
	arr_depth int
	reader    *bufio.Reader
}

type Value struct {
	typ  string
	str  string
	err  string
	num  int
	bulk string
	arr  []Value
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
		return r.readArray()
	default:
		fmt.Printf("Unknown type %c\n", b)
	}
	return Value{}, nil
}

func (r *Resp) readSimpleString() (Value, error) {
	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	fmt.Printf("This is a simple string %q\n", string(l))

	return Value{
		typ: "simple string",
		str: string(l),
	}, nil
}

func (r *Resp) readBulkString() (Value, error) {
	n, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	fmt.Printf("This a bulk string %q of size %d bytes\n", string(l), n)

	return Value{
		typ:  "bulk string",
		bulk: string(l),
	}, nil
}

func (r *Resp) readInt() (Value, error) {
	i, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	fmt.Printf("This is an integer '%d'\n", i)

	return Value{
		typ: "integer",
		num: i,
	}, nil
}

func (r *Resp) readError() (Value, error) {
	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	fmt.Printf("This is an error %q\n", string(l))

	return Value{
		typ: "error",
		err: string(l),
	}, nil
}

func (r *Resp) readArray() (Value, error) {
	r.arr_depth += 1
	defer func() {
		r.arr_depth -= 1
	}()

	n, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	val := []Value{}

	fmt.Println("This is an array with", n, "elements")
	for i := range n {
		fmt.Printf("%s[%d]", strings.Repeat("\t", r.arr_depth), i+1)

		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		val = append(val, v)
	}

	return Value{
		typ: "array",
		arr: val,
	}, nil
}
