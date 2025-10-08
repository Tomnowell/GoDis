package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
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

		echo := false
		set := false
		get := false
		setValueTrigger := false
		key := ""
		value := ""

		for _, element := range arr {
			fmt.Println(element)
			str := element.(string)
			bulkString := encodeBulkString(str)
			if echo {
				echo = false
				_, err := conn.Write(bulkString)
				if err != nil {
					fmt.Println(" Error parsing echo")
				}
			}
			if set {
				set = false
				setValueTrigger = true
				key = str
			} else if setValueTrigger {
				setValueTrigger = false
				value = str
				setStore(key, value)
				_, err := conn.Write(encodeSimpleString("OK"))
				if err != nil {
					fmt.Println(" Error parsing echo")
				}
			}
			if get {
				r := getStore(str)
				_, err := conn.Write([]byte(r))
				if err != nil {
					fmt.Println(" Error parsing GET")
				}
				get = false
			}

			if strings.ToUpper(str) == "PING" {
				i, err := conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					fmt.Println("Ërror parsing array")
					fmt.Println(err.Error())
					fmt.Println(i)
					return
				}
			}
			if strings.ToUpper(str) == "ECHO" {
				echo = true
			}

			if strings.ToUpper(str) == "SET" {
				set = true
			}

			if strings.ToUpper(str) == "GET" {
				get = true
			}

		}
	}

}
