package main

import "strings"

type Command struct {
	text   string
	origin string
}

const USER = "USER"
const GAME = "GAME"

func (c Command) Parse() (action string, argument string) {
	splitString := strings.SplitN(c.text, " ", 2)
	if len(splitString) == 2 {
		action, argument = splitString[0], splitString[1]
	} else if len(splitString) == 1 {
		action = splitString[0]
	}
	return action, argument
}
