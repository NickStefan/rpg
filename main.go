package main

import "time"

func gameLoop(game *Game, actions Actions, cc chan Command, mc chan<- Message) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case c := <-cc:
			action, argument := c.Parse()
			if actionFunction, ok := actions.user[action]; ok {
				actionFunction(argument, game, mc)
			} else {
				mc <- createMessage("Invalid Command")
			}
		case <-ticker.C:
			game.CheckTriggers(cc)
		}
	}
}

func tellStory(mc chan<- Message) {
	time.Sleep(2 * time.Second)
	mc <- createMessage(
		"Where there was once a house, there is now nothing.",
	)
	time.Sleep(6 * time.Second)
	mc <- createMessage(
		"Where there was once wood, there is now only charcoal glowing orange in the night.",
	)
	time.Sleep(10 * time.Second)
	mc <- createMessage(
		"You can still hear the clanking of swords, the howling of invaders some ways away. ",
		"If you move quickly, you might make it to the gatehouse and safety inside the castle. ",
		"If not, the invaders will be back shortly to take prisoners...\n",
	)
	time.Sleep(20 * time.Second)
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
	go gameLoop(game, actions, cc, mc)
	go func() {
		tellStory(mc)
		game.Spawn()
	}()
	createUI(cc, mc)
}
