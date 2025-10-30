package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"the-mind/cards"
	"the-mind/players"
)

type GameConfig struct {
	NumPlayers     int
	CardsPerPlayer int
}

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	playChan := make(chan int)
	cardsPlayedChan := make(chan []int)
	done := make(chan struct{})
	var wg sync.WaitGroup

	players := setUpGame(config, playChan, cardsPlayedChan)
	player := players[0]

	// Start the player's play routine in the background
	fmt.Printf("Player %d is starting to play...\n", player.Id)
	wg.Add(1)
	go player.Play(&wg, done)

	// Run a goroutine in the background to listen for played cards
	// Gets cards from playChan and notifies the player to continue
	var gameWg sync.WaitGroup
	var cardsPlayed []int
	gameWg.Add(1)
	go func() {
		defer gameWg.Done()
		for card := range playChan {
			fmt.Printf("Main received card %d from Player %d\n", card, player.Id)

			// Send an empty struct to notify the player that a card has been played
			if len(cardsPlayed) > 0 {
				last := cardsPlayed[len(cardsPlayed)-1]
				if card < last {
					fmt.Printf("Card %d played after %d! Game over.\n", card, last)
					player.CardsPlayedChan <- cardsPlayed
					close(done)
					return
				}
			}
			cardsPlayed = append(cardsPlayed, card)
			fmt.Printf("Cards played so far: %v\n", cardsPlayed)
			player.CardsPlayedChan <- cardsPlayed
		}
		// Once playChan is closed, we finish
		fmt.Println("No more cards to receive. Ending game.")
	}()

	// Wait for the player to finish, each player adds 1 to the wait group
	wg.Wait()
	// Once all players are done, close playChan to end the game goroutine
	close(playChan)
	gameWg.Wait()
	select {
	case <-done:
		fmt.Println("Game ended early due to invalid card.")
	default:
		fmt.Println("Game over — all cards played successfully.")
	}
}

func parseArgs() (*GameConfig, error) {
	// TODO: consider using flags
	args := os.Args[1:]
	if len(args) != 2 {
		return nil, fmt.Errorf("usage: go run main.go <numPlayers> <cardsPerPlayer>")

	}
	numPlayers, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("error: numPlayers must be an integer")
	}
	cardsPerPlayer, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("error: cardsPerPlayer must be an integer")
	}

	if numPlayers*cardsPerPlayer > 100 {
		return nil, fmt.Errorf("error: Not enough cards in the deck for the given number of players and cards per player")
	}
	return &GameConfig{
		NumPlayers:     numPlayers,
		CardsPerPlayer: cardsPerPlayer,
	}, nil
}

func setUpGame(config *GameConfig, playChan chan int, cardsPlayedChan chan []int) []*players.Player {
	hands := cards.GetHands(config.NumPlayers, config.CardsPerPlayer)
	for i, hand := range hands {
		fmt.Printf("Player %d's hand: %v\n", i+1, hand)
	}
	players := players.CreatePlayers(config.NumPlayers, hands, playChan, cardsPlayedChan)
	return players
}
