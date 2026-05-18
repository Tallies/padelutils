package padelgame

import (
    "fmt"
	"padelutils/internal/node"
	"padelutils/internal/padelbase"
)

// Game Types Enum
type PadelGameType int
const (
    StarPoint PadelGameType = iota
    OneDeuce
    GoldenPoint
    )

type PadelGame struct {
    root *node.Node
    score *node.Node
    Type PadelGameType

    //Methods
    padelbase.PadelBase
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

func reduceScoreA(scoreA string, scoreB string) (string, string) {
    return reduceScore(scoreA, scoreB)    
}

func reduceScoreB(scoreA string, scoreB string) (string, string) {
    newB, newA := reduceScore(scoreB, scoreA)
    return newA, newB
}

func reduceScore(toReduce string, toRemain string) (string, string) {
    switch toReduce {
        case "0" : return toReduce, toRemain
        case "15" : return "0", toRemain
        case "30": return "15", toRemain
        case "40": return "30", toRemain
        case "40*": return "40", toRemain
        case "D1" : return "30", "40"
        case "A1" : return "D1", "D1"
        case "A1*": return "A1", toRemain
        case "D2" : return "40", "A1"
        case "A2": return "D2", "D2"
        case "A2*": return "A2", toRemain
        case "SP": return "40", "A2"
        case "SP*": return "SP", "SP"
    }
    panic(fmt.Sprintf("Cannot reduce %v. Unknown value.", toReduce))
}

func (game *PadelGame) ReverseScoreForA() (string, string, bool) {
    scoreA := game.score.ScoreA
    scoreB := game.score.ScoreB
    if scoreA == "0" {
        return game.GetScore()
    }
    reverseToNode := game.root.FindNode(reduceScoreA(scoreA, scoreB))
    game.score = reverseToNode
    return game.GetScore()
}

func (game *PadelGame) ReverseScoreForB() (string, string, bool) {
    scoreA := game.score.ScoreA
    scoreB := game.score.ScoreB
    if scoreB == "0" {
        return game.GetScore()
    }
    reverseToNode := game.root.FindNode(reduceScoreB(scoreA, scoreB))
    game.score = reverseToNode
    return game.GetScore()
}

func (game PadelGame) GetScore() (string, string, bool) {
    return game.score.ScoreA, game.score.ScoreB, game.IsComplete()
}

func (game PadelGame) IsComplete() bool {
    return game.score.Left == nil && game.score.Right == nil
}

//Score tree for diffent game types
func CreatePadelGameStarPointScoreTree() *node.Node {
    root := node.CreateNode("0", "0") 
    s15_0 := node.CreateNode("15", "0")
    s30_0 := node.CreateNode("30", "0")
    s40_0 := node.CreateNode("40", "0")
    s40_0_Win := node.CreateNode("40*", "0")
    s0_15 := node.CreateNode("0", "15")
    s0_30 := node.CreateNode("0","30")
    s0_40 := node.CreateNode("0", "40")
    s0_40_Win := node.CreateNode("0", "40*")
    
    s15_15 := node.CreateNode("15", "15")
    s30_15 := node.CreateNode("30", "15")
    s40_15 := node.CreateNode("40", "15")
    s40_15_Win := node.CreateNode("40*", "15")
    s15_30 := node.CreateNode("15", "30")
    s15_40 := node.CreateNode("15", "40")
    s15_40_Win := node.CreateNode("15", "40*")

    s30_30 := node.CreateNode("30", "30")
    s40_30 := node.CreateNode("40", "30")
    s40_30_Win := node.CreateNode("40*", "30")
    s30_40 := node.CreateNode("30", "40")
    s30_40_Win := node.CreateNode("30", "40*")
    
    sD1_D1 := node.CreateNode("D1", "D1")
    sA1_40 := node.CreateNode("A1", "40")
    sA1_40_Win := node.CreateNode("A1*", "40")
    s40_A1 := node.CreateNode("40", "A1")
    s40_A1_Win := node.CreateNode("40", "A1*")
    
    sD2_D2 := node.CreateNode("D2", "D2")
    sA2_40 := node.CreateNode("A2", "40")
    sA2_40_Win := node.CreateNode("A2*", "40")
    s40_A2 := node.CreateNode("40", "A2")
    s40_A2_Win := node.CreateNode("40", "A2*")
    
    sSP_SP := node.CreateNode("SP", "SP")
    sSP_40_Win := node.CreateNode("SP*", "40")
    s40_SP_Win := node.CreateNode("40", "SP*")
    
    root.Left = s15_0
    root.Right = s0_15
    
    s15_0.Left = s30_0
    s15_0.Right = s15_15
    
    s0_15.Left = s15_15
    s0_15.Right = s0_30
    
    s15_15.Left = s30_15
    s15_15.Right = s15_30
    
    s30_15.Left = s40_15
    s30_15.Right = s30_30
    
    s15_30.Left = s30_30
    s15_30.Right = s15_40
    
    s30_30.Left = s40_30
    s30_30.Right = s30_40
    
    s30_0.Left = s40_0
    s30_0.Right = s30_15
    
    s0_30.Left = s15_30
    s0_30.Right = s0_40
    
    s40_0.Left = s40_0_Win
    s40_0.Right = s40_15
    s40_15.Left = s40_15_Win
    s40_15.Right = s40_30
    
    s0_40.Left = s15_40
    s0_40.Right = s0_40_Win
    s15_40.Left = s30_40
    s15_40.Right = s15_40_Win
    
    s40_30.Left = s40_30_Win
    s40_30.Right = sD1_D1
    s30_40.Left = sD1_D1
    s30_40.Right = s30_40_Win
    sD1_D1.Left = sA1_40
    sD1_D1.Right = s40_A1
    
    sA1_40.Left = sA1_40_Win
    sA1_40.Right = sD2_D2
    s40_A1.Left = sD2_D2
    s40_A1.Right = s40_A1_Win
    sD2_D2.Left = sA2_40
    sD2_D2.Right = s40_A2
    
    sA2_40.Left = sA2_40_Win
    sA2_40.Right = sSP_SP
    s40_A2.Left = sSP_SP
    s40_A2.Right = s40_A2_Win
    sSP_SP.Left = sSP_40_Win
    sSP_SP.Right = s40_SP_Win
    
    return root
}

// PadelGame type creation functions
func CreatePadelGameStarPoint() PadelGame {
    root := CreatePadelGameStarPointScoreTree()
    return PadelGame {
        Type: StarPoint,
        score: root,
        root: root,
    }
}
