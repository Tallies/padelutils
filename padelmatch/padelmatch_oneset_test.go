package padelmatch

import (
	"testing"
)

// ---------------------------------------------------------------------------------
// OneSet Padel Match Test
// ---------------------------------------------------------------------------------
func TestOneSet_CreatePadelMatchStandard(t *testing.T) {
	match := newMatch(OneSet)

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

func TestOneSet_ScoreForA_WinStraight(t *testing.T) {
	match := newMatch(OneSet)

	playSequence(t, &match, []scoreStep{
		{"A", "1", "0", true},
	})
}

func TestOneSet_ScoreForB_WinStraight(t *testing.T) {
	match := newMatch(OneSet)

	playSequence(t, &match, []scoreStep{
		{"B", "0", "1", true},
	})
}

func TestOneSet_ScoreForA_WinRemainUnchangedWithAdditionalAddScore(t *testing.T) {
	match := newMatch(OneSet)

	playSequence(t, &match, []scoreStep{
		{"A", "1", "0", true},
		{"A", "1", "0", true},
		{"B", "1", "0", true},
	})
}

func TestOneSet_ScoreForB_WinRemainUnchangedWithAdditionalAddScore(t *testing.T) {
	match := newMatch(OneSet)

	playSequence(t, &match, []scoreStep{
		{"B", "0", "1", true},
		{"A", "0", "1", true},
		{"B", "0", "1", true},
	})
}

// -----------------------------------------------------------------------------------------------------------
// ReverseScore Tests
// -----------------------------------------------------------------------------------------------------------
func TestOneSet_ReverseScoreForA_SingleStep(t *testing.T) {
	match := newMatch(OneSet)

	match.ScoreForA()
	a, b, done := match.ReverseScoreForA() // back to 0/0
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForA: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

func TestOneSet_ReverseScoreForB_SingleStep(t *testing.T) {
	match := newMatch(OneSet)
	match.ScoreForB()
	a, b, done := match.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("after ReverseScoreForB: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

// Rollback at root (no history) is a no-op
func TestOneSet_ReversePoint_AtRoot_IsNoOp(t *testing.T) {
	match := newMatch(OneSet)
	a, b, done := match.ReverseScoreForA()
	if a != "0" || b != "0" || done {
		t.Errorf("reverse at root: got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
	a, b, done = match.ReverseScoreForB()
	if a != "0" || b != "0" || done {
		t.Errorf("reverse at root (B): got (%s/%s done=%v), want (0/0 false)", a, b, done)
	}
}

func TestOneSet_IsCompleteInitiallyFalse(t *testing.T) {
	match := newMatch(OneSet)

	if match.IsComplete() {
		t.Errorf("expected match to not be complete")
	}
}
