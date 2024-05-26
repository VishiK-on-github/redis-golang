package main

import (
	"sync"
)

// handler function that maps the uppercase command
// to appropriate function
var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

// ping command, returns pong.
// used to check if the service is up or not
func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

// the main storage for set, get methods
var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	// filtering args, we need key and value for set
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	// getting key and value
	key := args[0].bulk
	value := args[1].bulk

	// locking the variable, writing to map and then unlocking
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	// filtering args, we need key for set
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	// getting the key
	key := args[0].bulk

	// setting read lock reading value and unlocking
	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	// filtering args, we need key, value and hash for hset
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	// getting required params
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	// locking, creating a new map, storing key, value in map and unlocking
	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	// filtering args, we need key and hash for hget
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	// getting hash and key
	hash := args[0].bulk
	key := args[1].bulk

	// read locking, getting hash set, getting
	// value from hash using key, unlocking
	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	// filtering args, we need hash for hgetall
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	// getting required hash
	hash := args[0].bulk

	// read locking, getting value for hash, unlocking
	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}

	// iterating over results, appending for returning
	for key, value := range value {
		values = append(values, Value{typ: "bulk", bulk: key})
		values = append(values, Value{typ: "bulk", bulk: value})
	}

	return Value{typ: "array", array: values}
}
