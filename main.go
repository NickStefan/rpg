package main

type Game struct {
	*Zone
	*Character
}

func main() {
	cc := make(chan Command, 1)
	mc := make(chan Message, 1)

	game := &Game{
		Character: &Character{health: 80, maxHealth: 100},
		Zone:      &Zone{maxY: 10, maxX: 0}, // make it a linear map for now...
	}

	actions := createActions()
	go consumeCommands(actions, cc, game, mc)

	cc <- Command{"/init", GAME}

	createUI(cc, mc)
}
