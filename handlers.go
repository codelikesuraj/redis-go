package main

const (
	ERR_KEY_NOT_FOUND   = "key not found"
	ERR_CMD_REQ_ONE_ARG = "command requires one argument"
	ERR_CMD_REQ_TWO_ARG = "command requires two arguments"
)

var (
	SETs     = map[string]string{}
	Handlers = map[string]func([]Value) Value{
		"PING": ping,
		"SET":  set,
		"GET":  get,
	}
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{
			typ: SIMPLE_STRING,
			str: "PONG",
		}
	}

	return Value{
		typ:  BULK_STRINGS,
		bulk: args[0].bulk,
	}
}

func set(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: SIMPLE_ERRORS, err: ERR_CMD_REQ_TWO_ARG}
	}

	key, value := args[0].bulk, args[1].bulk

	SETs[key] = value

	return Value{typ: SIMPLE_STRING, str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: SIMPLE_ERRORS, err: ERR_CMD_REQ_ONE_ARG}
	}

	value, ok := SETs[args[0].bulk]
	if !ok {
		return Value{typ: SIMPLE_ERRORS, err: ERR_KEY_NOT_FOUND}
	}

	return Value{typ: BULK_STRINGS, bulk: value}
}
