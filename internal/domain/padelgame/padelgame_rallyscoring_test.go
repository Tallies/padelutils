package padelgame

import (
	"fmt"
	"strconv"
	"testing"
)

// ---------------------------------------------------------------------------
// Initial state
// ---------------------------------------------------------------------------

func TestRallyScoringInitialScore(t *testing.T) {
	g := newGame(RallyScoring)
	a, b, done := g.GetScore()
	if a != "0" || b != "0" || done {
		t.Errorf("initial score: got (%s/%s done=%v), want (0/0 done=false)", a, b, done)
	}
	if g.IsComplete() {
		t.Error("new game should not be complete")
	}
}

func createSteps(scoreFor string, finalA int, finalB int) []scoreStep {
	countSub := 0
	scoreOther := "B"
	finalMain := finalA
	finalSub := finalB
	if scoreFor == "B" {
		scoreOther = "A"
		finalMain = finalB
		finalSub = finalA
	}

	steps := make([]scoreStep, 0)
	for count := 1; count <= finalMain; count++ {
		isComplete := count == finalMain
		if scoreFor == "A" {
			steps = append(steps, scoreStep{scoreFor, strconv.Itoa(count), strconv.Itoa(countSub), isComplete})
		} else {
			steps = append(steps, scoreStep{scoreFor, strconv.Itoa(countSub), strconv.Itoa(count), isComplete})
		}

		if(count <= finalSub) {
			steps = append(steps, scoreStep{scoreOther, strconv.Itoa(count), strconv.Itoa(count), isComplete})
			countSub++
		}
	}
	// fmt.Println(steps)
	return steps
}

// ---------------------------------------------------------------------------
// Team A wins without deuce (straight wins)
// ---------------------------------------------------------------------------

func TestRallyScoringTeamAWins_15_0(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("A", testRallyScore, 0))
}

func TestRallyScoringTeamAWins_15_5(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("A", testRallyScore, 5))
}

func TestRallyScoringTeamAWins_15_14(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("A", testRallyScore, testRallyScore - 1))
}

// ---------------------------------------------------------------------------
// Team B wins without deuce (straight wins)
// ---------------------------------------------------------------------------


func TestRallyScoringTeamBWins_15_0(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("B", 0, testRallyScore))
}

func TestRallyScoringTeamBWins_15_5(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("B", 5, testRallyScore))
}

func TestRallyScoringTeamBWins_15_14(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("B", testRallyScore - 1, testRallyScore))
}

// ---------------------------------------------------------------------------
// Score stays frozen after game is complete
// ---------------------------------------------------------------------------

func TestRallyScoringScoreFrozenAfterWin(t *testing.T) {
	g := newGame(RallyScoring)
	playSequence(t, &g, createSteps("A", testRallyScore, 1))
	
	a1, b1, done1 := g.GetScore()
	if !done1 {
		t.Fatal("game should be complete")
	}

	g.ScoreForA()
	g.ScoreForB()
	a2, b2, done2 := g.GetScore()
	if a1 != a2 || b1 != b2 || done1 != done2 {
		t.Errorf("score changed after game complete: (%s/%s) → (%s/%s)", a1, b1, a2, b2)
	}
}

// ---------------------------------------------------------------------------
// GetScore and IsComplete consistency
// ---------------------------------------------------------------------------

func TestRallyScoringGetScore_MatchesIsComplete(t *testing.T) {
	g := newGame(RallyScoring)
	steps := createSteps("A", testRallyScore - 1, 1)
	steps[len(steps)-1].wantDone = false
	playSequence(t, &g, steps)
	
	complete := g.IsComplete()
	_, _, done1 := g.GetScore()
	if complete != done1 {
		t.Errorf("GetScore done=%v bu IsComplet=%v", done1, complete)
	}

	g.ScoreForA()
	g.ScoreForA()
	g.ScoreForB()
		
	complete = g.IsComplete()
	_, _, done1 = g.GetScore()
	if done1 != complete {
		t.Errorf("GetScore done=%v but IsComplete=%v", done1, complete)
	}
}

// ---------------------------------------------------------------------------
// Rollback — ReverseScoreForA / ReverseScoreForB
// ---------------------------------------------------------------------------

// Single rollback from first point
func TestRallyScoringReverseScoreForA_SingleStep(t *testing.T) {
	g := newGame(RallyScoring)
	g.ScoreForA()                      // 0/0 → 1/0
	a, b, done := g.ReverseScoreForA() // back to 0/0
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForA: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

func TestRallyScoringReverseScoreForB_SingleStep(t *testing.T) {
	g := newGame(RallyScoring)
	g.ScoreForB()
	a, b, done := g.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForB: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback restores the previous node, not always 0/0
func TestRallyScoringReversePoint_MidGame(t *testing.T) {
	g := newGame(RallyScoring)
	g.ScoreForA() // 1/0
	g.ScoreForA() // 2/0
	g.ScoreForA() // 3/0

	a, b, _ := g.ReverseScoreForA() // back to 2/0
	if a != "2" || b != "0" {
		t.Errorf("rollback from 3/0: got %s/%s, want 2/0", a, b)
	}

	a, b, _ = g.ReverseScoreForA() // back to 15/0
	if a != "1" || b != "0" {
		t.Errorf("rollback from 2/0: got %s/%s, want 1/0", a, b)
	}
}

// Full rollback all the way to the start
func TestRallyScoringReversePoint_FullUnwind(t *testing.T) {
	g := newGame(RallyScoring)
	moves := []string{"A", "A", "B", "A"} // 1/0, 2/0, 2/1, 3,1
	for _, m := range moves {
		if m == "A" {
			g.ScoreForA()
		} else {
			g.ScoreForB()
		}
	}
	fmt.Printf("Before unwind: %v/%v\n", g.score.ScoreA, g.score.ScoreB)
	// Unwind all 4 moves
	for i := len(moves) - 1; i >= 0; i-- {
		if moves[i] == "A" {
			g.ReverseScoreForA()
		} else {
			g.ReverseScoreForB()
		}
		fmt.Printf("Reverse for %v: %v/%v\n", moves[i], g.score.ScoreA, g.score.ScoreB)
	}
	a, b, done := g.GetScore()
	if a != "0" || b != "0" || done {
		t.Errorf("after full unwind: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback at root (no history) is a no-op
func TestRallyScoringReversePoint_AtRoot_IsNoOp(t *testing.T) {
	g := newGame(RallyScoring)
	a, b, done := g.ReverseScoreForA()
	if a != "0" || b != "0" || done {
		t.Errorf("reverse at root: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
	a, b, done = g.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("reverse at root (B): got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback from a terminal win node restores the pre-win state
func TestRallyScoringReversePoint_FromWinNode(t *testing.T) {
	g := newGame(RallyScoring)
	steps := createSteps("A", testRallyScore, 5)
	playSequence(t, &g, steps)

	a, b, done := g.GetScore()
	if a != "15" || b != "5" || !done {
		t.Errorf("rollback from win: got (%s/%s done=%v), want (15/5 false)", a, b, done)
	}
	a, b, done = g.ReverseScoreForA()
	if a != "14" || b != "5" || done {
		t.Errorf("rollback from win: got (%s/%s done=%v), want (14/5 false)", a, b, done)
	}
	// Game is live again — should be able to continue
	g.ScoreForB()
	a, b, _ = g.GetScore()
	if a != "14" || b != "6" {
		t.Errorf("after resuming: got %s/%s, want 14/6", a, b)
	}
}

// Interleaving forward and reverse moves
func TestRallyScoringReversePoint_Interleaved(t *testing.T) {
	g := newGame(RallyScoring)
	g.ScoreForA()        // 1/0
	g.ScoreForA()        // 2/0
	g.ReverseScoreForA() // 1/0
	g.ScoreForB()        // 1/1
	g.ReverseScoreForB() // 1/0
	g.ScoreForA()        // 2/0

	a, b, _ := g.GetScore()
	if a != "2" || b != "0" {
		t.Errorf("interleaved result: got %s/%s, want 2/0", a, b)
	}
}

