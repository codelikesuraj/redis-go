package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	SIMPLE_ERRORS = '-'
	INTEGERS      = ':'
	BULK_STRINGS  = '$'
	ARRAYS        = '*'
	NULL          = '_'
)

type RespReader struct {
	arr_depth int
	reader    *bufio.Reader
}

type RespWriter struct {
	writer io.Writer
}

type Value struct {
	typ  byte
	str  string
	err  string
	num  int
	bulk string
	arr  []Value
}

func NewRespReader(rd io.Reader) *RespReader {
	return &RespReader{reader: bufio.NewReader(rd)}
}

func NewRespWriter(wr io.Writer) *RespWriter {
	return &RespWriter{writer: wr}
}

func (r *RespReader) readLine() (line []byte, n int, err error) {
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

func (r *RespReader) readInteger() (int, error) {
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

func (r *RespReader) Read() (Value, error) {
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
		return Value{}, errors.New("unknown value")
	}
}

func (r *RespReader) readSimpleString() (Value, error) {
	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	// fmt.Printf("This is a simple string %q\n", string(l))	// debugging info

	return Value{
		typ: SIMPLE_STRING,
		str: string(l),
	}, nil
}

func (r *RespReader) readBulkString() (Value, error) {
	// n, err := r.readInteger()	// debugging info

	_, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	// fmt.Printf("This a bulk string %q of size %d bytes\n", string(l), n)	// debuggin info

	return Value{
		typ:  BULK_STRINGS,
		bulk: string(l),
	}, nil
}

func (r *RespReader) readInt() (Value, error) {
	i, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	// fmt.Printf("This is an integer '%d'\n", i)	// debugging info

	return Value{
		typ: INTEGERS,
		num: i,
	}, nil
}

func (r *RespReader) readError() (Value, error) {
	l, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	// fmt.Printf("This is an error %q\n", string(l))	// debuggin info

	return Value{
		typ: SIMPLE_ERRORS,
		err: string(l),
	}, nil
}

func (r *RespReader) readArray() (Value, error) {
	r.arr_depth += 1
	defer func() {
		r.arr_depth -= 1
	}()

	n, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	val := []Value{}

	// fmt.Println("This is an array with", n, "elements")	// debugging info
	// for i := range n {	// debugging info
	// 	fmt.Printf("%s[%d]", strings.Repeat("\t", r.arr_depth), i+1)	// debugging info

	for range n {
		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		val = append(val, v)
	}

	return Value{
		typ: ARRAYS,
		arr: val,
	}, nil
}

func (v Value) marshalSimpleString() []byte {
	var bytes []byte
	bytes = append(bytes, SIMPLE_STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalBulkString() []byte {
	var bytes []byte
	bytes = append(bytes, BULK_STRINGS)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk[:len(v.bulk)]...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalSimpleError() []byte {
	var bytes []byte
	bytes = append(bytes, SIMPLE_ERRORS)
	bytes = append(bytes, v.err...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGERS)
	bytes = append(bytes, strconv.Itoa(v.num)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalArray() []byte {
	len := len(v.arr)
	var bytes []byte
	bytes = append(bytes, ARRAYS)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')
	for i := range len {
		bytes = append(bytes, v.arr[i].Marshal()...)
	}
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("_\r\n")
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case SIMPLE_STRING:
		return v.marshalSimpleString()
	case BULK_STRINGS:
		return v.marshalBulkString()
	case SIMPLE_ERRORS:
		return v.marshalSimpleError()
	case INTEGERS:
		return v.marshalInteger()
	case ARRAYS:
		return v.marshalArray()
	case NULL:
		return v.marshalNull()
	default:
		// fmt.Printf("Unknown type %q\n", v.typ)	// debugging info
		return []byte("-ERROR unknown type\r\n")
	}
}

func (wr *RespWriter) Write(v Value) error {
	bytes := v.Marshal()

	// debugging info
	// for b := range bytes {
	// 	fmt.Printf("%c", b)
	// }

	_, err := wr.writer.Write(bytes)
	if err != nil {
		log.Println("Error writing to client: ", err.Error())
		return err
	}
	return nil
}
