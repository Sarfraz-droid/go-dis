package resp

import (
	"strings"
	"sync"
)

func HandleValue(v Value) (string, []Value) {
	command := strings.ToUpper(v.array[0].bulk)
	args := v.array[1:]

	return command, args
}

var db = map[string]string{}
var dbMutex = sync.RWMutex{}

var hdb = map[string]map[string]string{}
var hdbMutex = sync.RWMutex{}

var Handlers = map[string]func([]Value) Value{
	"COMMAND": command,
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HGET":    hget,
	"HSET":    hset,
	"HGETALL": hgetAll,
}

func command(args []Value) Value {
	return Value{typ: "string", str: "OK"}
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	print(key, value)

	dbMutex.Lock()
	db[key] = value
	dbMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	dbMutex.RLock()
	value, ok := db[key]
	dbMutex.RUnlock()

	print(value, ok)

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "string", str: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	key := args[0].bulk
	field := args[1].bulk
	value := args[2].bulk

	hdbMutex.Lock()
	_, ok := hdb[key]

	if !ok {
		hdb[key] = map[string]string{}
	}
	hdb[key][field] = value
	hdbMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	hdbMutex.RLock()
	value, ok := hdb[key][value]
	hdbMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "string", str: value}
}

func hgetAll(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetAll' command"}
	}

	key := args[0].bulk

	var data []Value

	hdbMutex.RLock()
	for k := range hdb[key] {
		data = append(data, Value{
			typ: "string",
			str: hdb[key][k],
		})
	}
	hdbMutex.RUnlock()

	return Value{
		typ:   "array",
		array: data,
	}
}
