package padelgame

import (
	"fmt"
	"testing"
)

// ---------------------------------------------------------------------------
// Initial state
// ---------------------------------------------------------------------------

func TestGoldenPointInitialScore(t *testing.T) {
	g := newGame(GoldenPoint)
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

func TestGoldenPointTeamAWins_40_0(t *testing.T) {
	g := newGame(GoldenPoint)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"A", "40*", "0", true}, // terminal win node
	})
}

func TestGoldenPointTeamAWins_40_15(t *testing.T) {
	g := newGame(GoldenPoint)
	playSequence(t, &g, []scoreStep{
		{"B", "0", "15", false},
		{"A", "15", "15", false},
		{"A", "30", "15", false},
		{"A", "40", "15", false},
		{"A", "40*", "15", true},
	})
}

func TestGoldenPointTeamAWins_40_30(t *testing.T) {
	g := newGame(GoldenPoint)
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

func TestGoldenPointTeamBWins_0_40(t *testing.T) {
	g := newGame(GoldenPoint)
	playSequence(t, &g, []scoreStep{
		{"B", "0", "15", false},
		{"B", "0", "30", false},
		{"B", "0", "40", false},
		{"B", "0", "40*", true},
	})
}

func TestGoldenPointTeamBWins_15_40(t *testing.T) {
	g := newGame(GoldenPoint)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"B", "15", "15", false},
		{"B", "15", "30", false},
		{"B", "15", "40", false},
		{"B", "15", "40*", true},
	})
}

func TestGoldenPointTeamBWins_30_40(t *testing.T) {
	g := newGame(GoldenPoint)
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
func TestGoldenPointGoldenPoint_AWins(t *testing.T) {
	g := newGame(GoldenPoint)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"B", "40", "15", false},
		{"B", "40", "30", false},
		{"B", "GP", "GP", false},
		{"A", "GP*", "40", true},
	})
}

// 40-30 → B scores → D1/D1 → B scores → Advantage B (40/A1) → B wins
func TestGoldenPointGoldenPoint_BWins(t *testing.T) {
	g := newGame(GoldenPoint)
	playSequence(t, &g, []scoreStep{
		{"A", "15", "0", false},
		{"A", "30", "0", false},
		{"A", "40", "0", false},
		{"B", "40", "15", false},
		{"B", "40", "30", false},
		{"B", "GP", "GP", false},
		{"B", "40", "GP*", true},
	})
}

// ---------------------------------------------------------------------------
// Score stays frozen after game is complete
// ---------------------------------------------------------------------------

func TestGoldenPointScoreFrozenAfterWin(t *testing.T) {
	g := newGame(GoldenPoint)
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

func TestGoldenPointGetScore_MatchesIsComplete(t *testing.T) {
	g := newGame(GoldenPoint)
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
func TestGoldenPointReverseScoreForA_SingleStep(t *testing.T) {
	g := newGame(GoldenPoint)
	g.ScoreForA()                      // 0/0 → 15/0
	a, b, done := g.ReverseScoreForA() // back to 0/0
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForA: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

func TestGoldenPointReverseScoreForB_SingleStep(t *testing.T) {
	g := newGame(GoldenPoint)
	g.ScoreForB()
	a, b, done := g.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForB: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback restores the previous node, not always 0/0
func TestGoldenPointReversePoint_MidGame(t *testing.T) {
	g := newGame(GoldenPoint)
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
func TestGoldenPointReversePoint_FullUnwind(t *testing.T) {
	g := newGame(GoldenPoint)
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
func TestGoldenPointReversePoint_AtRoot_IsNoOp(t *testing.T) {
	g := newGame(GoldenPoint)
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
func TestGoldenPointReversePoint_FromWinNode(t *testing.T) {
	g := newGame(GoldenPoint)
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
func TestGoldenPointReversePoint_ThroughGoldenPoint(t *testing.T) {
	g := newGame(GoldenPoint)
	// Advance to GP/GP
	g.ScoreForA()
	g.ScoreForA()
	g.ScoreForA() // 40/0
	g.ScoreForB()
	g.ScoreForB() // 40/30
	g.ScoreForB() // GP/GP
	g.ScoreForB() // 40/GP*

	a, b, _ := g.GetScore()
	if a != "40" || b != "GP*" {
		t.Fatalf("expected 40/GP*, got %s/%s", a, b)
	}

	g.ReverseScoreForB() // back to GP/GP
	a, b, _ = g.GetScore()
	if a != "GP" || b != "GP" {
		t.Errorf("rollback from 40/GP*: got %s/%s, want GP/GP", a, b)
	}
}

// Rollback through Deuce 1
func TestGoldenPointReversePoint_ThroughGoldenPointFromLowerSide(t *testing.T) {
	g := newGame(GoldenPoint)
	// Advance to GP/GP
	g.ScoreForA()
	g.ScoreForA()
	g.ScoreForA() // 40/0
	g.ScoreForB()
	g.ScoreForB() // 40/30
	g.ScoreForB() // GP/GP
	g.ScoreForB() // 40/GP*

	a, b, _ := g.GetScore()
	if a != "40" || b != "GP*" {
		t.Fatalf("expected 40/GP*, got %s/%s", a, b)
	}

	g.ReverseScoreForA() // back to 30/40
	a, b, _ = g.GetScore()
	if a != "GP" || b != "GP" {
		t.Errorf("rollback A from 40/GP*: got %s/%s, want GP/GP", a, b)
	}
}

// Interleaving forward and reverse moves
func TestGoldenPointReversePoint_Interleaved(t *testing.T) {
	g := newGame(GoldenPoint)
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

func TestGoldenPointAllStraightWinPaths(t *testing.T) {
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
			g := newGame(GoldenPoint)
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
