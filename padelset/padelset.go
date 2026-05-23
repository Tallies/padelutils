// Package padelset contains structures and functions that encapsulate the scoring for a padel set
package padelset

import (
	"fmt"
	"padelutils/internal/node"
	"padelutils/internal/padelbase"
	"strconv"
)

const (
	firstToKey = "FirstTo"
	bestOfKey  = "BestOf"
)

type PadelSetType int

const (
	Standard PadelSetType = iota // 6-7
	BestOf                       // Best of x, eg 7 (3-4, 5-2, 7,0)
	FirstTo                      // First to x, eg 7 (7-0, 7-3, 7-6)
)

type PadelSet struct {
	padelbase.PadelBase
	root     *node.Node
	score    *node.Node
	metadata map[string]any
	Type     PadelSetType
}

func (set *PadelSet) ScoreForA() (string, string, bool) {
	if !set.IsComplete() {
		set.score = set.score.Left
	}
	return set.GetScore()
}

func (set *PadelSet) ScoreForB() (string, string, bool) {
	if !set.IsComplete() {
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
	switch set.Type {
	case Standard:
		return isStandardComplete(set)
	case FirstTo:
		return isFirstToComplete(set, set.metadata[firstToKey].(int))
	case BestOf:
		return isBestOfComplete(set, set.metadata[bestOfKey].(int))
	}
	panic(fmt.Sprintf("Invalid padel set type '%v'", set.Type))
}

func isStandardComplete(set PadelSet) bool {
	// Helpers
	isSetScore := func(scoreOne int, scoreTwo int) bool {
		return (scoreOne == 6 && scoreTwo < 5) || (scoreOne == 7 && (scoreTwo == 5 || scoreTwo == 6))
	}

	// Implementation
	scoreAInt, scoreBInt := getScoresAsInt(set.score.ScoreA, set.score.ScoreB)
	return isSetScore(scoreAInt, scoreBInt) || isSetScore(scoreBInt, scoreAInt)
}

func isFirstToComplete(set PadelSet, firstTo int) bool {
	// Helpers
	isSetScore := func(scoreOne int, scoreTwo int) bool {
		return scoreOne == firstTo && scoreTwo < firstTo
	}

	// Implementation
	scoreAInt, scoreBInt := getScoresAsInt(set.score.ScoreA, set.score.ScoreB)
	return isSetScore(scoreAInt, scoreBInt) || isSetScore(scoreBInt, scoreAInt)
}

func isBestOfComplete(set PadelSet, bestOf int) bool {
	// Helpers
	isSetScore := func(scoreOne int, scoreTwo int) bool {
		return scoreOne+scoreTwo == bestOf
	}

	// Implementation
	scoreAInt, scoreBInt := getScoresAsInt(set.score.ScoreA, set.score.ScoreB)
	return isSetScore(scoreAInt, scoreBInt) || isSetScore(scoreBInt, scoreAInt)
}

func createNodes(scoreMax int, excludeFn func(int, int) bool) map[string]*node.Node {
	return node.CreateNodes(
		scoreMax,
		excludeFn,
	)
}

func getScoresAsInt(scoreA string, scoreB string) (int, int) {
	scoreAInt, errA := strconv.Atoi(scoreA)
	if errA != nil {
		panic(fmt.Sprintf("The value provided, %v, is not convertible to an int.", scoreA))
	}

	scoreBInt, errB := strconv.Atoi(scoreB)
	if errB != nil {
		panic(fmt.Sprintf("The value provided, %v, is not convertible to an int.", scoreB))
	}

	return scoreAInt, scoreBInt
}

func linkNodesStandardScoresTree(nodeCache map[string]*node.Node) *node.Node {
	// Helpers
	parentScoresFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt, scoreBInt := getScoresAsInt(scoreA, scoreB)
		if scoreAInt == 7 && scoreBInt >= 5 {
			return strconv.Itoa(scoreAInt - 1), scoreB, true
		} else if scoreAInt > 0 && scoreAInt <= 6 && scoreBInt <= 6 {
			return strconv.Itoa(scoreAInt - 1), scoreB, true
		}
		return "", "", false
	}

	parentScoresWithLeftChildFn := func(scoreA string, scoreB string) (string, string, bool) {
		parentA, parentB, ok := parentScoresFn(scoreA, scoreB)
		//fmt.Printf("%v-%v left child of %v-%v\n", scoreA, scoreB, parentA, parentB)
		return parentA, parentB, ok
	}

	parentScoresWithRightChildFn := func(scoreA string, scoreB string) (string, string, bool) {
		parentB, parentA, ok := parentScoresFn(scoreB, scoreA)
		//fmt.Printf("%v-%v right child of %v-%v\n", scoreA, scoreB, parentA, parentB)
		return parentA, parentB, ok
	}

	// Implementation
	return node.LinkNodes(nodeCache, parentScoresWithLeftChildFn, parentScoresWithRightChildFn)
}

func linkNodesFirstToScoresTree(nodeCache map[string]*node.Node, firstTo int) *node.Node {
	// Helpers
	parentScoresWithLeftChildFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt, scoreBInt := getScoresAsInt(scoreA, scoreB)
		if scoreAInt > 0 && scoreAInt <= firstTo && scoreBInt < firstTo {
			//fmt.Printf("Left: %v_%v - Parent: %v_%v\n", scoreA, scoreB, scoreAInt-1, scoreB)
			return strconv.Itoa(scoreAInt - 1), scoreB, true
		}
		return "", "", false
	}

	parentScoresWithRightChildFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt, scoreBInt := getScoresAsInt(scoreA, scoreB)
		if scoreBInt > 0 && scoreAInt < firstTo && scoreBInt <= firstTo {
			//fmt.Printf("Right: %v_%v - Parent: %v_%v\n", scoreA, scoreB, scoreA, scoreBInt-1)
			return scoreA, strconv.Itoa(scoreBInt - 1), true
		}
		return "", "", false
	}

	// Implementation
	return node.LinkNodes(nodeCache, parentScoresWithLeftChildFn, parentScoresWithRightChildFn)
}

func linkNodesBestOfScoresTree(nodeCache map[string]*node.Node, bestOf int) *node.Node {
	// Helpers
	parentScoresWithLeftChildFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt, scoreBInt := getScoresAsInt(scoreA, scoreB)
		if scoreAInt > 0 && scoreAInt+scoreBInt <= bestOf {
			return strconv.Itoa(scoreAInt - 1), scoreB, true
		}
		return "", "", false
	}

	parentScoresWithRightChildFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt, scoreBInt := getScoresAsInt(scoreA, scoreB)
		if scoreBInt > 0 && scoreAInt+scoreBInt <= bestOf {
			return scoreA, strconv.Itoa(scoreBInt - 1), true
		}
		return "", "", false
	}

	// Implementation
	return node.LinkNodes(nodeCache, parentScoresWithLeftChildFn, parentScoresWithRightChildFn)
}

func createPadelSetStandardScoresTree() *node.Node {
	nodeCache := createNodes(7,
		func(scoreA int, scoreB int) bool {
			return (scoreA == 7 && scoreB < 5) ||
				(scoreA < 5 && scoreB == 7) ||
				(scoreA == 7 && scoreB == 7)
		})
	return linkNodesStandardScoresTree(nodeCache)
}

func printAllNodes(padelSetType string, nodeCache map[string]*node.Node) {
	fmt.Println("Nodes for: ", padelSetType)
	for key := range nodeCache {
		fmt.Printf("%v\n", key)
	}
	fmt.Println("Done - Nodes for: ", padelSetType)
}

func createPadelSetFirstToScoresTree(firstTo int) *node.Node {
	nodeCache := createNodes(firstTo,
		// Works just like standard, except for special win by 2 rule.
		func(scoreA int, scoreB int) bool {
			return (scoreA == firstTo && scoreB == firstTo)
		})
	//printAllNodes(firstTo, nodeCache)
	return linkNodesFirstToScoresTree(nodeCache, firstTo)
}

func createPadelSetBestOfScoresTree(bestOf int) *node.Node {
	nodeCache := createNodes(bestOf,
		// Only sets where the score is less or equal to the bestOf
		func(scoreA int, scoreB int) bool {
			return scoreA+scoreB > bestOf
		})
	//printAllNodes(bestOf, nodeCache)
	return linkNodesBestOfScoresTree(nodeCache, bestOf)
}

// Public creation Methods

// CreatePadelSetStandard is the entry creation method for standard padel set of first to 6 games (win by 2), tie breaker at 6/6
func CreatePadelSetStandard() PadelSet {
	root := createPadelSetStandardScoresTree()
	return PadelSet{
		Type:  Standard,
		score: root,
		root:  root,
	}
}

// CreatePadelSetFirstTo is the entry creation method for a padel set that is a simple race to X number of games. Eg, first to 5 can be 5/0, 5/1, 5/2, 5/3, 5/4 for either team.
func CreatePadelSetFirstTo(firstTo int) PadelSet {
	root := createPadelSetFirstToScoresTree(firstTo)
	return PadelSet{
		Type:     FirstTo,
		score:    root,
		root:     root,
		metadata: map[string]any{firstToKey: firstTo},
	}
}

func CreatePadelSetBestOf(bestOf int) PadelSet {
	root := createPadelSetBestOfScoresTree(bestOf)
	return PadelSet{
		Type:     BestOf,
		score:    root,
		root:     root,
		metadata: map[string]any{bestOfKey: bestOf},
	}
}
