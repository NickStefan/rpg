package main

import (
	"context"
	"fmt"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
	"math/rand"
	"time"
)

// quotations are used as text that is rolled up in a text widget.
var quotations = []string{
	"When some see coincidence, I see consequence. When others see chance, I see cost.",
}

//go writeLines(ctx, chatLogUI, 10*time.Second)

// writeLines writes a line of text to the text widget every delay.
// Exits when the context expires.
func writeLines(ctx context.Context, t *text.Text, delay time.Duration) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			i := r.Intn(len(quotations))
			if err := t.Write(fmt.Sprintf("%s\n", quotations[i])); err != nil {
				panic(err)
			}

		case <-ctx.Done():
			return
		}
	}
}

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
func createChatLogUI(rc <-chan Result) *text.Text {
	rolled, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			for r := range rc {
				var err error
				for _, txt := range r.texts {
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
				cc <- Command{text: text}
			}
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}
	return input
}
