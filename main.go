package main

import "time"

func gameLoop(game *Game, actions Actions, cc chan Command, mc chan<- Message) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case c := <-cc:
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
		case <-ticker.C:
			game.CheckTriggers(cc)
		}
	}
}

func main() {
	cc := make(chan Command, 1)
	mc := make(chan Message, 1)

	game := &Game{
		playerHealth:            80,
		playerLoc:               0,
		playerInventoryCapacity: 1,
		playerInventory:         make(map[string]Item),
		itemsByLoc:              make(map[int][]Item),
		npcsByLoc:               make(map[int][]NPC),
	}
	actions := createActions()
	//go consumeCommands(actions, cc, game, mc)
	go gameLoop(game, actions, cc, mc)
	cc <- Command{"/init", GAME}

	createUI(cc, mc)
}
