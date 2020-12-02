package main

import "time"

func createActions() map[string]func(string, chan<- Message) {
	actions := make(map[string]func(string, chan<- Message))
	actions["/_init"] = actionInit
	actions["/start"] = actionStart
	actions["/echo"] = actionEcho
	actions["/inspect"] = actionInspect
	actions["/take"] = actionTake
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

func actionInit(arg string, mc chan<- Message) {
	mc <- createMessage(newline)
	mc <- createMessage("(click the Chat Input...)")
	pacing()
	mc <- createMessage("(type '/start' ...)")
	pacing()
	mc <- createMessage("(and hit Enter)", newline)
}

func actionStart(arg string, mc chan<- Message) {
	mc <- createMessage("Where there was once a house, there is now nothing.")
	pacing()
	mc <- createMessage("Where there was once wood, there is now only chamcoal.")
	pacing()
	mc <- createMessage("Where there was once no choices to be made, there are now many...", newline)

	pacing(6)

	mc <- createMessage("You can't stay long--the invaders could be back any minute.")
	pacing()
	mc <- createMessage("You'll have to choose between your last few belongings found in the fire:", newline)
	pacing()
	mc <- createMessage(" - A roll of [burned bandages]")
	pacing()
	mc <- createMessage(" - A [blackened sword]")
	pacing()
	mc <- createMessage(" - An [indecipherable book]", newline)

	pacing(4)

	mc <- createMessage("(type /inspect [item] e.g. '/inspect indecipherable book' to learn more)")
	pacing()
	mc <- createMessage("(type /take [item] to choose an item and flee)", newline)
}

func actionInspect(arg string, mc chan<- Message) {
	// TODO validate their argument is something nearby our character Entity...

	actionTaken := "You inspect the " + arg
	mc <- createMessage(actionTaken)
	pacing()

	switch arg {
	case "burned bandages":
		mc <- createMessage(
			"They appear to be wrapped around several bottles that jangle inside--like a makeshift pack.",
		)
		pacing()
		mc <- createMessage(
			"You're afraid to unwrap it further without somewhere safe to inspect the contents. They could be tonics or salves--perfect for healing the wounded.",
		)
	case "blackened sword":
		mc <- createMessage(
			"It appears to be burned--and dull.",
		)
		pacing()
		mc <- createMessage(
			"It'll be heavy to carry, but also heavy to receive a blow from--perfect for protecting oneself on hostile roads.",
		)
	case "indecipherable book":
		mc <- createMessage(
			"The words are pronounceable--maybe. It's some language you aren't familiar with.",
		)
		pacing()
		mc <- createMessage(
			"Seems impracticle to take with you, but maybe it's valuable. After all, it was your grandfather's most prized possession.",
		)
	}
}

func actionTake(arg string, mc chan<- Message) {
	mc <- createMessage("You take the " + arg + " and flee.")
}

func actionEcho(arg string, mc chan<- Message) {
	mc <- createMessage(arg)
}
