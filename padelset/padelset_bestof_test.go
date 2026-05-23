package padelset

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Initial State
// ---------------------------------------------------------------------------

func TestBestOfInitialScore(t *testing.T) {
	set := newSet(BestOf)

	assertScore(t, &set, "0", "0", false)

	if set.IsComplete() {
		t.Fatal("new set should not be completed")
	}
}

// ---------------------------------------------------------------------------
// Team A wins
// ---------------------------------------------------------------------------

func TestBestOfTeamAWins_5_0(t *testing.T) {
	set := newSet(BestOf)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", true},
	})
}

func TestBestOfTeamAWins_4_1(t *testing.T) {
	set := newSet(BestOf)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"A", "1", "1", false},
		{"A", "2", "1", false},
		{"A", "3", "1", false},
		{"A", "4", "1", true},
	})
}

func TestBestOfTeamAWins_3_2(t *testing.T) {
	set := newSet(BestOf)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"A", "1", "2", false},
		{"A", "2", "2", false},
		{"A", "3", "2", true},
	})
}

// ---------------------------------------------------------------------------
// Team B wins
// ---------------------------------------------------------------------------

func TestBestOfTeamBWins_5_0(t *testing.T) {
	set := newSet(BestOf)

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"B", "0", "5", true},
	})
}

func TestBestOfTeamBWins_4_1(t *testing.T) {
	set := newSet(BestOf)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"B", "1", "1", false},
		{"B", "1", "2", false},
		{"B", "1", "3", false},
		{"B", "1", "4", true},
	})
}

func TestBestOfTeamBWins_3_2(t *testing.T) {
	set := newSet(BestOf)

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"B", "2", "1", false},
		{"B", "2", "2", false},
		{"B", "2", "3", true},
	})
}

// ---------------------------------------------------------------------------
// Rollback Tests
// ---------------------------------------------------------------------------

func TestBestOfRollbackSingleStep(t *testing.T) {
	set := newSet(BestOf)

	set.ScoreForA()
	assertScore(t, &set, "1", "0", false)

	set.ReverseScoreForA()
	assertScore(t, &set, "0", "0", false)
}

func TestBestOfRollbackMultipleSteps(t *testing.T) {
	set := newSet(BestOf)

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

func TestBestOfRollbackAtRootDoesNothing(t *testing.T) {
	set := newSet(BestOf)

	set.ReverseScoreForA()

	assertScore(t, &set, "0", "0", false)
}

func TestBestOfRollbackAfterCompletedSet(t *testing.T) {
	set := newSet(BestOf)

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

func TestBestOfCompletedSetCannotAdvance(t *testing.T) {
	set := newSet(BestOf)

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
