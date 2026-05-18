// Package padelsett contains structures and functions that encapsulate the scoring for ar a pa padel se set
package padelset

import (
	"errors"
	"fmt"
	"strconv"

	"padelutils/internal/node"
	"padelutils/internal/padelbase"
)

type PadelSetType int

const (
	Standard PadelSetType = iota // 6-7
	BestOf                       // Best of x, eg 7 (3-4, 5-2, 7,0)
	FirstTo                      // First to x, eg 7 (7-0, 7-3, 7-6)
)

type PadelSet struct {
	// Fields
	root  *node.Node
	score *node.Node
	Type  PadelSetType

	// Methods
	padelbase.PadelBase
}

func (set *PadelSet) ScoreForA() (string, string, bool) {
	if set.score.Left != nil {
		set.score = set.score.Left
	}
	return set.GetScore()
}

func (set *PadelSet) ScoreForB() (string, string, bool) {
	if set.score.Right != nil {
		set.score = set.score.Right
	}
	return set.GetScore()
}

func reduceScore(score string) string {
	scoreInt, err := strconv.Atoi(score)
	if err != nil {
		panic(fmt.Sprintf("'%v' is an invalid score for a set.", score))
	}
	return strconv.Itoa(scoreInt - 1)
}

func (set *PadelSet) ReverseScoreForA() (string, string, bool) {
	scoreA := set.score.ScoreA
	scoreB := set.score.ScoreB
	if scoreA == "0" {
		return set.GetScore()
	}
	reverseToNode := set.root.FindNode(reduceScore(scoreA), scoreB)
	set.score = reverseToNode
	return set.GetScore()
}

func (set *PadelSet) ReverseScoreForB() (string, string, bool) {
	scoreA := set.score.ScoreA
	scoreB := set.score.ScoreB
	if scoreB == "0" {
		return set.GetScore()
	}
	reverseToNode := set.root.FindNode(scoreA, reduceScore(scoreB))
	set.score = reverseToNode
	return set.GetScore()
}

func (set PadelSet) GetScore() (string, string, bool) {
	return set.score.ScoreA, set.score.ScoreB, set.IsComplete()
}

func (set PadelSet) IsComplete() bool {
	return set.score.IsLeaf()
}

func createNodes() map[string]*node.Node {
	return node.CreateNodes(
		7,
		func(scoreA int, scoreB int) bool {
			return (scoreA == 7 && scoreB < 5) ||
				(scoreA < 5 && scoreB == 7) ||
				(scoreA == 7 && scoreB == 7)
		},
	)
}

func linkNodes(nodeCache map[string]*node.Node) *node.Node {
	// Helpers
	parentScoresFn := func(scoreA string, scoreB string) (string, string, bool, error) {
		scoreAInt, errA := strconv.Atoi(scoreA)
		scoreBInt, errB := strconv.Atoi(scoreB)
		if errA != nil || errB != nil {
			return "", "", false, errors.New("Error converting score string to int.")
		}

		if (scoreAInt > 0 && scoreAInt <= 6 && scoreBInt <= 5) ||
			(scoreAInt == 6 && scoreBInt == 6) ||
			(scoreAInt == 7 && scoreBInt >= 5) {
			return strconv.Itoa(scoreAInt - 1), scoreB, true, nil
		}
		return "", "", false, nil
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
	return node.LinkNodes(nodeCache, leftChildOfFn, rightChildOfFn)
}

func CreatePadelSetScoresTree() *node.Node {
	nodeCache := createNodes()
	return linkNodes(nodeCache)
}

func CreatePadelSetStandard() PadelSet {
	root := CreatePadelSetScoresTree()
	return PadelSet{
		Type:  Standard,
		score: root,
		root:  root,
	}
}
