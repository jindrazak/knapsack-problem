package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BruteforceBagSolver struct {
	VisitedConfigurations int
}

func (solver *BruteforceBagSolver) Reset() {
	solver.VisitedConfigurations = 0
}

func (solver *BruteforceBagSolver) GetSolution(problemInstance model.ProblemInstance) *model.FinalConfiguration {
	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
	return solver.getSolutionRec(problemInstance, configuration)
}

func (solver *BruteforceBagSolver) getSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) *model.FinalConfiguration {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		solver.VisitedConfigurations++
		totalWeight := problemInstance.CalculateTotalWeight(configuration.Flags)
		totalPrice := problemInstance.CalculateTotalPrice(configuration.Flags)
		fitsInBag := totalWeight <= problemInstance.Bag.Capacity
		fulfilsMinimumPrice := totalPrice >= problemInstance.MinimumPrice
		if fitsInBag && fulfilsMinimumPrice {
			var configurationCopy = configuration.Flags.Clone()
			return &configurationCopy
		} else {
			return nil
		}
	} else {
		//continue with recursion, set currently changed boolean to false
		originalMask := configuration.MaskIndex
		configuration.SetNextFlag(false)
		leftSolution := solver.getSolutionRec(problemInstance, configuration)
		configuration.MaskIndex = originalMask

		//continue with recursion, set currently changed boolean to true
		configuration.SetNextFlag(true)
		rightSolution := solver.getSolutionRec(problemInstance, configuration)
		return pickBetterConfiguration(leftSolution, rightSolution, problemInstance)
	}

}

func pickBetterConfiguration(a, b *model.FinalConfiguration, instance model.ProblemInstance) *model.FinalConfiguration {
	if a != nil && b != nil {
		aPrice := instance.CalculateTotalPrice(*a)
		bPrice := instance.CalculateTotalPrice(*b)
		if aPrice < bPrice {
			return b
		} else {
			return a
		}

	} else if a != nil && b == nil {
		return a
	} else if a == nil && b != nil {
		return b
	} else {
		return nil
	}
}
