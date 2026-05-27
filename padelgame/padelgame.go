package padelgame

import (
	"fmt"
	"strconv"

	"padelutils/internal/node"
	"padelutils/internal/padelbase"
)

const rallyScoreKey = "rallyScore"

// Game Types Enum
type PadelGameType int

const (
	Advantage PadelGameType = iota
	StarPoint
	OneDeuce
	GoldenPoint
	RallyScoring
)

type PadelGame struct {
	root     *node.Node
	score    *node.Node
	metadata map[string]any
	Type     PadelGameType

	// Methods
	padelbase.PadelBase
}

func scoreAsInt(score string) int {
	scoreInt, err := strconv.Atoi(score)
	if err != nil {
		panic(fmt.Sprintf("Score '%v' is invalid for RallyScoring", score))
	}
	return scoreInt
}

func (game *PadelGame) ScoreForA() (string, string, bool) {
	if game.score.Left != nil {
		game.score = game.score.Left
	}
	return game.GetScore()
}

func (game *PadelGame) ScoreForB() (string, string, bool) {
	if game.score.Right != nil {
		game.score = game.score.Right
	}
	return game.GetScore()
}

func reduceScoreA(scoreA string, scoreB string, gameType PadelGameType, metadata map[string]any) (string, string) {
	return reduceScore(scoreA, scoreB, gameType, metadata)
}

func reduceScoreB(scoreA string, scoreB string, gameType PadelGameType, metadata map[string]any) (string, string) {
	// Scoring is symmetric. Just swap A & B and then swap the return A & B
	newScoreB, newScoreA := reduceScore(scoreB, scoreA, gameType, metadata)
	return newScoreA, newScoreB
}

func reduceScore(toReduce string, toRemain string, gameType PadelGameType, metadata map[string]any) (string, string) {
	var newToReduce string
	var newToRemain string
	switch gameType {
	case Advantage:
		newToReduce, newToRemain = reduceScoreAdvantage(toReduce, toRemain)
	case StarPoint:
		newToReduce, newToRemain = reduceScoreStarPoint(toReduce, toRemain)
	case OneDeuce:
		newToReduce, newToRemain = reduceScoreOneDeuce(toReduce, toRemain)
	case GoldenPoint:
		newToReduce, newToRemain = reduceScoreGoldenPoint(toReduce, toRemain)
	case RallyScoring:
		_, ok := metadata[rallyScoreKey].(int)
		if !ok {
			panic("No upper score specified for rally scoring")
		}
		newToReduce, newToRemain = reduceScoreRallyScoring(toReduce, toRemain)
	default:
		newToReduce, newToRemain = toReduce, toRemain
	}
	return newToReduce, newToRemain
}

func reduceScoreStandard(toReduce string, toRemain string) (string, string, bool) {
	switch toReduce {
	case "0":
		return toReduce, toRemain, true
	case "15":
		return "0", toRemain, true
	case "30":
		return "15", toRemain, true
	case "40":
		return "30", toRemain, true
	case "40*":
		return "40", toRemain, true
	}
	return "", "", false
}

func reduceScoreAdvantageToDeuce(toReduce string, toRemain string) (string, string, bool) {
	switch {
	case toReduce == "40" && toRemain == "A1*":
		return "40", "A1", true
	case toReduce == "40" && toRemain == "A1":
		return "D1", "D1", true
	} 

	newToReduce, newToRemain, ok := reduceScoreStandard(toReduce, toRemain)
	if ok {
		return newToReduce, newToRemain, ok
	}

	switch toReduce {
	case "D1":
		return "30", "40", true
	case "A1":
		return "D1", "D1", true
	case "A1*":
		return "A1", "40", true
	}
	return "", "", false
}

func reduceScoreAdvantage(toReduce string, toRemain string) (string, string) {
	newToReduce, newToRemain, ok := reduceScoreAdvantageToDeuce(toReduce, toRemain)
	if ok {
		return newToReduce, newToRemain
	}
	panic(fmt.Sprintf("Invalid scores recieved to reduce for an Advantage game: '%v', '%v'", toReduce, toRemain))
}

func reduceScoreStarPoint(toReduce string, toRemain string) (string, string) {
	switch {
	case toReduce == "40" && toRemain == "SP*":
		return "SP", "SP"
	case toReduce == "40" && toRemain == "A2*":
		return "40", "A2"
	case toReduce == "40" && toRemain == "A2":
		return "D2", "D2"
	}

	newToReduce, newToRemain, ok := reduceScoreAdvantageToDeuce(toReduce, toRemain)
	if ok {
		return newToReduce, newToRemain
	}
	switch toReduce {
	case "D2":
		return "40", "A1"
	case "A2":
		return "D2", "D2"
	case "A2*":
		return "A2", "40"
	case "SP":
		return "40", "A2"
	case "SP*":
		return "SP", "SP"
	}
	panic(fmt.Sprintf("Invalid scores recieved to reduce for an StarPoint game: '%v', '%v'", toReduce, toRemain))
}

func reduceScoreOneDeuce(toReduce string, toRemain string) (string, string) {
	//Special case for when reducing the lower score from a win
	switch {
	case toReduce == "40" && toRemain == "GP*":
		return "GP", "GP"
	} 

	newToReduce, newToRemain, ok := reduceScoreAdvantageToDeuce(toReduce, toRemain)
	if ok {
		return newToReduce, newToRemain
	}
	
	switch toReduce {
	case "GP":
		return "40", "A1"
	case "GP*":
		return "GP", "GP"
	}
	
	panic(fmt.Sprintf("Invalid scores recieved to reduce for an OneDeuce game: '%v', '%v'", toReduce, toRemain))
}

func reduceScoreGoldenPoint(toReduce string, toRemain string) (string, string) {
	if toReduce == "40" && toRemain =="GP*"{
		return "GP", "GP"
	}

	newToReduce, newToRemain, ok := reduceScoreStandard(toReduce, toRemain)
	if ok {
		return newToReduce, newToRemain
	}

	switch toReduce {
	case "GP":
		return "30", "40"
	case "GP*":
		return "GP", "GP"
	}

	panic(fmt.Sprintf("Invalid scores recieved to reduce for an GoldenPoint game: '%v', '%v'", toReduce, toRemain))
}

func reduceScoreRallyScoring(toReduce string, toRemain string) (string, string) {
	if(toReduce == "0") {
		return toReduce, toRemain
	}
	toReduceInt := scoreAsInt(toReduce)
	return strconv.Itoa(toReduceInt -1), toRemain
}

func (game *PadelGame) ReverseScoreForA() (string, string, bool) {
	scoreA := game.score.ScoreA
	scoreB := game.score.ScoreB
	if scoreA == "0" {
		return game.GetScore()
	}
	newScoreA, newScoreB := reduceScoreA(scoreA, scoreB, game.Type, game.metadata)
	reverseToNode := game.root.FindNode(newScoreA, newScoreB)
	game.score = reverseToNode
	return game.GetScore()
}

func (game *PadelGame) ReverseScoreForB() (string, string, bool) {
	scoreA := game.score.ScoreA
	scoreB := game.score.ScoreB
	if scoreB == "0" {
		return game.GetScore()
	}
	newScoreA, newScoreB := reduceScoreB(scoreA, scoreB, game.Type, game.metadata)
	reverseToNode := game.root.FindNode(newScoreA, newScoreB)
	game.score = reverseToNode
	return game.GetScore()
}

func (game PadelGame) GetScore() (string, string, bool) {
	return game.score.ScoreA, game.score.ScoreB, game.IsComplete()
}

func (game PadelGame) IsComplete() bool {
	if game.Type == RallyScoring {
		rallyScore := strconv.Itoa(game.metadata[rallyScoreKey].(int))
		return game.score.ScoreA == rallyScore || game.score.ScoreB == rallyScore
	}
	return game.score.Left == nil && game.score.Right == nil
}

// Score tree for diffent game types
func createStandardScoreTree() map[string]*node.Node {
	nodeCache := make(map[string]*node.Node)

	nodeCache[a0b0] = node.CreateNode("0", "0")

	nodeCache[a15b0] = node.CreateNode("15", "0")
	nodeCache[a30b0] = node.CreateNode("30", "0")
	nodeCache[a40b0] = node.CreateNode("40", "0")
	nodeCache[a40wb0] = node.CreateNode("40*", "0")

	nodeCache[a0b15] = node.CreateNode("0", "15")
	nodeCache[a0b30] = node.CreateNode("0", "30")
	nodeCache[a0b40] = node.CreateNode("0", "40")
	nodeCache[a0b40w] = node.CreateNode("0", "40*")

	nodeCache[a15b15] = node.CreateNode("15", "15")
	nodeCache[a30b15] = node.CreateNode("30", "15")
	nodeCache[a40b15] = node.CreateNode("40", "15")
	nodeCache[a40wb15] = node.CreateNode("40*", "15")

	nodeCache[a15b30] = node.CreateNode("15", "30")
	nodeCache[a15b40] = node.CreateNode("15", "40")
	nodeCache[a15b40w] = node.CreateNode("15", "40*")

	nodeCache[a30b30] = node.CreateNode("30", "30")
	nodeCache[a40b30] = node.CreateNode("40", "30")
	nodeCache[a40wb30] = node.CreateNode("40*", "30")

	nodeCache[a30b40] = node.CreateNode("30", "40")
	nodeCache[a30b40w] = node.CreateNode("30", "40*")

	return nodeCache
}

func createDeuceToAdvantageScoreTree() map[string]*node.Node {
	nodeCache := createStandardScoreTree()

	nodeCache[aD1bD1] = node.CreateNode("D1", "D1")
	nodeCache[aA1b40] = node.CreateNode("A1", "40")
	nodeCache[aA1wb40] = node.CreateNode("A1*", "40")
	nodeCache[a40bA1] = node.CreateNode("40", "A1")
	nodeCache[a40bA1w] = node.CreateNode("40", "A1*")

	return nodeCache
}

func createAdvantageScoreTree() map[string]*node.Node {
	nodeCache := createDeuceToAdvantageScoreTree()

	return nodeCache
}

func createStarPointScoreTree() map[string]*node.Node {
	nodeCache := createDeuceToAdvantageScoreTree()

	nodeCache[aD2bD2] = node.CreateNode("D2", "D2")
	nodeCache[aA2b40] = node.CreateNode("A2", "40")
	nodeCache[aA2wb40] = node.CreateNode("A2*", "40")
	nodeCache[a40bA2] = node.CreateNode("40", "A2")
	nodeCache[a40bA2w] = node.CreateNode("40", "A2*")

	nodeCache[aSPbSP] = node.CreateNode("SP", "SP")
	nodeCache[aSPwb40] = node.CreateNode("SP*", "40")
	nodeCache[a40bSPw] = node.CreateNode("40", "SP*")

	return nodeCache
}

func createOneDeuceScoreTree() map[string]*node.Node {
	nodeCache := createDeuceToAdvantageScoreTree()

	nodeCache[aGPbGP] = node.CreateNode("GP", "GP")
	nodeCache[aGPwb40] = node.CreateNode("GP*", "40")
	nodeCache[a40bGPw] = node.CreateNode("40", "GP*")

	return nodeCache
}

func createGoldenPointScoreTree() map[string]*node.Node {
	nodeCache := createStandardScoreTree()

	nodeCache[aGPbGP] = node.CreateNode("GP", "GP")
	nodeCache[aGPwb40] = node.CreateNode("GP*", "40")
	nodeCache[a40bGPw] = node.CreateNode("40", "GP*")

	return nodeCache
}

func createRallyScoringScoreTree(rallyScore int) map[string]*node.Node {
	nodeCache := node.CreateNodes(rallyScore,
		func(scoreA int, scoreB int) bool { return scoreA == rallyScore && scoreB == rallyScore })
	return nodeCache
}

func linkStandardScoreTree(nodeCache map[string]*node.Node) *node.Node {
	root := nodeCache[a0b0]

	root.Left = nodeCache[a15b0]
	root.Right = nodeCache[a0b15]

	nodeCache[a15b0].Left = nodeCache[a30b0]
	nodeCache[a15b0].Right = nodeCache[a15b15]

	nodeCache[a0b15].Left = nodeCache[a15b15]
	nodeCache[a0b15].Right = nodeCache[a0b30]

	nodeCache[a15b15].Left = nodeCache[a30b15]
	nodeCache[a15b15].Right = nodeCache[a15b30]

	nodeCache[a30b15].Left = nodeCache[a40b15]
	nodeCache[a30b15].Right = nodeCache[a30b30]

	nodeCache[a15b30].Left = nodeCache[a30b30]
	nodeCache[a15b30].Right = nodeCache[a15b40]

	nodeCache[a30b30].Left = nodeCache[a40b30]
	nodeCache[a30b30].Right = nodeCache[a30b40]

	nodeCache[a30b0].Left = nodeCache[a40b0]
	nodeCache[a30b0].Right = nodeCache[a30b15]

	nodeCache[a0b30].Left = nodeCache[a15b30]
	nodeCache[a0b30].Right = nodeCache[a0b40]

	nodeCache[a40b0].Left = nodeCache[a40wb0]
	nodeCache[a40b0].Right = nodeCache[a40b15]

	nodeCache[a40b15].Left = nodeCache[a40wb15]
	nodeCache[a40b15].Right = nodeCache[a40b30]

	nodeCache[a0b40].Left = nodeCache[a15b40]
	nodeCache[a0b40].Right = nodeCache[a0b40w]

	nodeCache[a15b40].Left = nodeCache[a30b40]
	nodeCache[a15b40].Right = nodeCache[a15b40w]

	nodeCache[a40b30].Left = nodeCache[a40wb30]
	nodeCache[a30b40].Right = nodeCache[a30b40w]

	return nodeCache[a0b0]
}

func linkDeuceToAdvantageScoreTree(nodeCache map[string]*node.Node) *node.Node {
	root := linkStandardScoreTree(nodeCache)

	nodeCache[a40b30].Right = nodeCache[aD1bD1]
	nodeCache[a30b40].Left = nodeCache[aD1bD1]

	nodeCache[aD1bD1].Left = nodeCache[aA1b40]
	nodeCache[aD1bD1].Right = nodeCache[a40bA1]

	nodeCache[aA1b40].Left = nodeCache[aA1wb40]
	nodeCache[a40bA1].Right = nodeCache[a40bA1w]

	return root
}

func linkAdvantageScoreTree(nodeCache map[string]*node.Node) *node.Node {
	root := linkDeuceToAdvantageScoreTree(nodeCache)

	// Loop advantage back to duece
	nodeCache[aA1b40].Right = nodeCache[aD1bD1]
	nodeCache[a40bA1].Left = nodeCache[aD1bD1]

	return root
}

func linkStarPointScoreTree(nodeCache map[string]*node.Node) *node.Node {
	root := linkDeuceToAdvantageScoreTree(nodeCache)

	nodeCache[aA1b40].Right = nodeCache[aD2bD2]
	nodeCache[a40bA1].Left = nodeCache[aD2bD2]

	nodeCache[aD2bD2].Left = nodeCache[aA2b40]
	nodeCache[aD2bD2].Right = nodeCache[a40bA2]

	nodeCache[aA2b40].Left = nodeCache[aA2wb40]
	nodeCache[aA2b40].Right = nodeCache[aSPbSP]
  
	nodeCache[a40bA2].Left = nodeCache[aSPbSP]
	nodeCache[a40bA2].Right = nodeCache[a40bA2w]

	nodeCache[aSPbSP].Left = nodeCache[aSPwb40]
	nodeCache[aSPbSP].Right = nodeCache[a40bSPw]

	return root
}

func linkOneDeuceScoreTree(nodeCache map[string]*node.Node) *node.Node {
	root := linkDeuceToAdvantageScoreTree(nodeCache)

	nodeCache[aA1b40].Right = nodeCache[aGPbGP]
	nodeCache[a40bA1].Left = nodeCache[aGPbGP]
	nodeCache[aGPbGP].Left = nodeCache[aGPwb40]
	nodeCache[aGPbGP].Right = nodeCache[a40bGPw]
	return root
}

func linkGoldenPointScoreTree(nodeCache map[string]*node.Node) *node.Node {
	root := linkStandardScoreTree(nodeCache)

	nodeCache[a40b30].Right = nodeCache[aGPbGP]
	nodeCache[a30b40].Left = nodeCache[aGPbGP]
	nodeCache[aGPbGP].Left = nodeCache[aGPwb40]
	nodeCache[aGPbGP].Right = nodeCache[a40bGPw]

	return root
}

func linkRallyScoringScoreTree(nodeCache map[string]*node.Node, rallyScore int) *node.Node {
	// Helpers
	leftChildOfFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt := scoreAsInt(scoreA)
		scoreBInt := scoreAsInt(scoreB)

		if scoreAInt > 0 && scoreBInt < rallyScore {
			return strconv.Itoa(scoreAInt - 1), scoreB, true
		}

		return "", "", false
	}

	rightChildOfFn := func(scoreA string, scoreB string) (string, string, bool) {
		scoreAInt := scoreAsInt(scoreA)
		scoreBInt := scoreAsInt(scoreB)
		if scoreBInt > 0 && scoreAInt < rallyScore {
			return scoreA, strconv.Itoa(scoreBInt - 1), true
		}
		return "", "", false
	}

	// Implementation
	root := node.LinkNodes(nodeCache, leftChildOfFn, rightChildOfFn)
	return root
}

// PadelGame type creation functions
func CreatePadelGameAdvantage() PadelGame {
	nodeCache := createAdvantageScoreTree()
	root := linkAdvantageScoreTree(nodeCache)
	return PadelGame{
		Type:     StarPoint,
		score:    root,
		root:     root,
		metadata: make(map[string]any),
	}
}

func CreatePadelGameStarPoint() PadelGame {
	nodeCache := createStarPointScoreTree()
	root := linkStarPointScoreTree(nodeCache)
	return PadelGame{
		Type:     StarPoint,
		score:    root,
		root:     root,
		metadata: make(map[string]any),
	}
}

func CreatePadelGameOneDeuce() PadelGame {
	nodeCache := createOneDeuceScoreTree()
	root := linkOneDeuceScoreTree(nodeCache)
	return PadelGame{
		Type:     OneDeuce,
		score:    root,
		root:     root,
		metadata: make(map[string]any),
	}
}

func CreatePadelGameGoldenPoint() PadelGame {
	nodeCache := createGoldenPointScoreTree()
	root := linkGoldenPointScoreTree(nodeCache)
	return PadelGame{
		Type:     GoldenPoint,
		score:    root,
		root:     root,
		metadata: make(map[string]any),
	}
}

func CreatePadelGameRallyScoring(rallyScore int) PadelGame {
	nodeCache := createRallyScoringScoreTree(rallyScore)
	root := linkRallyScoringScoreTree(nodeCache, rallyScore)
	return PadelGame{
		Type:     RallyScoring,
		score:    root,
		root:     root,
		metadata: map[string]any{rallyScoreKey: rallyScore},
	}
}
