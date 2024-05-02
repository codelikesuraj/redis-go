package main

import (
	"bufio"
	"io"
	"os"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := Aof{file: f, rd: bufio.NewReader(f)}

	return &aof, nil
}

func (a *Aof) Close() error {
	return a.file.Close()
}

func (a *Aof) Write(value Value) error {
	_, err := a.file.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}

func (a *Aof) Read() error {
	for {
		val, err := NewRespReader(a.rd).Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		Handlers["SET"](val.arr[1:])
	}

	return nil
}
