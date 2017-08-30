package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	kLowestRank = 6
	kAceRank = 14
	kNumRanks = (kAceRank - kLowestRank) + 1
	
	kNumSuits = 4
	kNumCards = kNumSuits * kNumRanks
)

type card struct {
	iSuit	int
	iRank	int
}

func assert(valid bool, errorMsg string) {
	if !valid {
		panic(errorMsg)
	}
}

func Card(perm int) (c card) {
	assert(0 <= perm && perm < kNumCards, "perm is out of range")
	
	c.iSuit = perm % kNumSuits
	c.iRank = perm / kNumSuits
	
	assert(0 <= c.iSuit && c.iSuit < kNumSuits, "suit is out of range")
	assert(0 <= c.iRank && c.iRank < kNumRanks, "rank is out of range")

	return c
}

func (c card) String() string {
	var rank, suit string

	r := c.iRank + kLowestRank
	switch r {
		case kAceRank:		rank = "A";
		case kAceRank - 1:	rank = "K";
		case kAceRank - 2:	rank = "Q";
		case kAceRank - 3:	rank = "J";
		default:			rank = strconv.Itoa(r);
	}
	
	switch c.iSuit {
		case 0: suit = "c";
		case 1: suit = "d";
		case 2: suit = "h";
		case 3: suit = "s";
		default: assert(false, "Invalid suit");
	}
	
	return rank + suit
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 10; i++ {
		deck := rand.Perm(kNumCards)

		//fmt.Printf("deck = %v\n", deck)

		hand := deck[0:5]

		for _, h := range hand {
			c := Card(h)
			fmt.Printf("%v ", c)
		}
		
		fmt.Print("\n")
	}
}