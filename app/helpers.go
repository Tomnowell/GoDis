package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	read              inputState = iota
	echo              inputState = iota
	setDetect         inputState = iota
	setValue          inputState = iota
	getValue          inputState = iota
	optionDetect      inputState = iota
	optionValueDetect inputState = iota
	setList           inputState = iota
)

func encodeBulkString(s string) []byte {
	length := len(s)
	// convert the length to a string using strconv.Itoa
	bulkString := "$" + strconv.Itoa(length) + "\r\n" + s + "\r\n"
	return []byte(bulkString)
}

func encodeSimpleString(s string) []byte {
	simpleString := "+" + s + "\r\n"
	return []byte(simpleString)
}

func encodeRespInt(n int) RespInt {

	stringN := ":" + strconv.Itoa(n) + "\r\n"
	return RespInt(stringN)
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimRight(line, "\r\n")
	fmt.Println("line:", []byte(trimmed)) // debug: show CR if it exists
	return trimmed, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var tokenInterface interface{}
	var tokenErr error
	for {
		tokenInterface, tokenErr = parseRESP(reader)
		if tokenErr != nil {

			if tokenErr.Error() != "EOF" {
				fmt.Println("EOF ", tokenErr.Error())
			}
			fmt.Println("Error parsing token: ", tokenErr.Error())
			return
		}

		arr, ok := tokenInterface.([]interface{})
		if !ok {
			fmt.Println("Error parsing token: ", tokenErr.Error())
			return
		}
		processTokens(arr, conn)
	}
}

func processTokens(arr []interface{}, conn net.Conn) {
	// FSM partially implemented to make reading a little easier.
	// I think there's still more refactoring needed here.
	var state inputState = read

	key := ""
	value := ""
	option := ""
	optionValue := 0

	for _, element := range arr {
		fmt.Println(element)
		str := element.(string)
		switch state {
		case echo:
			state = read
			bulkString := encodeBulkString(str)
			_, err := conn.Write(bulkString)
			if err != nil {
				fmt.Println(" Error parsing echo")
			}

		case setDetect:
			state = setValue
			key = str

		case setValue:
			state = optionDetect
			value = str
			setStore(key, value)
			_, err := conn.Write(encodeSimpleString("OK"))
			if err != nil {
				fmt.Println(" Error parsing echo")
			}

		case optionDetect:
			state = optionValueDetect
			option = strings.ToUpper(str)

		case optionValueDetect:
			optionValue, _ = strconv.Atoi(str)
			switch option {
			case "EX":
				// Delayed delete (temporary key). Delay in seconds
				delay := time.Duration(optionValue) * time.Second
				go deleteStore(key, delay)

			case "PX":
				// Delayed delete (temporary key). Delay in milliseconds
				delay := time.Duration(optionValue) * time.Millisecond
				go deleteStore(key, delay)
			}

		case getValue:
			response := getStore(str)
			_, err := conn.Write(response)
			if err != nil {
				fmt.Println(" Error parsing GET")
			}
			state = read

		case setList:
			response, _ := makeList(str)
			_, err := conn.Write(response)

			if err != nil {
				fmt.Println("Error writing response")
			}

		case read:
			strUpper := strings.ToUpper(str)
			if strUpper == "PING" {
				i, err := conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					fmt.Println("Ërror parsing array")
					fmt.Println(err.Error())
					fmt.Println(i)
				}
			}
			if strUpper == "ECHO" {
				state = echo
			}

			if strUpper == "SET" {
				state = setDetect
			}

			if strUpper == "GET" {
				state = getValue
			}

			if strUpper == "RPUSH" {
				state = setList
			}
		}

	}
}

func temporaryKey() {
	// TODO refactor the EX and PX cases here if possible

}
