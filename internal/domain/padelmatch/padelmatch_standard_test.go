package padelmatch

import (
	"testing"
)

// ---------------------------------------------------------------------------------
// Standard Padel Match Test
// ---------------------------------------------------------------------------------
func TestStandard_CreatePadelMatchStandard(t *testing.T) {
	match := newMatch(Standard)

	scoreA, scoreB, complete := match.GetScore()

	if scoreA != "0" {
		t.Errorf("expected scoreA to be '0', got '%s'", scoreA)
	}

	if scoreB != "0" {
		t.Errorf("expected scoreB to be '0', got '%s'", scoreB)
	}

	if complete {
		t.Errorf("expected new match to not be complete")
	}
}

func TestStandard_ScoreForA_WinStraight(t *testing.T) {
	match := newMatch(Standard)

	playSequence(t, &match, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", true},
	})
}

func TestStandard_ScoreForA_WinNonStraight(t *testing.T) {
	match := newMatch(Standard)

	playSequence(t, &match, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"A", "2", "1", true},
	})
}

func TestStandard_ScoreForB_WinStraight(t *testing.T) {
	match := newMatch(Standard)

	playSequence(t, &match, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", true},
	})
}

func TestStandard_ScoreForB_WinNonStraight(t *testing.T) {
	match := newMatch(Standard)

	playSequence(t, &match, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"B", "1", "2", true},
	})
}

func TestStandard_ScoreForA_WinRemainUnchangedWithAdditionalAddScore(t *testing.T) {
	match := newMatch(Standard)

	playSequence(t, &match, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", true},
		{"A", "2", "0", true},
		{"B", "2", "0", true},
	})
}

func TestStandard_ScoreForB_WinRemainUnchangedWithAdditionalAddScore(t *testing.T) {
	match := newMatch(Standard)

	playSequence(t, &match, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", true},
		{"A", "0", "2", true},
		{"B", "0", "2", true},
	})
}

// -----------------------------------------------------------------------------------------------------------
// ReverseScore Test
// -----------------------------------------------------------------------------------------------------------
func TestStandard_ReverseScoreForA_SingleStep(t *testing.T) {
	match := newMatch(Standard)

	match.ScoreForA()
	a, b, done := match.ReverseScoreForA() // back to 0/0
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForA: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

func TestStandard_ReverseScoreForB_SingleStep(t *testing.T) {
	match := newMatch(Standard)
	match.ScoreForB()
	a, b, done := match.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForB: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback restores the previous node, not always 0/0
func TestStandard_ReversePoint_MidGame(t *testing.T) {
	match := newMatch(Standard)
	match.ScoreForA() // 1/0
	match.ScoreForB() // 1/1

	a, b, _ := match.ReverseScoreForA() // back to 1/0
	if a != "0" || b != "1" {
		t.Errorf("rollback from 1/1: got %s/%s, want 0/1", a, b)
	}

	a, b, _ = match.ReverseScoreForA() // back to 15/0
	if a != "0" || b != "1" {
		t.Errorf("rollback from 0/1: got %s/%s, want 0/1", a, b)
	}
}

// Rollback at root (no history) is a no-op
func TestStandard_ReversePoint_AtRoot_IsNoOp(t *testing.T) {
	match := newMatch(Standard)
	a, b, done := match.ReverseScoreForA()
	if a != "0" || b != "0" || done {
		t.Errorf("reverse at root: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
	a, b, done = match.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("reverse at root (B): got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback from a terminal win node restores the pre-win state
func TestStandard_ReversePoint_FromWinNode(t *testing.T) {
	match := newMatch(Standard)
	match.ScoreForA() // 1/0
	match.ScoreForA() // 2/0

	a, b, done := match.ReverseScoreForA()
	if a != "1" || b != "0" || done {
		t.Errorf("rollback from win: got (%s/%s done=%v), want (1/0 false)", a, b, done)
	}
	// Game is live again — should be able to continue
	match.ScoreForB() // 1/1 (not terminal)
	a, b, _ = match.GetScore()
	if a != "1" || b != "1" {
		t.Errorf("after resuming: got %s/%s, want 1/1", a, b)
	}
}

func TestStandard_IsCompleteInitiallyFalse(t *testing.T) {
	match := newMatch(Standard)

	if match.IsComplete() {
		t.Errorf("expected match to not be complete")
	}
}
