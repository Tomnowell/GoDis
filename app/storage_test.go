package main

import (
	"testing"
)

func TestSetStore(t *testing.T) {
	key := "foo"
	value := "bar"

	setStore(key, value)
	if store[key] != value {
		t.Errorf("failed to set store")
	}
}

func TestGetStore(t *testing.T) {
	key := "foo"
	value := "bar"
	bulkString := encodeBulkString(value)
	setStore(key, value)
	result := getStore(key)

	if string(result) != string(bulkString) {
		t.Errorf("failed to getValue store")
	}
}

func TestMakeList(t *testing.T) {

	key := "foo"
	response, list := makeList(key)
	if string(response) != string(encodeRespInt(1)) {
		t.Errorf("failed response, %s", response)
	}
	if list.Front().Value != "foo" {

		t.Errorf("Failed list: %s", list.Front().Value)
	}

}
