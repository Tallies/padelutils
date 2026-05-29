package padelmatch

import (
	"fmt"
	"testing"
)

type scoreStep struct {
	scoreFor string // "A" or "B"
	wantA    string
	wantB    string
	wantDone bool
}

func newMatch(matchType PadelMatchType) PadelMatch {
	switch matchType {
	case Standard:
		return CreateStandardPadelMatch()
	case OneSet:
		return CreateOneSetPadelMatch()
	}
	panic(fmt.Sprintf("invalid type '%v'", matchType))
}

// playSequence drives a game through a sequence of points and asserts the
// score after every single point.
func playSequence(t *testing.T, match *PadelMatch, steps []scoreStep) {
	t.Helper()
	for i, s := range steps {
		var gotA, gotB string
		var done bool
		switch s.scoreFor {
		case "A":
			gotA, gotB, done = match.ScoreForA()
		case "B":
			gotA, gotB, done = match.ScoreForB()
		default:
			t.Fatalf("step %d: unknown scoreFor %q", i, s.scoreFor)
		}
		if gotA != s.wantA || gotB != s.wantB || done != s.wantDone {
			t.Errorf(
				"step %d (point→%s): got (%s/%s done=%v), want (%s/%s done=%v)",
				i, s.scoreFor,
				gotA, gotB, done,
				s.wantA, s.wantB, s.wantDone,
			)
		}
	}
}
