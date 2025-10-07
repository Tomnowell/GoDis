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

		for _, element := range arr {
			fmt.Println(element)
			str := element.(string)
			bulkString := string(encodeBulkString(str))
			if echo {
				echo = false
				_, err := conn.Write([]byte(bulkString))
				if err != nil {
					fmt.Println(" Error parsing echo")
				}
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

		}
	}

}
