package padelbase

import (
	"fmt"
	"strconv"
)

func reduceScore(score string) string {
	scoreInt, err := strconv.Atoi(score)
	if err != nil {
		panic(fmt.Sprintf("'%v' is an invalid score for a set.", score))
	}
	return strconv.Itoa(scoreInt-1)
}

