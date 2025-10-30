package players

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Player struct {
	Id              int
	Hand            []int
	PlayChan        chan Move
	CardsPlayedChan chan []int
}

type Move struct {
	PlayerId int
	Card     int
}

func CreatePlayers(numPlayers int, hands [][]int, playChan chan Move, cardsPlayedChan chan []int) []*Player {
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = &Player{
			Id:              i + 1,
			Hand:            hands[i],
			PlayChan:        playChan,
			CardsPlayedChan: cardsPlayedChan,
		}
	}
	return players
}

func (p *Player) Play(wg *sync.WaitGroup, done <-chan struct{}) {
	sort.Ints(p.Hand)
	fmt.Printf("Player %d's sorted hand: %v\n", p.Id, p.Hand)
	defer wg.Done()
	for len(p.Hand) > 0 {
		card := p.Hand[0]
		fmt.Printf("Player %d is preparing to play card %d\n", p.Id, card)
		delay := time.Duration(card) * time.Second / 200
		timer := time.NewTimer(delay)

		for {
			select {
			case <-timer.C:
				fmt.Printf("Player %d is playing card %d\n", p.Id, card)
				p.Hand = p.Hand[1:]
				p.PlayChan <- Move{PlayerId: p.Id, Card: card}

				select {
				case _, ok := <-p.CardsPlayedChan:
					if !ok {
						fmt.Printf("Player %d: CardsPlayedChan closed. Stopping play.\n", p.Id)
						return
					}
				case <-done:
					fmt.Printf("Player %d received done signal. Stopping.\n", p.Id)
					return
				}
			case played, ok := <-p.CardsPlayedChan:
				if !ok {
					fmt.Printf("Player %d: CardsPlayedChan closed. Stopping play.\n", p.Id)
					return
				}
				last := played[len(played)-1]
				fmt.Printf("Player %d received notification of played card %d while waiting to play %d\n", p.Id, last, card)
				diff := card - last
				if !timer.Stop() {
					<-timer.C
				}
				newDelay := time.Duration(diff) * 50 * time.Millisecond
				if diff > 20 {
					newDelay = 2 * newDelay
				} else if diff < 10 {
					newDelay = newDelay / 2
				} else if diff < 0 {
					newDelay = 0
				}

				fmt.Printf("Player %d adjusting timer for card %d to %v\n", p.Id, card, newDelay)
				timer.Reset(newDelay)
			case <-done:
				fmt.Printf("Player %d received done signal. Stopping play.\n", p.Id)
				if !timer.Stop() {
					<-timer.C
				}
				return
			}
			fmt.Printf("Player %d has no more cards to play.\n", p.Id)
		}
	}
}
