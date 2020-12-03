package main

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
	go consumeCommands(actions, cc, game, mc)
	cc <- Command{"/init", GAME}

	createUI(cc, mc)
}
