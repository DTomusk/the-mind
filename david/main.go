package main

import (
	"fmt"
	"os"
	"strconv"
	"the-mind/cards"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <numPlayers> <cardsPerPlayer>")
		return
	}
	numPlayers, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error: numPlayers must be an integer.")
		return
	}
	cardsPerPlayer, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: cardsPerPlayer must be an integer.")
		return
	}

	if numPlayers*cardsPerPlayer > 100 {
		fmt.Println("Error: Not enough cards in the deck for the given number of players and cards per player.")
		return
	}

	hands := cards.GetHands(numPlayers, cardsPerPlayer)
	for i, hand := range hands {
		fmt.Printf("Player %d's hand: %v\n", i+1, hand)
	}
}
