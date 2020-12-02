package main

type Message struct {
	texts []ColorText
}

type ColorText struct {
	text  string
	color string
}

func createMessage(texts ...string) Message {
	var _texts []ColorText
	for _, t := range texts {
		_texts = append(_texts, ColorText{text: t})
	}
	return Message{texts: _texts}
}
