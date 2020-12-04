package main

type Actions struct {
	user map[string]func(string, *Game, chan<- Message)
}

func createActions() Actions {
	actions := Actions{
		user: map[string]func(string, *Game, chan<- Message){
			"/echo":      actionEcho,
			"/inspect":   actionInspect,
			"/inventory": actionInventory,
			"/take":      actionTake,
			"/consider":  actionConsider,
			"/attack":    actionAttack,
		},
	}
	return actions
}

func actionEcho(arg string, game *Game, mc chan<- Message) {
	mc <- createMessage(arg)
}

func actionInspect(arg string, game *Game, mc chan<- Message) {
	game.InspectItem(arg, mc)
}

func actionTake(arg string, game *Game, mc chan<- Message) {
	game.TakeItem(arg, mc)
	game.playerLoc++
}

func actionInventory(arg string, game *Game, mc chan<- Message) {
	game.ListInventory(mc)
}

func actionConsider(arg string, game *Game, mc chan<- Message) {
	game.ConsiderNPC(arg, mc)
}

func actionAttack(arg string, game *Game, mc chan<- Message) {
	game.ToggleAttack(arg, mc)
}
