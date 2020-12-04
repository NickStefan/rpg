package main

import "strings"

type Actions struct {
	user map[string]func(string, *Game, chan<- Message)
}

func createActions() Actions {
	actions := Actions{
		user: map[string]func(string, *Game, chan<- Message){
			"/echo":        actionEcho,
			"/inspect":     actionInspect,
			"/inventory":   actionInventory,
			"/take":        actionTake,
			"/find":        actionFind,
			"/attackByNPC": actionAttackByNPC,
			"/deathByNPC":  actionDeathByNPC,
		},
	}
	return actions
}

func actionEcho(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage(arg)
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

func actionFind(arg string, game *Game, mc chan<- Message) {
	if item, ok := game.GetItem(game.playerLoc, arg); ok {
		mc <- createMessage("You find a [" + item.name + "].")

	} else if npc, ok := game.GetNPC(game.playerLoc, arg); ok {
		mc <- createMessage("Just ahead, you see a [" + npc.name + "]. You might need to /inspect it to be sure.")

	} else {
		mc <- createMessage("You find nothing.")
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
