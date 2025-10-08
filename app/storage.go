package main

import (
	"fmt"
)

const (
	nullBulkString = "$-1\r\n"
)

var store = make(map[string]string)

func setStore(key string, value string) {
	// Initial behaviour: overwrite any existing values
	store[key] = value
	fmt.Println(store)
}

func getStore(key string) []byte {
	value, ok := store[key]
	fmt.Printf("getStore(%q) ok=%v value=%q store=%v\n", key, ok, value, store)
	if ok {
		fmt.Printf("store pointer(storage): %p\n", &store)
		return encodeBulkString(value)
	}
	fmt.Printf("store pointer(storage): %p\n", &store)
	return []byte(nullBulkString)
}
