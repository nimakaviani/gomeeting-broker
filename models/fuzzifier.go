package models

import "github.com/sajari/fuzzy"

type Fuzzifier interface {
	Suggest(word string) string
}

type fuzzifier struct {
	model  *fuzzy.Model
	config Config
}

func NewFuzzifier(config Config) Fuzzifier {
	model := fuzzy.NewModel()
	f := fuzzifier{config: config, model: model}
	f.train()
	return f
}

func (f fuzzifier) train() {
	f.model.SetThreshold(1)
	f.model.SetDepth(5)

	var words = make([]string, len(f.config.Rooms))
	for _, room := range f.config.Rooms {
		words = append(words, room.Name)
	}
	f.model.Train(words)
}

func (f fuzzifier) Suggest(word string) string {
	if f.containsRoom(word) {
		return word
	}
	return f.model.Suggestions(word, false)[0]
}

func (f fuzzifier) containsRoom(roomName string) bool {
	for _, room := range f.config.Rooms {
		if room.Name == roomName {
			return true
		}
	}
	return false
}
