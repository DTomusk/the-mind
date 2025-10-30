package players

import (
	"fmt"
	"math/rand"
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
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(p.Id)))
	defer wg.Done()
	for len(p.Hand) > 0 {
		time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)
		select {
		case <-done:
			fmt.Printf("Player %d received done signal. Stopping play.\n", p.Id)
			return
		default:
		}
		card := p.Hand[0]
		p.Hand = p.Hand[1:]
		fmt.Printf("Player %d is playing card %d\n", p.Id, card)
		p.PlayChan <- Move{PlayerId: p.Id, Card: card}
		<-p.CardsPlayedChan
	}
	fmt.Printf("Player %d has no more cards to play.\n", p.Id)
}
