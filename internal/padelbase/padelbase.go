package padelbase

type PadelBase interface {
	ScoreForA() (string, string, bool)
	ScoreForB() (string, string, bool)
	ReverseScoreForA() (string, string, bool)
	ReverseScoreForB() (string, string, bool)
	GetScore() (string, string, bool)
	IsMatchComplete() bool
}
