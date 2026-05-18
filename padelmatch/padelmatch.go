// Package padelmatch implements the structures and scoring logic to encapsulate scoring for a padel match
package padelmatch

import (
	"errors"
	"fmt"
	"strconv"

	"padelutils/internal/node"
	"padelutils/internal/padelbase"
)

// PadelMatchType is the kind match that the scoring will represent
type PadelMatchType int

const (
	Standard PadelMatchType = iota
	OneSet
)

// PadelMatch representss a single match set score.
type PadelMatch struct {
	// Fields
	root  *node.Node
	score *node.Node
	Type  PadelMatchType

	// Methods
	padelbase.PadelBase
}

func (match *PadelMatch) ScoreForA() (string, string, bool) {
	if match.score.Left != nil {
		match.score = match.score.Left
	}
	return match.GetScore()
}

func (match *PadelMatch) ScoreForB() (string, string, bool) {
	if match.score.Right != nil {
		match.score = match.score.Right
	}
	return match.GetScore()
}

func reduceScore(score string) string {
	scoreInt, err := strconv.Atoi(score)
	if err != nil {
		panic(fmt.Sprintf("'%v' is an invalid score for a set.", score))
	}
	return strconv.Itoa(scoreInt - 1)
}

func (match *PadelMatch) ReverseScoreForA() (string, string, bool) {
	scoreA := match.score.ScoreA
	scoreB := match.score.ScoreB
	if scoreA == "0" {
		return match.GetScore()
	}
	reverseToNode := match.root.FindNode(reduceScore(scoreA), scoreB)
	match.score = reverseToNode
	return match.GetScore()
}

func (match *PadelMatch) ReverseScoreForB() (string, string, bool) {
	scoreA := match.score.ScoreA
	scoreB := match.score.ScoreB
	if scoreB == "0" {
		return match.GetScore()
	}
	reverseToNode := match.root.FindNode(scoreA, reduceScore(scoreB))
	match.score = reverseToNode
	return match.GetScore()
}

func (match PadelMatch) GetScore() (string, string, bool) {
	return match.score.ScoreA, match.score.ScoreB, match.IsComplete()
}

func (match PadelMatch) IsComplete() bool {
	return match.score.IsLeaf()
}

func CreateStandardPadelMatchTree() *node.Node {
	// Helpers
	parentScoresFn := func(scoreA string, scoreB string) (string, string, bool, error) {
		scoreAInt, errA := strconv.Atoi(scoreA)
		if errA != nil {
			return "", "", false, errors.New("error converting score string to int")
		}

		return strconv.Itoa(scoreAInt - 1), scoreB, true, nil
	}

	leftChildOfFn := func(scoreA string, scoreB string) (string, string, bool) {
		newA, newB, ok, err := parentScoresFn(scoreA, scoreB)
		if err != nil {
			panic(fmt.Sprintf("Invalid scores to find parent on: %v-%v", scoreA, scoreB))
		}
		//fmt.Printf("%v-%v left child of %v-%v\n", scoreA, scoreB, newA, newB)
		return newA, newB, ok
	}

	rightChildOfFn := func(scoreA string, scoreB string) (string, string, bool) {
		newB, newA, ok, err := parentScoresFn(scoreB, scoreA)
		if err != nil {
			panic(fmt.Sprintf("Invalid scores to find parent on: %v-%v", scoreA, scoreB))
		}
		//fmt.Printf("%v-%v right child of %v-%v\n", scoreA, scoreB, newA, newB)
		return newA, newB, ok
	}

	// Implementation
	nodeCache := node.CreateNodes(2, func(scoreA int, scoreB int) bool { return false })
	root := node.LinkNodes(nodeCache, leftChildOfFn, rightChildOfFn)
	return root
}

func CreateStandardPadelMatch() PadelMatch {
	root := CreateStandardPadelMatchTree()
	return PadelMatch{
		Type:  Standard,
		root:  root,
		score: root,
	}
}
