package main

import (
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

// the "wrapped" widget
func createWrapped() *text.Text {
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}
	if err := wrapped.Write("Supports", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write(" colors", text.WriteCellOpts(cell.FgColor(cell.ColorNumber(33)))); err != nil {
		panic(err)
	}
	return wrapped
}

// ChatLogUI widget
func createChatLogUI(mc <-chan Message) *text.Text {
	rolled, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			for m := range mc {
				var err error
				for _, txt := range m.texts {
					if txt.text != "" {
						err = rolled.Write(txt.text)
					}
				}
				err = rolled.Write("\n")
				if err != nil {
					panic(err)
				}
			}
		}
	}()
	return rolled
}

// ChatInputUI widget
func createChatInputUI(cc chan<- Command) *textinput.TextInput {
	input, err := textinput.New(
		textinput.FillColor(cell.ColorDefault),
		textinput.TextColor(cell.ColorDefault),
		textinput.ClearOnSubmit(),
		textinput.OnSubmit(func(text string) error {
			if text != "" {
				cc <- Command{text: text, origin: USER}
			}
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}
	return input
}

func createUI(cc chan<- Command, mc <-chan Message) {
	ctx, cancel := context.WithCancel(context.Background())

	wrapped := createWrapped()

	chatInputUI := createChatInputUI(cc)
	chatLogUI := createChatLogUI(mc)

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
