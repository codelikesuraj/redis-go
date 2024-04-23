# Redis-Go - A simple implementation of a REDIS server in Golang

## Usage
- go run main.go

## Todo
- [x] create a tcp listener
- [x] parse incoming RESP data
    - [x] simple string
    - [x] bulk string
    - [x] simple error
    - [x] integer
    - [x] arrays
        - [x] single data type
        - [x] multiple data type
- [x] format outgoing RESP data
    - [x] simple string
    - [x] bulk string
    - [x] simple error
    - [x] integer
    - [x] arrays
    - [x] null
        - [x] single data type
        - [x] multiple data type
- [ ] accept basic commands
    - [ ] SET {data} {value}
    - [ ] GET {data}
- [ ] setup data persistence (in-memory database)