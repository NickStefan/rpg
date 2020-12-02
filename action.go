package main

import "time"

type Actions struct {
	user map[string]func(string, *Game, chan<- Message)
	game map[string]func(string, *Game, chan<- Message)
}

func createActions() Actions {
	actions := Actions{
		game: map[string]func(string, *Game, chan<- Message){
			"/init": actionInit,
			"/find": actionFind,
		},
		user: map[string]func(string, *Game, chan<- Message){
			"/start":   actionStart,
			"/echo":    actionEcho,
			"/inspect": actionInspect,
			"/take":    actionTake,
		},
	}
	return actions
}

const newline = "\n"

func pacing(seconds ...int) {
	var pace int
	if len(seconds) == 0 {
		pace = 3
	} else {
		pace = seconds[0]
	}
	time.Sleep(time.Duration(pace) * time.Second)
}

func actionInit(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage(newline)
	pacing()
	mc <- createMessage("(click the Chat Input, type '/start', and hit Enter)")
}

func actionStart(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage("Where there was once a house, there is now nothing.")
	pacing()
	mc <- createMessage("Where there was once wood, there is now only charcoal.")
	pacing()
	mc <- createMessage("Take only what you can carry--the invaders could be back any minute.")
	pacing()
	mc <- createMessage(newline)

	game.Zone.addItem(Item{
		LocatedEntity{
			name:        "roll of burned bandages",
			description: "They appear to be wrapped around several bottles--like a makeshift pack--probably tonics and salves meant for healing. Probably.",
		},
	})
	game.Zone.addItem(Item{
		LocatedEntity{
			name:        "blackened sword",
			description: "It appears to be burned--and dull. It'll be heavy to carry, but also heavy to receive a blow from--perfect for protecting oneself on hostile roads.",
		},
	})
	game.Zone.addItem(Item{
		LocatedEntity{
			name:        "indecipherable book",
			description: "The words are pronounceable--maybe. It's some language you aren't familiar with. Seems impracticle to take with you, but maybe it's valuable. Maybe.",
		},
	})

	pacing()
	mc <- createMessage(newline)
	mc <- createMessage("Use '/inspect [item]' to look more closely at any items you find.")
	mc <- createMessage("Use '/take [item]' to take one. Keep in mind: you'll have limited carrying space.")
}

func actionFind(arg string, game *Game, mc chan<- Message) {
	//item := game.Zone.getItem(arg)
	//mc <- createMessage("You find a " + item.name)
}

func actionInspect(arg string, game *Game, mc chan<- Message) {
	// TODO validate their argument is something nearby our character Entity...

	actionTaken := "You inspect the " + arg
	mc <- createMessage(actionTaken)

	switch arg {
	case "burned bandages":
		mc <- createMessage(
			"They appear to be wrapped around several bottles that jangle inside--like a makeshift pack.",
		)
		mc <- createMessage(
			"You're afraid to unwrap it further without somewhere safe to inspect the contents. They could be tonics or salves--perfect for healing the wounded.",
		)
	case "blackened sword":
		mc <- createMessage(
			"It appears to be burned--and dull.",
		)
		mc <- createMessage(
			"It'll be heavy to carry, but also heavy to receive a blow from--perfect for protecting oneself on hostile roads.",
		)
	case "indecipherable book":
		mc <- createMessage(
			"The words are pronounceable--maybe. It's some language you aren't familiar with.",
		)
		mc <- createMessage(
			"Seems impracticle to take with you, but maybe it's valuable. After all, it was your grandfather's most prized possession.",
		)
	}
}

func actionTake(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage("You take the " + arg + " and flee.")
}

func actionEcho(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage(arg)
}
