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

func newSet(padelSetType PadelSetType) PadelSet {
	switch padelSetType {
	case Standard:
		// print("Creating Standard PadelSet\n")
		return CreatePadelSetStandard()
	case FirstTo:
		// print("Creating FirstTo PadelSet[n")
		return CreatePadelSetFirstTo(5)
	case BestOf:
		// print("Creating BestOf PadelSet\n")
		return CreatePadelSetBestOf(5)
	}
	panic("Provided PadelSetType is invalid")
}
