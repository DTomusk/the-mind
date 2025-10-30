package main

import "fmt"

func main() {
	fmt.Println("Hello, David!")
	deck := createDeck()
	fmt.Println("Deck of cards:", deck)
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
