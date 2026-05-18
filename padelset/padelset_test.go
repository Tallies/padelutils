package padelset

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type scoreStep struct {
	pointFor string // "A" or "B"
	wantA    string
	wantB    string
	wantDone bool
}

func playSequence(t *testing.T, set *PadelSet, steps []scoreStep) {
	t.Helper()

	for i, s := range steps {
		var gotA, gotB string
		var done bool

		switch s.pointFor {
		case "A":
			gotA, gotB, done = set.ScoreForA()
		case "B":
			gotA, gotB, done = set.ScoreForB()
		default:
			t.Fatalf("step %d: unknown pointFor %q", i, s.pointFor)
		}

		if gotA != s.wantA || gotB != s.wantB || done != s.wantDone {
			t.Fatalf(
				"step %d (point→%s): got (%s/%s done=%v), want (%s/%s done=%v)",
				i,
				s.pointFor,
				gotA,
				gotB,
				done,
				s.wantA,
				s.wantB,
				s.wantDone,
			)
		}
	}
}

func assertScore(
	t *testing.T,
	set *PadelSet,
	wantA string,
	wantB string,
	wantDone bool,
) {
	t.Helper()

	gotA, gotB, gotDone := set.GetScore()

	if gotA != wantA || gotB != wantB || gotDone != wantDone {
		t.Fatalf(
			"got (%s/%s done=%v), want (%s/%s done=%v)",
			gotA,
			gotB,
			gotDone,
			wantA,
			wantB,
			wantDone,
		)
	}
}

func newSet() PadelSet {
	return CreatePadelSetStandard()
}

// ---------------------------------------------------------------------------
// Initial State
// ---------------------------------------------------------------------------

func TestInitialScore(t *testing.T) {
	set := newSet()

	assertScore(t, &set, "0", "0", false)

	if set.IsComplete() {
		t.Fatal("new set should not be completed")
	}
}

// ---------------------------------------------------------------------------
// Team A wins
// ---------------------------------------------------------------------------

func TestTeamAWins_6_0(t *testing.T) {
	set := newSet()

	playSequence(t, &set, []scoreStep{
		{"A", "1", "0", false},
		{"A", "2", "0", false},
		{"A", "3", "0", false},
		{"A", "4", "0", false},
		{"A", "5", "0", false},
		{"A", "6", "0", true},
	})
}

func TestTeamAWins_6_1(t *testing.T) {
	set := newSet()

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

func TestTeamAWins_6_2(t *testing.T) {
	set := newSet()

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

func TestTeamAWins_6_3(t *testing.T) {
	set := newSet()

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

func TestTeamAWins_6_4(t *testing.T) {
	set := newSet()

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

func TestTeamAWins_7_5(t *testing.T) {
	set := newSet()

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

func TestTeamAWins_7_6(t *testing.T) {
	set := newSet()

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

func TestTeamBWins_6_0(t *testing.T) {
	set := newSet()

	playSequence(t, &set, []scoreStep{
		{"B", "0", "1", false},
		{"B", "0", "2", false},
		{"B", "0", "3", false},
		{"B", "0", "4", false},
		{"B", "0", "5", false},
		{"B", "0", "6", true},
	})
}

func TestTeamBWins_7_5(t *testing.T) {
	set := newSet()

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

func TestTeamBWins_7_6(t *testing.T) {
	set := newSet()

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

func TestRollbackSingleStep(t *testing.T) {
	set := newSet()

	set.ScoreForA()
	assertScore(t, &set, "1", "0", false)

	set.ReverseScoreForA()
	assertScore(t, &set, "0", "0", false)
}

func TestRollbackMultipleSteps(t *testing.T) {
	set := newSet()

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

func TestRollbackAtRootDoesNothing(t *testing.T) {
	set := newSet()

	set.ReverseScoreForA()

	assertScore(t, &set, "0", "0", false)
}

func TestRollbackAfterCompletedSet(t *testing.T) {
	set := newSet()

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

func TestCompletedSetCannotAdvance(t *testing.T) {
	set := newSet()

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
