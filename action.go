package main

import "time"
import "strings"

type Actions struct {
	user map[string]func(string, *Game, chan<- Message)
	game map[string]func(string, *Game, chan<- Message)
}

func createActions() Actions {
	actions := Actions{
		game: map[string]func(string, *Game, chan<- Message){
			"/init":        actionInit,
			"/find":        actionFind,
			"/attackByNPC": actionAttackByNPC,
			"/deathByNPC":  actionDeathByNPC,
		},
		user: map[string]func(string, *Game, chan<- Message){
			"/start":     actionStart,
			"/echo":      actionEcho,
			"/inspect":   actionInspect,
			"/inventory": actionInventory,
			"/take":      actionTake,
		},
	}
	return actions
}

const newline = "\n"

func pacing(seconds ...int) {
	var pace int
	if len(seconds) == 0 {
		pace = 0
	} else {
		pace = seconds[0]
	}
	time.Sleep(time.Duration(pace) * time.Second)
}

func actionInit(arg string, game *Game, mc chan<- Message) {
	pacing()
	mc <- createMessage(newline, "(click the Chat Input, type '/start', and hit Enter)")
}

func actionStart(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage(newline, "Where there was once a house, there is now nothing.")
	pacing()
	mc <- createMessage("Where there was once wood, there is now only charcoal.")
	pacing()
	mc <- createMessage("Take only what you can carry--the invaders could be back any minute.")

	pacing()
	mc <- createMessage(newline, "Use '/inspect [item]' to look more closely at any items you find.")
	mc <- createMessage("Use '/take [item]' to take one. Keep in mind: you'll have limited carrying space.")

	game.AddItem(1, Item{
		name:        "roll of burned bandages",
		description: "They appear to be wrapped around several bottles--like a makeshift pack--probably tonics and salves meant for healing. Probably.",
	})
	game.AddItem(1, Item{
		name:        "blackened sword",
		description: "It appears to be burned--and dull. It'll be heavy to carry, but also heavy to receive a blow from--perfect for protecting oneself on hostile roads.",
	})
	game.AddItem(1, Item{
		name:        "indecipherable book",
		description: "The words are pronounceable--maybe. It's some language you aren't familiar with. Seems impracticle to take with you, but maybe it's valuable. Maybe.",
	})
	game.playerLoc = 1
}

func actionFind(arg string, game *Game, mc chan<- Message) {
	if item, ok := game.GetItem(game.playerLoc, arg); ok {
		mc <- createMessage("You find a [" + item.name + "].")

	} else if npc, ok := game.GetNPC(game.playerLoc, arg); ok {
		mc <- createMessage("Just ahead, you see a [" + npc.name + "]. You might need to /inspect it to be sure.")

	} else {
		mc <- createMessage("You find nothing.")
	}
}

func actionInspect(arg string, game *Game, mc chan<- Message) {
	if item, ok := game.GetItem(game.playerLoc, arg); ok {
		mc <- createMessage("You inspect the " + arg + ". " + item.description)
	} else if npc, ok := game.GetNPC(game.playerLoc, arg); ok {
		mc <- createMessage("You consider the " + arg + ". " + npc.description)
	} else {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
	}
}

func actionTake(arg string, game *Game, mc chan<- Message) {
	ok, err := game.TakeItem(arg)
	if ok {
		mc <- createMessage("You take the " + arg + ".")
	} else if err == CAPACITY {
		mc <- createMessage("You can't carry anything more.")
	} else if err == NOTFOUND {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
	}

	game.AddNPC(2, NPC{
		name:         "invading warrior",
		description:  "He looks dangerous and will attack any second.",
		weaponDamage: 30,
	})
	game.playerLoc++
}

func actionInventory(arg string, game *Game, mc chan<- Message) {
	if item, ok := game.GetInventoryItem(arg); ok {
		mc <- createMessage("You inspect the " + arg + ". " + item.description)
	} else {
		mc <- createMessage("There isn't anything by that name here. Maybe it's gone?")
	}
}

func actionAttackByNPC(arg string, game *Game, mc chan<- Message) {
	splitString := strings.SplitN(arg, " ", 2)
	dmg, name := splitString[0], splitString[1]
	mc <- createMessage("A " + name + " hits your for " + dmg + " damage.")
}

func actionDeathByNPC(arg string, game *Game, mc chan<- Message) {
	splitString := strings.SplitN(arg, " ", 2)
	dmg, name := splitString[0], splitString[1]
	mc <- createMessage("A " + name + " hits your for " + dmg + " damage AND KILLS YOU.")
}

func actionEcho(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage(arg)
}
