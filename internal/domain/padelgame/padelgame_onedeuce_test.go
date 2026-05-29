package padelgame

import (
	"fmt"
	"testing"
)

// ---------------------------------------------------------------------------
// Initial state
// ---------------------------------------------------------------------------

func TestOneDeuceInitialScore(t *testing.T) {
	g := newGame(OneDeuce)
	a, b, done := g.GetScore()
	if a != "0" || b != "0" || done {
		t.Errorf("initial score: got (%s/%s done=%v), want (0/0 done=false)", a, b, done)
	}
	if g.IsComplete() {
		t.Error("new game should not be complete")
	}
}

// ---------------------------------------------------------------------------
// Team A wins without deuce (straight wins)
// ---------------------------------------------------------------------------

func TestOneDeuceTeamAWins_40_0(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"A", "40*", "0", true}, // terminal win node
	})
}

func TestOneDeuceTeamAWins_40_15(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"B", "0", "15", false},
		{"A", "15", "15", false},
		{"A", "30", "15", false},
		{"A", "40", "15", false},
		{"A", "40*", "15", true},
	})
}

func TestOneDeuceTeamAWins_40_30(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"B", "0", "15", false},
		{"B", "0", "30", false},
		{"A", "15", "30", false},
		{"A", "30", "30", false},
		{"A", "40", "30", false},
		{"A", "40*", "30", true},
	})
}

// ---------------------------------------------------------------------------
// Team B wins without deuce (straight wins)
// ---------------------------------------------------------------------------

func TestOneDeuceTeamBWins_0_40(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"B", "0", "15", false},
		{"B", "0", "30", false},
		{"B", "0", "40", false},
		{"B", "0", "40*", true},
	})
}

func TestOneDeuceTeamBWins_15_40(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"B", "15", "15", false},
		{"B", "15", "30", false},
		{"B", "15", "40", false},
		{"B", "15", "40*", true},
	})
}

func TestOneDeuceTeamBWins_30_40(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"B", "30", "15", false},
		{"B", "30", "30", false},
		{"B", "30", "40", false},
		{"B", "30", "40*", true},
	})
}

// ---------------------------------------------------------------------------
// Deuce 1 (D1/D1) paths
// ---------------------------------------------------------------------------

// 40-30 → B scores → D1/D1 → A scores → Advantage A (A1/40) → A wins
func TestOneDeuceDeuce1_AAdvantage_AWins(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"B", "40", "15", false},
		{"B", "40", "30", false},
		{"B", "D1", "D1", false},
		{"A", "A1", "40", false},
		{"A", "A1*", "40", true},
	})
}

// 40-30 → B scores → D1/D1 → B scores → Advantage B (40/A1) → B wins
func TestOneDeuceDeuce1_BAdvantage_BWins(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"B", "40", "15", false},
		{"B", "40", "30", false},
		{"B", "D1", "D1", false},
		{"B", "40", "A1", false},
		{"B", "40", "A1*", true},
	})
}

// Arrive at D1 via 30-40 path
func TestOneDeuceDeuce1_Via30_40(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"B", "0", "15", false},
		{"B", "0", "30", false},
		{"B", "0", "40", false},
		{"A", "15", "40", false},
		{"A", "30", "40", false},
		{"A", "D1", "D1", false}, // arrives from 30-40 side
		{"A", "A1", "40", false},
		{"A", "A1*", "40", true},
	})
}

// ---------------------------------------------------------------------------
// Deuce 2 (D2/D2) paths
// ---------------------------------------------------------------------------

// D1 → A adv → B scores → D2 → A wins
func TestOneDeuceGoldenPoint_AWins(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"B", "40", "15", false},
		{"B", "40", "30", false},
		{"B", "D1", "D1", false},
		{"A", "A1", "40", false},
		{"B", "GP", "GP", false},
		{"A", "GP*", "40", true},
	})
}

// D1 → B adv → A scores → D2 → B wins
func TestOneDeuceGoldenPoint_BWins(t *testing.T) {
	g := newGame(OneDeuce)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"B", "40", "15", false},
		{"B", "40", "30", false},
		{"B", "D1", "D1", false},
		{"B", "40", "A1", false},
		{"A", "GP", "GP", false},
		{"B", "40", "GP*", true},
	})
}

// ---------------------------------------------------------------------------
// Score stays frozen after game is complete
// ---------------------------------------------------------------------------

func TestOneDeuceScoreFrozenAfterWin(t *testing.T) {
	g := newGame(OneDeuce)
	// A wins 40-0
	for range 4 {
		g.ScoreForA()
	}
	a1, b1, done1 := g.GetScore()
	if !done1 {
		t.Fatal("game should be complete")
	}
	// Extra points must not change score
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

func TestOneDeuceGetScore_MatchesIsComplete(t *testing.T) {
	g := newGame(OneDeuce)
	steps := []string{"A", "A", "A", "B", "B", "B", "A", "A"} // path to D1 then A wins
	for i, who := range steps {
		var _, _, done bool
		if who == "A" {
			_, _, done = g.ScoreForA()
		} else {
			_, _, done = g.ScoreForB()
		}
		complete := g.IsComplete()
		if done != complete {
			t.Errorf("step %d: GetScore done=%v but IsComplete=%v", i, done, complete)
		}
	}
}

// ---------------------------------------------------------------------------
// Rollback — ReverseScoreForA / ReverseScoreForB
// ---------------------------------------------------------------------------

// Single rollback from first point
func TestOneDeuceReverseScoreForA_SingleStep(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForA()                      // 0/0 → 15/0
	a, b, done := g.ReverseScoreForA() // back to 0/0
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForA: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

func TestOneDeuceReverseScoreForB_SingleStep(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForB()
	a, b, done := g.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForB: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback restores the previous node, not always 0/0
func TestOneDeuceReversePoint_MidGame(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForA() // 15/0
	g.ScoreForA() // 30/0
	g.ScoreForA() // 40/0

	a, b, _ := g.ReverseScoreForA() // back to 30/0
	if a != "30" || b != "0" {
		t.Errorf("rollback from 40/0: got %s/%s, want 30/0", a, b)
	}

	a, b, _ = g.ReverseScoreForA() // back to 15/0
	if a != "15" || b != "0" {
		t.Errorf("rollback from 30/0: got %s/%s, want 15/0", a, b)
	}
}

// Full rollback all the way to the start
func TestOneDeuceReversePoint_FullUnwind(t *testing.T) {
	g := newGame(OneDeuce)
	moves := []string{"A", "A", "B", "A"} // 15/0 -> 30/0 -> 30/15 -> 40/15
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
func TestOneDeuceReversePoint_AtRoot_IsNoOp(t *testing.T) {
	g := newGame(OneDeuce)
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
func TestOneDeuceReversePoint_FromWinNode(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForA() // 15/0
	g.ScoreForA() // 30/0
	g.ScoreForA() // 40/0
	g.ScoreForA() // win (40/0 terminal)

	a, b, done := g.ReverseScoreForA()
	if a != "40" || b != "0" || done {
		t.Errorf("rollback from win: got (%s/%s done=%v), want (40/0 false)", a, b, done)
	}
	// Game is live again — should be able to continue
	g.ScoreForB() // 40/15 (not terminal)
	a, b, _ = g.GetScore()
	if a != "40" || b != "15" {
		t.Errorf("after resuming: got %s/%s, want 40/15", a, b)
	}
}

// Rollback through Deuce 1
func TestOneDeuceReversePoint_ThroughDeuce1(t *testing.T) {
	g := newGame(OneDeuce)
	// Advance to D1
	g.ScoreForA()
	g.ScoreForA()
	g.ScoreForA() // 40/0
	g.ScoreForB()
	g.ScoreForB() // 40/30
	g.ScoreForB() // D1/D1

	a, b, _ := g.GetScore()
	if a != "D1" || b != "D1" {
		t.Fatalf("expected D1/D1, got %s/%s", a, b)
	}

	g.ReverseScoreForB() // back to 40/30
	a, b, _ = g.GetScore()
	if a != "40" || b != "30" {
		t.Errorf("rollback from D1: got %s/%s, want 40/30", a, b)
	}
}

// Rollback through Star Point
func TestOneDeuceReversePoint_ThroughGoldenPoint(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForA()
	g.ScoreForA()
	g.ScoreForA() // 40/0
	g.ScoreForB()
	g.ScoreForB() // 40/30
	g.ScoreForB() // D1/D1
	g.ScoreForA() // A1/40
	g.ScoreForB() // GP/GP
	g.ScoreForA() // GP*/40

	a, b, _ := g.GetScore()
	if a != "GP*" || b != "40" {
		t.Fatalf("expected GP*/40, got %s/%s", a, b)
	}

	g.ReverseScoreForA() // back to GP/40 -> GP/GP
	a, b, _ = g.GetScore()
	if a != "GP" || b != "GP" {
		t.Errorf("rollback A from GP*/40: got %s/%s, want GP/GP", a, b)
	}
}

func TestOneDeuceReversePoint_ThroughGoldenPointFromLowerSide(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForA()
	g.ScoreForA()
	g.ScoreForA() // 40/0
	g.ScoreForB()
	g.ScoreForB() // 40/30
	g.ScoreForB() // D1/D1
	g.ScoreForA() // A1/40
	g.ScoreForB() // GP/GP
	g.ScoreForB() // 40/GP*

	a, b, _ := g.GetScore()
	if a != "40" || b != "GP*" {
		t.Fatalf("expected 40/GP*, got %s/%s", a, b)
	}

	g.ReverseScoreForA() // back to GP/40 -> 40/30
	a, b, _ = g.GetScore()
	if a != "GP" || b != "GP" {
		t.Errorf("rollback A from 40/GP*: got %s/%s, want GP/GP", a, b)
	}

	g.ReverseScoreForA() // back to GP/GP -> 40/A1
	a, b, _ = g.GetScore()
	if a != "40" || b != "A1" {
		t.Errorf("rollback A from GP/GP: got %s/%s, want 40/A1", a, b)
	}

	g.ScoreForB() //Win 40/A1*
	g.ReverseScoreForA() // back to GP/GP -> 40/A1
	a, b, _ = g.GetScore()
	if a != "40" || b != "A1" {
		t.Errorf("rollback A from 40/A1*: got %s/%s, want 40/A1", a, b)
	}

	g.ReverseScoreForA() // back to 40/A1 -> D1/D1
	a, b, _ = g.GetScore()
	if a != "D1" || b != "D1" {
		t.Errorf("rollback A from 40/A1: got %s/%s, want D1/D1", a, b)
	}
}

// Interleaving forward and reverse moves
func TestOneDeuceReversePoint_Interleaved(t *testing.T) {
	g := newGame(OneDeuce)
	g.ScoreForA()        // 15/0
	g.ScoreForA()        // 30/0
	g.ReverseScoreForA() // 15/0
	g.ScoreForB()        // 15/15
	g.ReverseScoreForB() // 15/0
	g.ScoreForA()        // 30/0

	a, b, _ := g.GetScore()
	if a != "30" || b != "0" {
		t.Errorf("interleaved result: got %s/%s, want 30/0", a, b)
	}
}

// ---------------------------------------------------------------------------
// Table-driven: all straight-win paths (no deuce)
// ---------------------------------------------------------------------------

func TestOneDeuceAllStraightWinPaths(t *testing.T) {
	type winPath struct {
		name   string
		points string // "A"/"B" sequence
		finalA string
		finalB string
	}
	tests := []winPath{
		{"A wins 4-0", "AAAA", "40*", "0"},
		{"A wins 4-1", "BAAAA", "40*", "15"},
		{"A wins 4-2 (path1)", "BBAAAA", "40*", "30"},
		{"A wins 4-2 (path2)", "BABAAA", "40*", "30"},
		{"B wins 0-4", "BBBB", "0", "40*"},
		{"B wins 1-4", "ABBBB", "15", "40*"},
		{"B wins 2-4 (path1)", "AABBBB", "30", "40*"},
		{"B wins 2-4 (path2)", "ABABBB", "30", "40*"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := newGame(OneDeuce)
			var a, b string
			var done bool
			for _, ch := range tc.points {
				if ch == 'A' {
					a, b, done = g.ScoreForA()
				} else {
					a, b, done = g.ScoreForB()
				}
			}
			if a != tc.finalA || b != tc.finalB || !done {
				t.Errorf("got (%s/%s done=%v), want (%s/%s done=true)",
					a, b, done, tc.finalA, tc.finalB)
			}
		})
	}
}
