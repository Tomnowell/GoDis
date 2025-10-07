package main

import (
	"bufio"
	"fmt"
	"strconv"
)

const (
	Integer    = ':'
	String     = '+'
	BulkString = '$'
	Array      = '*'
	Error      = '-'
)

func parseRESP(reader *bufio.Reader) (interface{}, error) {

	prefix, err := reader.ReadByte()
	fmt.Println("prefix:", string(prefix))
	if err != nil {
		fmt.Println("Error reading prefix")
		return "", err
	}

	switch prefix {
	case Integer:
		fmt.Println("Parsing Integer")
		return readInteger(reader)

	case String:
		fmt.Println("Parsing String")
		return readSimpleString(reader)

	case BulkString:
		fmt.Println("Parsing BulkString")
		return readBulkString(reader)

	case Array:
		fmt.Println("Parsing array")
		return readArray(reader)

	case Error:
		fmt.Println("Parsing error")
		return readError(reader)

	default:
		return "", fmt.Errorf("Unrecognized response type: %v", prefix)
	}
}

func readBulkString(reader *bufio.Reader) (interface{}, error) {
	line, _ := readLine(reader)

	length, _ := strconv.Atoi(line)
	if length < 0 {
		return "", fmt.Errorf("invalid length: %d", length)
	}

	data := make([]byte, length)
	_, err := reader.Read(data)
	if err != nil {
		return "", err
	}

	// Consume CLRF
	_, err = reader.Discard(2)
	if err != nil {
		return "", err
	}
	fmt.Println(string(data))
	return string(data), nil
}

func readArray(reader *bufio.Reader) (interface{}, error) {
	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	length, _ := strconv.Atoi(line)
	fmt.Println(string(line), length)
	if length < 0 {
		return "", fmt.Errorf("invalid length: %d", length)
	}
	arr := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		val, err := parseRESP(reader)
		if err != nil {
			return "", err
		}
		arr = append(arr, val)
	}
	return arr, nil
}

func readSimpleString(reader *bufio.Reader) (interface{}, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	return line, nil
}

func readInteger(reader *bufio.Reader) (interface{}, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	val, err := strconv.Atoi(line)
	if val > 9223372036854775807 || val < -9223372036854775808 {
		return nil, fmt.Errorf("invalid integer: %v", val)
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func readError(reader *bufio.Reader) (interface{}, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	// Usually treated just like a string, but may be wrapped in a special type
	return fmt.Errorf("%s", line), nil
}
