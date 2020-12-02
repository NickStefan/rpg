package main

import (
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"strings"
)

type Command struct {
	text string
}

func (c Command) ParseCommand() (action string, argument string) {
	splitString := strings.SplitN(c.text, " ", 2)
	if len(splitString) == 2 {
		action, argument = splitString[0], splitString[1]
	} else if len(splitString) == 1 {
		action = splitString[0]
	}
	return action, argument
}

type ColorText struct {
	text  string
	color string
}

type Message struct {
	texts []ColorText
}

func createMessage(texts ...string) Message {
	var _texts []ColorText
	for _, t := range texts {
		_texts = append(_texts, ColorText{text: t})
	}
	return Message{texts: _texts}
}

func consumeCommands(actions map[string]func(string, chan<- Message), cc <-chan Command, mc chan<- Message) {
	for {
		for c := range cc {
			action, argument := c.ParseCommand()
			if actionFunction, ok := actions[action]; ok {
				go actionFunction(argument, mc)
			} else {
				mc <- createMessage("Invalid Command")
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wrapped := createWrapped()

	cc := make(chan Command, 1)
	chatInputUI := createChatInputUI(cc)

	mc := make(chan Message, 1)
	chatLogUI := createChatLogUI(mc)

	actions := createActions()
	go consumeCommands(actions, cc, mc)

	cc <- Command{"/_init"}

	rootUI, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer rootUI.Close()

	innerUI, err := container.New(
		rootUI,
		container.Border(linestyle.Light),
		container.BorderTitle("RPG"),
		container.SplitVertical(
			container.Left(
				container.SplitHorizontal(
					container.Top(
						container.Border(linestyle.Light),
						container.BorderTitle("Game Log"),
						container.PlaceWidget(chatLogUI),
					),
					container.Bottom(
						container.Border(linestyle.Light),
						container.BorderTitle("Chat Input"),
						container.PlaceWidget(chatInputUI),
					),
					container.SplitPercent(85),
				),
			),
			container.Right(
				container.Border(linestyle.Light),
				container.BorderTitle("Wraps lines at rune boundaries"),
				container.PlaceWidget(wrapped),
			),
			container.SplitPercent(70),
		),
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, rootUI, innerUI, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}
