package padelgame

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const testRallyScore = 15 

type scoreStep struct {
	pointFor string // "A" or "B"
	wantA    string
	wantB    string
	wantDone bool
}

// playSequence drives a game through a sequence of points and asserts the
// score after every single point.
func playSequence(t *testing.T, game *PadelGame, steps []scoreStep) {
	t.Helper()
	for i, s := range steps {
		var gotA, gotB string
		var done bool
		switch s.pointFor {
		case "A":
			gotA, gotB, done = game.ScoreForA()
		case "B":
			gotA, gotB, done = game.ScoreForB()
		default:
			t.Fatalf("step %d: unknown pointFor %q", i, s.pointFor)
		}
		if gotA != s.wantA || gotB != s.wantB || done != s.wantDone {
			t.Errorf(
				"step %d (point→%s): got (%s/%s done=%v), want (%s/%s done=%v)",
				i, s.pointFor,
				gotA, gotB, done,
				s.wantA, s.wantB, s.wantDone,
			)
		}
	}
}

func newGame(gameType PadelGameType) PadelGame {
	switch gameType {
	case Advantage:
		return CreatePadelGameAdvantage()
	case StarPoint:
		return CreatePadelGameStarPoint()
	case OneDeuce:
		return CreatePadelGameOneDeuce()
	case GoldenPoint:
		return CreatePadelGameGoldenPoint()
	case RallyScoring:
		return CreatePadelGameRallyScoring(testRallyScore)
	}
	panic("Invalid PadelGameType specified")
}
