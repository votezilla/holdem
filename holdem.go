package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const (
	kLowestRank = 6
	kAceRank = 14
	kNumRanks = (kAceRank - kLowestRank) + 1
	
	kNumSuits = 4
	kNumCards = kNumSuits * kNumRanks
	
	kNumCardsPerHand = 5
)

type card struct {
	Suit	int
	Rank	int
}

func assert(valid bool, errorMsg string) {
	if !valid {
		panic(errorMsg)
	}
}

func Card(perm int) (c card) {
	assert(0 <= perm && perm < kNumCards, "perm is out of range")
	
	c.Suit = perm % kNumSuits
	c.Rank = perm / kNumSuits
	
	assert(0 <= c.Suit && c.Suit < kNumSuits, "suit is out of range")
	assert(0 <= c.Rank && c.Rank < kNumRanks, "rank is out of range")

	return c
}

func (c card) String() string {
	var rank, suit string

	r := c.Rank + kLowestRank
	switch r {
		case kAceRank:		rank = "A";
		case kAceRank - 1:	rank = "K";
		case kAceRank - 2:	rank = "Q";
		case kAceRank - 3:	rank = "J";
		default:			rank = strconv.Itoa(r);
	}
	
	switch c.Suit {
		case 0: suit = "c";
		case 1: suit = "d";
		case 2: suit = "h";
		case 3: suit = "s";
		default: assert(false, "Invalid suit");
	}
	
	return rank + suit
}

type ByRank []card
func (c ByRank) Len() int           { return len(c) }
func (c ByRank) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByRank) Less(i, j int) bool { return c[i].Rank < c[j].Rank }

type BySuit []card
func (c BySuit) Len() int           { return len(c) }
func (c BySuit) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c BySuit) Less(i, j int) bool { return c[i].Suit < c[j].Suit }

func main() {
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 10; i++ {
		deck := rand.Perm(kNumCards)

		fmt.Printf("New hand:\n")

		hand := make([]card, kNumCardsPerHand)
		for j := 0; j < kNumCardsPerHand; j++ {
			hand[j] = Card(deck[j])
			fmt.Printf("%v ", hand[j])
		}
		
		fmt.Printf("Sorted by rank:\n")
		sort.Sort(ByRank(hand))
		for _, h := range hand {
			fmt.Printf("%v ", h)
		}
		
		fmt.Printf("Sorted by suit:\n")
		sort.Sort(BySuit(hand))
		for _, h := range hand {
			fmt.Printf("%v ", h)
		}
		
		fmt.Print("\n\n")
	}
}