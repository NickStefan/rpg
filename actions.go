package main

import "time"

func createActions() map[string]func(string, chan<- Result) {
	actions := make(map[string]func(string, chan<- Result))
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

func actionInit(arg string, rc chan<- Result) {
	rc <- createResult(newline)
	rc <- createResult("(click the Chat Input...)")
	pacing()
	rc <- createResult("(type '/start' ...)")
	pacing()
	rc <- createResult("(and hit Enter)", newline)
}

func actionStart(arg string, rc chan<- Result) {
	rc <- createResult("Where there was once a house, there is now nothing.")
	pacing()
	rc <- createResult("Where there was once wood, there is now only charcoal.")
	pacing()
	rc <- createResult("Where there was once no choices to be made, there are now many...", newline)

	pacing(6)

	rc <- createResult("You can't stay long--the invaders could be back any minute.")
	pacing()
	rc <- createResult("You'll have to choose between your last few belongings found in the fire:", newline)
	pacing()
	rc <- createResult(" - A roll of [burned bandages]")
	pacing()
	rc <- createResult(" - A [blackened sword]")
	pacing()
	rc <- createResult(" - An [indecipherable book]", newline)

	pacing(4)

	rc <- createResult("(type /inspect [item] e.g. '/inspect indecipherable book' to learn more)")
	pacing()
	rc <- createResult("(type /take [item] to choose an item and flee)", newline)
}

func actionInspect(arg string, rc chan<- Result) {
	// TODO validate their argument is something nearby our character Entity...

	actionTaken := "You inspect the " + arg
	rc <- createResult(actionTaken)
	pacing()

	switch arg {
	case "burned bandages":
		rc <- createResult(
			"They appear to be wrapped around several bottles that jangle inside--like a makeshift pack.",
		)
		pacing()
		rc <- createResult(
			"You're afraid to unwrap it further without somewhere safe to inspect the contents. They could be tonics or salves--perfect for healing the wounded.",
		)
	case "blackened sword":
		rc <- createResult(
			"It appears to be burned--and dull.",
		)
		pacing()
		rc <- createResult(
			"It'll be heavy to carry, but also heavy to receive a blow from--perfect for protecting oneself on hostile roads.",
		)
	case "indecipherable book":
		rc <- createResult(
			"The words are pronounceable--maybe. It's some language you aren't familiar with.",
		)
		pacing()
		rc <- createResult(
			"Seems impracticle to take with you, but maybe it's valuable. After all, it was your grandfather's most prized possession.",
		)
	}
}

func actionTake(arg string, rc chan<- Result) {
	rc <- createResult("You take the " + arg + " and flee.")
}

func actionEcho(arg string, rc chan<- Result) {
	rc <- createResult(arg)
}
