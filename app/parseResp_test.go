package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestParseResp(t *testing.T) {
	str := "*1\r\n$4\r\nPING\r\n"
	reader := bufio.NewReader(strings.NewReader(str))

	received, err := parseRESP(reader)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(received)

	arr, ok := received.([]interface{})
	if !ok {
		t.Errorf("expected []interface{}, got %T", received)
	}
	expected := []interface{}{"PING"}

	for i := range expected {
		if arr[i] != expected[i] {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], arr[i])
		}
	}

	str = "*2\r\n$4\r\nECHO\r\n$6\r\norange\r\n"
	reader = bufio.NewReader(strings.NewReader(str))
	received, err = parseRESP(reader)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(received)
	arr, ok = received.([]interface{})
	if !ok {
		t.Errorf("expected []interface{}, got %T", received)
	}
	expected = []interface{}{"ECHO", "orange"}
	for i := range expected {
		if arr[i] != expected[i] {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], arr[i])
		}
	}
}

func TestReadBulkString(t *testing.T) {
	str := "5\r\nhello\r\n"
	emptyStr := "0\r\n\r\n"

	reader := bufio.NewReader(strings.NewReader(str))
	reader2 := bufio.NewReader(strings.NewReader(emptyStr))

	received, err := readBulkString(reader)
	if err != nil {
		t.Error(err)
	}

	if received != "hello" {
		t.Errorf("received: %s, expected: hello", received)
	}

	received, err = readBulkString(reader2)
	if err != nil {
		t.Error(err)
	}
	if received != "" {
		t.Errorf("received: %s, expected: \"\"", received)
	}
}

func TestReadArray(t *testing.T) {
	str := "2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	reader := bufio.NewReader(strings.NewReader(str))

	received, err := readArray(reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert type: we expect readArray to return a []interface{}
	arr, ok := received.([]interface{})
	if !ok {
		t.Errorf("expected []interface{}, got %T", received)
	}

	expected := []interface{}{"hello", "world"}

	if len(arr) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(arr))
	}

	for i := range expected {
		if arr[i] != expected[i] {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], arr[i])
		}
	}
}

func TestReadString(t *testing.T) {
	str := "hello\r\n"
	reader := bufio.NewReader(strings.NewReader(str))

	received, err := readSimpleString(reader)
	expected := "hello"
	if err != nil {
		t.Errorf("unexpected error in read String: %v", err)
	}
	if received != expected {
		t.Errorf("received: %s, expected: %s", received, expected)
	}
}

func TestReadInteger(t *testing.T) {
	str := "+5\r\n"
	reader := bufio.NewReader(strings.NewReader(str))
	received, err := readInteger(reader)
	if err != nil {
		t.Errorf("unexpected error in read Integer: %v", err)
	}
	expected := 5
	if received != expected {
		t.Errorf("received: %d, expected: %d", received, expected)
	}

	str = "-5\r\n"
	reader = bufio.NewReader(strings.NewReader(str))
	received, err = readInteger(reader)
	if err != nil {
		t.Errorf("unexpected error in read Integer: %v", err)
	}
	expected = -5
	if received != expected {
		t.Errorf("received: %d, expected: %d", received, expected)
	}
	str = "42\r\n"
	reader = bufio.NewReader(strings.NewReader(str))
	received, err = readInteger(reader)
	if err != nil {
		t.Errorf("unexpected error in read Integer: %v", err)
	}
	expected = 42
	if received != expected {
		t.Errorf("received: %d, expected: %d", received, expected)
	}

	str = "2147483647\r\n"
	reader = bufio.NewReader(strings.NewReader(str))
	received, err = readInteger(reader)
	if err != nil {
		t.Errorf("unexpected error in read Integer: %v", err)
	}
	expected = 2147483647
	if received != expected {
		t.Errorf("received: %d, expected: %d", received, expected)
	}
	str = "9223372036854775807\r\n"
	reader = bufio.NewReader(strings.NewReader(str))
	received, err = readInteger(reader)
	if err != nil {
		t.Errorf("unexpected error in read Integer: %v", err)
	}
	expected = 9223372036854775807
	if received != expected {
		t.Errorf("received: %d, expected: %d", received, expected)
	}

	str = "9223372036854775808\r\n"
	reader = bufio.NewReader(strings.NewReader(str))
	received, err = readInteger(reader)
	if err != nil {
		if _, ok := err.(*strconv.NumError); !ok { // Check specifically for strconv.NumError to handle invalid integer cases properly
			t.Errorf("invalid integer not caught: %v", err)
		}
	} else {
		t.Errorf("invalid integer should have been caught but was not: %v", received)
	}

}
