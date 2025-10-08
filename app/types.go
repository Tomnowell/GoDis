package main

type token struct {
	dataType rune
	command  string
	body     string
}

type tokenStack []token

type inputState int
