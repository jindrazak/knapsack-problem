package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BagSolver interface {
	Reset()
	GetSolution(problemInstance model.ProblemInstance) bool
}
