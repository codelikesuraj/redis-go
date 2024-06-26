# Redis-Go - A simple implementation of a REDIS server in Golang

## Usage
- ### Start the Redis-Go server
    ```go run main.go```
- ### Run your redis commands
    ```redis-cli [command] [args...]```

## Todo
    
- ✅ create a tcp listener
- ✅ parse incoming RESP data
    - ✅ simple string
    - ✅ bulk string
    - ✅ simple error
    - ✅ integer
    - ✅ arrays
        - ✅ single data type
        - ✅ multiple data type
- ✅ format outgoing RESP data
    - ✅ simple string
    - ✅ bulk string
    - ✅ simple error
    - ✅ integer
    - ✅ arrays
    - ✅ null
        - ✅ single data type
        - ✅ multiple data type
- ✅ accept basic commands
    - ✅ SET {data} {value}
    - ✅ GET {data}
    - ✅ PING
- ✅ setup data persistence (using Append-Only File)