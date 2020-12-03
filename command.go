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

func consumeCommands(actions Actions, cc chan Command, game *Game, mc chan<- Message) {
	for c := range cc {
		mc <- createMessage("Starting Command loop...")
		action, argument := c.Parse()
		switch c.origin {
		case USER:
			if actionFunction, ok := actions.user[action]; ok {
				actionFunction(argument, game, mc)
			} else {
				mc <- createMessage("Invalid Command")
			}
		case GAME:
			if actionFunction, ok := actions.game[action]; ok {
				actionFunction(argument, game, mc)
			} else {
				mc <- createMessage("Invalid Command")
			}
		default:
			mc <- createMessage("Invalid Command: Missing Origin")
		}
		mc <- createMessage("About to check triggers...")
		game.CheckTriggers(cc)
		mc <- createMessage("Checked triggers.")
		mc <- createMessage("Ending Command loop.")
	}
}
