package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BagSolver interface {
	CalculateSolution(problemInstance model.ProblemInstance) *model.FinalConfiguration
}
