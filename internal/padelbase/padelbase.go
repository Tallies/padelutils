package padelbase

type PadelBase interface {
	// ScoreForA advance the score for the game/set/match with a point for A.
	ScoreForA() (string, string, bool)
	// ScoreForA advance the score for the game/set/matchi with a point for B.
	ScoreForB() (string, string, bool)
	// ReverseScoreForA reduces the score for the game/set/match with a point for A.
	ReverseScoreForA() (string, string, bool)
	// ReverseScoreForA reduces the score for the game/set/match with a point for B.
	ReverseScoreForB() (string, string, bool)
	// GetScore returns the current game/set/match state as the current points 
	// for A, the current points for B, and a boolean indicating whether
	// the game/set/match is complete or not.
	GetScore() (string, string, bool)
	// Returns a boolean indicating whether the game/set/match is complete.
	IsMatchComplete() bool
}
