package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hello, David!")
	deck := createDeck()
	fmt.Println("Deck of cards:", deck)
	shuffleDeck(deck)
	fmt.Println("Shuffled deck of cards:", deck)
}

// initialize a new deck of cards numbered 1 to 100
func createDeck() []int {
	fmt.Println("Creating a new deck of cards...")
	deck := make([]int, 100)
	for i := 0; i < 100; i++ {
		deck[i] = i + 1
	}
	return deck
}

func shuffleDeck(deck []int) []int {
	// Placeholder for shuffle logic
	fmt.Println("Shuffling the deck...")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
