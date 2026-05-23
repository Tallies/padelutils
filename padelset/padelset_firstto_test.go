package padelset

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Initial State
// ---------------------------------------------------------------------------

func TestFirstToInitialScore(t *testing.T) {
	set := newSet(FirstTo)

	assertScore(t, &set, "0", "0", false)

	if set.IsComplete() {
		t.Fatal("new set should not be completed")
	}
}

// ---------------------------------------------------------------------------
// Team A wins
// ---------------------------------------------------------------------------

func TestFirstToTeamAWins_5_0(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", true},
	})
}

func TestFirstToTeamAWins_5_1(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"A", "1", "1", false},
		{"A", "2", "1", false},
		{"A", "3", "1", false},
		{"A", "4", "1", false},
		{"A", "5", "1", true},
	})
}

func TestFirstToTeamAWins_5_2(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"A", "1", "2", false},
		{"A", "2", "2", false},
		{"A", "3", "2", false},
		{"A", "4", "2", false},
		{"A", "5", "2", true},
	})
}

func TestFirstToTeamAWins_5_3(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"A", "1", "3", false},
		{"A", "2", "3", false},
		{"A", "3", "3", false},
		{"A", "4", "3", false},
		{"A", "5", "3", true},
	})
}

func TestFirstToTeamAWins_5_4(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"A", "1", "4", false},
		{"A", "2", "4", false},
		{"A", "3", "4", false},
		{"A", "4", "4", false},
		{"A", "5", "4", true},
	})
}

// ---------------------------------------------------------------------------
// Team B wins
// ---------------------------------------------------------------------------

func TestFirstToTeamBWins_5_0(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"B", "0", "5", true},
	})
}

func TestFirstToTeamBWins_5_1(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"B", "1", "2", false},
		{"B", "1", "3", false},
		{"B", "1", "4", false},
		{"B", "1", "5", true},
	})
}

func TestFirstToTeamBWins_5_2(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"B", "2", "1", false},
		{"B", "2", "2", false},
		{"B", "2", "3", false},
		{"B", "2", "4", false},
		{"B", "2", "5", true},
	})
}

func TestFirstToTeamBWins_5_3(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"B", "3", "1", false},
		{"B", "3", "2", false},
		{"B", "3", "3", false},
		{"B", "3", "4", false},
		{"B", "3", "5", true},
	})
}

func TestFirstToTeamBWins_5_4(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"B", "4", "1", false},
		{"B", "4", "2", false},
		{"B", "4", "3", false},
		{"B", "4", "4", false},
		{"B", "4", "5", true},
	})
}

// ---------------------------------------------------------------------------
// Rollback Tests
// ---------------------------------------------------------------------------

func TestFirstToRollbackSingleStep(t *testing.T) {
	set := newSet(FirstTo)

	set.ScoreForA()
	assertScore(t, &set, "1", "0", false)

	set.ReverseScoreForA()
	assertScore(t, &set, "0", "0", false)
}

func TestFirstToRollbackMultipleSteps(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"A", "2", "1", false},
	})

	set.ReverseScoreForA()
	assertScore(t, &set, "1", "1", false)

	set.ReverseScoreForB()
	assertScore(t, &set, "1", "0", false)

	set.ReverseScoreForA()
	assertScore(t, &set, "0", "0", false)
}

func TestFirstToRollbackAtRootDoesNothing(t *testing.T) {
	set := newSet(FirstTo)

	set.ReverseScoreForA()

	assertScore(t, &set, "0", "0", false)
}

func TestFirstToRollbackAfterCompletedSet(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", true},
	})

	set.ReverseScoreForA()

	assertScore(t, &set, "4", "0", false)
}

// ---------------------------------------------------------------------------
// Completion State
// ---------------------------------------------------------------------------

func TestFirstToCompletedSetCannotAdvance(t *testing.T) {
	set := newSet(FirstTo)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", true},
	})

	set.ScoreForA()

	assertScore(t, &set, "5", "0", true)

	set.ScoreForB()

	assertScore(t, &set, "5", "0", true)
}
