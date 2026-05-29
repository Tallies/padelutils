package node

import (
	"fmt"
	"padelutils/internal/stack"
	"strconv"
)

// Soring tree node
type Node struct {
	ScoreA string
	ScoreB string
	Left   *Node
	Right  *Node
}

func (node *Node) GetScores() (int, int) {
	scoreA, errA := strconv.Atoi(node.ScoreA)
	scoreB, errB := strconv.Atoi(node.ScoreB)

	if errA != nil || errB != nil {
		panic(fmt.Sprintf("Invalid score for %v_%v", node.ScoreA, node.ScoreB))
	}
	return scoreA, scoreB
}

func (root *Node) FindNode(scoreA string, scoreB string) *Node {
	//fmt.Printf("Start FindNode. scoreA=%v, scoreB=%v\n", scoreA, scoreB)
	node := root
	nodeStack := stack.Stack[*Node]{}
	checkedCache := make(map[string]bool)

	for {
		if node.ScoreA == scoreA && node.ScoreB == scoreB {
			return node
		}

		nodeName := node.ScoreA + "_" + node.ScoreB
		done, ok := checkedCache[nodeName]

		if !ok { // Haven't found this node at all.
			nodeStack.Push(node)
			checkedCache[nodeName] = false // if it exist and value is false, the only left so far has been checked
			node = node.Left
		} else if !done {
			checkedCache[nodeName] = true
			node = node.Right
		} else {
			if nodeStack.IsEmpty() {
				panic(fmt.Sprintf("Cannot find node for %v-%v", scoreA, scoreB))
			}
			node, _ = nodeStack.Pop()
		}
		if node == nil {
			if nodeStack.IsEmpty() {
				panic(fmt.Sprintf("Cannot find node for %v-%v", scoreA, scoreB))
			}
			node, _ = nodeStack.Pop()
		}
	}
}

func (node *Node) IsLeaf() bool {
	return node.Left == nil && node.Right == nil
}

// General node helper functions
func GetNodeName(teamA int, teamB int) string {
	return fmt.Sprintf("%d_%d", teamA, teamB)
}

func CreateNode(teamAScore string, teamBScore string) *Node {
	return &Node{
		ScoreA: teamAScore,
		ScoreB: teamBScore,
	}
}

func CreateNodes(scoreMax int, excludeFn func(int, int) bool) map[string]*Node {
	nodeCache := make(map[string]*Node)

	for scoreA := 0; scoreA <= scoreMax; scoreA++ {
		for scoreB := 0; scoreB <= scoreMax; scoreB++ {
			// fmt.Printf("CreateNode for %v_%v? %v\n", scoreA, scoreB, excludeFn(scoreA, scoreB))
			if excludeFn(scoreA, scoreB) {
				continue
			}

			nodeName := GetNodeName(scoreA, scoreB)
			node, ok := nodeCache[nodeName]
			if !ok {
				node = CreateNode(strconv.Itoa(scoreA), strconv.Itoa(scoreB))
				nodeCache[nodeName] = node
			}
		}
	}
	return nodeCache
}

func LinkNodes(nodeCache map[string]*Node,
	leftChildOfFn func(string, string) (string, string, bool),
	rightChildOfFn func(string, string) (string, string, bool),
) *Node {
	for key, node := range nodeCache {
		if key == "0_0" {
			continue
		}

		parentScoreA, parentScoreB, ok := leftChildOfFn(node.ScoreA, node.ScoreB)
		if ok {
			// fmt.Printf("%v_%v is left to parent %v_%v.\n",
			// 	node.ScoreA,
			// 	node.ScoreB,
			// 	parentScoreA,
			// 	parentScoreB)
			parentNode := nodeCache[parentScoreA+"_"+parentScoreB]
			parentNode.Left = node
		}

		parentScoreA, parentScoreB, ok = rightChildOfFn(node.ScoreA, node.ScoreB)
		if ok {
			// fmt.Printf("%v_%v is left to parent %v_%v.\n",
			// 	node.ScoreA,
			// 	node.ScoreB,
			// 	parentScoreA,
			// 	parentScoreB)
			parentNode := nodeCache[parentScoreA+"_"+parentScoreB]
			parentNode.Right = node
		}
	}
	return nodeCache["0_0"]
}
