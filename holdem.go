package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const (
	kLowestRank = 7
	kAceRank = 14
	kNumRanks = (kAceRank - kLowestRank) + 1
	
	kNumSuits = 4
	kNumCards = kNumSuits * kNumRanks
	
	kNumCardsPerHand = 5
	
	kNumHands = 300000
)

type HandRank int
const (
	NoPair HandRank = iota
	Pair
	TwoPair
	SmallStraight
	ThreeOfAKind
	SmallFlush
	Straight
	FullHouse
	SmallStraightFlush
	FourOfAKind
	Flush
	StraightFlush
	
	NumHandRanks
)

var HandRankNames = [NumHandRanks]string {
	"NoPair",
	"Pair",
	"TwoPair",
	"SmallStraight",
	"ThreeOfAKind",
	"SmallFlush",
	"Straight",
	"FullHouse",
	"SmallStraightFlush",
	"FourOfAKind",
	"Flush",
	"StraightFlush",
}

type card struct {
	Suit	int
	Rank	int
	Index	int
}

func assert(valid bool, errorMsg string) {
	if !valid {
		panic(errorMsg)
	}
}

func Card(index int) (c card) {
	assert(0 <= index && index < kNumCards, "perm is out of range")
	
	c.Suit = index % kNumSuits
	c.Rank = index / kNumSuits
	c.Index = index
	
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

type ByIndex []card
func (c ByIndex) Len() int           { return len(c) }
func (c ByIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByIndex) Less(i, j int) bool { return c[i].Index < c[j].Index }

type Hand []card
func CalcHandRanks(hand Hand) []HandRank {
	handRanks := []HandRank{}

	// Check for all Pair / Triple / Quad-based hands
	numPairs := 0
	numTriples := 0
	numQuads := 0
	
	sort.Sort(ByRank(hand))
	
	lastRank := hand[0].Rank
	streak := 1
	for i := 1; i < kNumCardsPerHand; i++ {
		rank := hand[i].Rank
		
		if rank == lastRank {
			streak++
		} else {
			if streak == 4 {
				numQuads++
			} else if streak == 3 {
				numTriples++ 
			} else if streak == 2 {
				numPairs++
			}
			
			lastRank = rank
			streak = 1
		}
	}
	if streak == 4 {
		numQuads++
	} else if streak == 3 {
		numTriples++ 
	} else if streak == 2 {
		numPairs++
	}
	
	if numQuads == 1 {
		handRanks = append(handRanks, FourOfAKind)
	} else if numTriples == 1 && numPairs == 1 {
		handRanks = append(handRanks, FullHouse)
	} else if numTriples == 1 {
		handRanks = append(handRanks, ThreeOfAKind)
	} else if numPairs == 2 {
		handRanks = append(handRanks, TwoPair)
	} else if numPairs == 1 {
		handRanks = append(handRanks, Pair)
	} else {
		handRanks = append(handRanks, NoPair)
	}
	
	// Check for Straights
	if hand[0].Rank == hand[1].Rank-1 &&
	   hand[0].Rank == hand[2].Rank-2 &&
	   hand[0].Rank == hand[3].Rank-3 &&
	   hand[0].Rank == hand[4].Rank-4 {
		handRanks = append(handRanks, Straight)
	} else {
		for cardToSkip := 0; cardToSkip < 5; cardToSkip++ {
			i0 := 0
			i1 := 1
			i2 := 2
			i3 := 3
			if cardToSkip <= 0 { i0++ }
			if cardToSkip <= 1 { i1++ }
			if cardToSkip <= 2 { i2++ }
			if cardToSkip <= 3 { i3++ }
			
			if hand[i0].Rank == hand[i1].Rank-1 &&
			   hand[i0].Rank == hand[i2].Rank-2 &&
	           hand[i0].Rank == hand[i3].Rank-3 {
				handRanks = append(handRanks, SmallStraight)
				break
			}
		}
	}
	   		
	// Check for Flushes
	sort.Sort(BySuit(hand))
	
	if hand[0].Suit == hand[1].Suit &&
	   hand[0].Suit == hand[2].Suit &&
	   hand[0].Suit == hand[3].Suit &&
	   hand[0].Suit == hand[4].Suit {
		handRanks = append(handRanks, Flush)
	} else if hand[0].Suit == hand[1].Suit &&
			  hand[0].Suit == hand[2].Suit &&
			  hand[0].Suit == hand[3].Suit ||
			  hand[1].Suit == hand[2].Suit &&
			  hand[1].Suit == hand[3].Suit &&
			  hand[1].Suit == hand[4].Suit {
		handRanks = append(handRanks, SmallFlush)
	}
	
	// Check for Straight Flushes
	sort.Sort(ByIndex(hand))
	
	if hand[0].Suit == hand[1].Suit && hand[0].Rank == hand[1].Rank-1 &&
	   hand[0].Suit == hand[2].Suit && hand[0].Rank == hand[2].Rank-2 &&
	   hand[0].Suit == hand[3].Suit && hand[0].Rank == hand[3].Rank-3 &&
	   hand[0].Suit == hand[4].Suit && hand[0].Rank == hand[4].Rank-4 {
		handRanks = append(handRanks, StraightFlush)
	} else {
		for cardToSkip := 0; cardToSkip < 5; cardToSkip++ {
			i0 := 0
			i1 := 1
			i2 := 2
			i3 := 3
			if cardToSkip <= 0 { i0++ }
			if cardToSkip <= 1 { i1++ }
			if cardToSkip <= 2 { i2++ }
			if cardToSkip <= 3 { i3++ }

			if hand[i0].Suit == hand[i1].Suit && hand[i0].Rank == hand[i1].Rank-1 &&
		       hand[i0].Suit == hand[i2].Suit && hand[i0].Rank == hand[i2].Rank-2 &&
		       hand[i0].Suit == hand[i3].Suit && hand[i0].Rank == hand[i3].Rank-3 {
				handRanks = append(handRanks, SmallStraightFlush)
				break
			}
		}
	}
	
	return handRanks
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numMadeHands := make([]int, NumHandRanks)
	
	for i := 0; i < kNumHands; i++ {
		deck := rand.Perm(kNumCards)

//		fmt.Printf("New hand:\n")

		hand := make([]card, kNumCardsPerHand)
		
		for j := 0; j < kNumCardsPerHand; j++ {
			hand[j] = Card(deck[j])
		}
		
/*		for j := 0; j < kNumCardsPerHand; j++ {
			fmt.Printf("%v ", hand[j])
		}
		
		fmt.Printf("Sorted:\n")
		sort.Sort(ByRank(hand))
		for j := 0; j < kNumCardsPerHand; j++ {
			fmt.Printf("%v ", hand[j])
		}
*/		
		handRanks := CalcHandRanks(hand) 
		for _, hr := range(handRanks) {
//			fmt.Printf(" <-- %s ", HandRankNames[hr])
			
			numMadeHands[hr]++
		}

//		fmt.Print("\n\n")
	}
	
	handRankProbabilities := make(map[string]float32)
	
	for hr := NoPair; hr < NumHandRanks; hr++ {
		fmt.Printf("%s : %.3f%%\n", HandRankNames[hr], float32(numMadeHands[hr]) / float32(kNumHands) * 100.0)
		
		handRankProbabilities[HandRankNames[hr]] = float32(numMadeHands[hr]) / float32(kNumHands) * 100.0
	}
	
	
}