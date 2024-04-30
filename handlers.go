package main

var SETs = map[string]string{}

var Handlers = map[string]func([]Value) Value{
	"PING": pong,
	"SET":  set,
	"GET":  get,
}

func pong(args []Value) Value {
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
		return Value{typ: SIMPLE_ERRORS, err: "'set' command should include a key and value arguments"}
	}

	key, value := args[0].bulk, args[1].bulk

	SETs[key] = value

	return Value{typ: SIMPLE_STRING, str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: SIMPLE_ERRORS, err: "'get' command should include only a key argument"}
	}

	value, ok := SETs[args[0].bulk]
	if !ok {
		return Value{typ: SIMPLE_ERRORS, err: "key not found"}
	}

	return Value{typ: SIMPLE_STRING, str: value}
}
