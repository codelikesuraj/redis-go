# Redis-Go - A simple implementation of a REDIS server in Golang

## Usage
- go run main.go

## Todo
- [x] create a tcp listener
- [ ] parse incoming RESP data
    - [x] simple string
    - [x] bulk string
    - [x] simple error
    - [x] integer
    - [x] arrays
        - [x] single data type
        - [x] multiple data type
- [ ] format outgoing RESP data
    - [ ] simple string
    - [ ] bulk string
    - [ ] simple error
    - [ ] integer
    - [ ] arrays
        - [ ] single data type
        - [ ] multiple data type
- [ ] accept basic commands
    - [ ] SET {data} {value}
    - [ ] GET {data}
- [ ] setup data persistence (in-memory database)