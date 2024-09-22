package models

type Message struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
}
type MessageHashCalculator interface {
	Calculate(data string) string
}
