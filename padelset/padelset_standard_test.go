package padelset

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Initial State
// ---------------------------------------------------------------------------

func TestStandardInitialScore(t *testing.T) {
	set := newSet(Standard)

	assertScore(t, &set, "0", "0", false)

	if set.IsComplete() {
		t.Fatal("new set should not be completed")
	}
}

// ---------------------------------------------------------------------------
// Team A wins
// ---------------------------------------------------------------------------

func TestStandardTeamAWins_6_0(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", false},
		{"A", "6", "0", true},
	})
}

func TestStandardTeamAWins_6_1(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"A", "1", "1", false},
		{"A", "2", "1", false},
		{"A", "3", "1", false},
		{"A", "4", "1", false},
		{"A", "5", "1", false},
		{"A", "6", "1", true},
	})
}

func TestStandardTeamAWins_6_2(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"A", "1", "2", false},
		{"A", "2", "2", false},
		{"A", "3", "2", false},
		{"A", "4", "2", false},
		{"A", "5", "2", false},
		{"A", "6", "2", true},
	})
}

func TestStandardTeamAWins_6_3(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"A", "1", "3", false},
		{"A", "2", "3", false},
		{"A", "3", "3", false},
		{"A", "4", "3", false},
		{"A", "5", "3", false},
		{"A", "6", "3", true},
	})
}

func TestStandardTeamAWins_6_4(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"A", "1", "4", false},
		{"A", "2", "4", false},
		{"A", "3", "4", false},
		{"A", "4", "4", false},
		{"A", "5", "4", false},
		{"A", "6", "4", true},
	})
}

func TestStandardTeamAWins_7_5(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"B", "0", "5", false},
		{"A", "1", "5", false},
		{"A", "2", "5", false},
		{"A", "3", "5", false},
		{"A", "4", "5", false},
		{"A", "5", "5", false},
		{"A", "6", "5", false},
		{"A", "7", "5", true},
	})
}

func TestStandardTeamAWins_7_6(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"A", "2", "1", false},
		{"B", "2", "2", false},
		{"A", "3", "2", false},
		{"B", "3", "3", false},
		{"A", "4", "3", false},
		{"B", "4", "4", false},
		{"A", "5", "4", false},
		{"B", "5", "5", false},
		{"A", "6", "5", false},
		{"B", "6", "6", false},
		{"A", "7", "6", true},
	})
}

// ---------------------------------------------------------------------------
// Team B wins
// ---------------------------------------------------------------------------

func TestStandardTeamBWins_6_0(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"B", "0", "5", false},
		{"B", "0", "6", true},
	})
}

func TestStandardTeamBWins_7_5(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", false},
		{"B", "5", "1", false},
		{"B", "5", "2", false},
		{"B", "5", "3", false},
		{"B", "5", "4", false},
		{"B", "5", "5", false},
		{"B", "5", "6", false},
		{"B", "5", "7", true},
	})
}

func TestStandardTeamBWins_7_6(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"A", "2", "1", false},
		{"B", "2", "2", false},
		{"A", "3", "2", false},
		{"B", "3", "3", false},
		{"A", "4", "3", false},
		{"B", "4", "4", false},
		{"A", "5", "4", false},
		{"B", "5", "5", false},
		{"A", "6", "5", false},
		{"B", "6", "6", false},
		{"B", "6", "7", true},
	})
}

// ---------------------------------------------------------------------------
// Rollback Tests
// ---------------------------------------------------------------------------

func TestStandardRollbackSingleStep(t *testing.T) {
	set := newSet(Standard)

	set.ScoreForA()
	assertScore(t, &set, "1", "0", false)

	set.ReverseScoreForA()
	assertScore(t, &set, "0", "0", false)
}

func TestStandardRollbackMultipleSteps(t *testing.T) {
	set := newSet(Standard)

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

func TestStandardRollbackAtRootDoesNothing(t *testing.T) {
	set := newSet(Standard)

	set.ReverseScoreForA()

	assertScore(t, &set, "0", "0", false)
}

func TestStandardRollbackAfterCompletedSet(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", false},
		{"A", "6", "0", true},
	})

	set.ReverseScoreForA()

	assertScore(t, &set, "5", "0", false)
}

// ---------------------------------------------------------------------------
// Completion State
// ---------------------------------------------------------------------------

func TestStandardCompletedSetCannotAdvance(t *testing.T) {
	set := newSet(Standard)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", false},
		{"A", "6", "0", true},
	})

	set.ScoreForA()

	assertScore(t, &set, "6", "0", true)

	set.ScoreForB()

	assertScore(t, &set, "6", "0", true)
}
