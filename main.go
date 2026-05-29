package main

import (
	"fmt"
	"padelutils/padelset"
)

func main() {
	padelSet, _ := padelset.CreatePadelSetStandard()
	teamA, teamB, setCompleted := padelSet.ScoreForA()
	teamA, teamB, setCompleted = padelSet.ScoreForA()
	teamA, teamB, setCompleted = padelSet.ScoreForB()

	teamA, teamB, setCompleted = padelSet.ScoreForB()
	teamA, teamB, setCompleted = padelSet.ScoreForB()
	teamA, teamB, setCompleted = padelSet.ScoreForB()
	teamA, teamB, setCompleted = padelSet.ScoreForB()
	teamA, teamB, setCompleted = padelSet.ScoreForA()
	teamA, teamB, setCompleted = padelSet.ScoreForA()
	teamA, teamB, setCompleted = padelSet.ScoreForA()

	fmt.Printf("%v-%v (%v)\n", teamA, teamB, setCompleted)
	teamA, teamB, setCompleted = padelSet.ScoreForA()
	fmt.Printf("%v-%v (%v)\n", teamA, teamB, setCompleted)

	teamA, teamB, setCompleted = padelSet.ScoreForB()
	fmt.Printf("%v-%v (%v)\n", teamA, teamB, setCompleted)
	teamA, teamB, setCompleted = padelSet.ScoreForA()
	fmt.Printf("%v-%v (%v)\n", teamA, teamB, setCompleted)
}
