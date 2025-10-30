package main

import (
	"fmt"
	"os"
	"strconv"
	"the-mind/cards"
	"the-mind/players"
	"time"
)

func main() {
	// TODO: consider using flags
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

	playChan := make(chan int)
	done := make(chan struct{})

	// Create one player for demonstration
	// All players share the playChan which takes the card played (an int)
	fmt.Println("Starting a single player for demonstration...")
	player := &players.Player{
		Id:         1,
		Hand:       hands[0],
		PlayChan:   playChan,
		NotifyChan: make(chan struct{}),
	}
	// Start the player's play routine
	fmt.Printf("Player %d is starting to play...\n", player.Id)
	go player.Play()

	// Run a goroutine to listen for played cards
	// Gets cards from playChan and notifies the player to continue
	go func() {
		for card := range playChan {
			fmt.Printf("Main received card %d from Player %d\n", card, player.Id)

			player.NotifyChan <- struct{}{}
		}
		// Once all cards are played, signal done
		close(done)
	}()

	time.Sleep(1 * time.Second) // Give some time for the player to start
	// Wait for the game to finish
	close(playChan)
	<-done
	fmt.Println("Game over.")
}
