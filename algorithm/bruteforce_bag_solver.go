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
	resultChannel := make(chan *model.FinalConfiguration)
	go solver.goGetSolutionRec(problemInstance, configuration, resultChannel)
	result := <-resultChannel
	return result
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

func (solver *BruteforceBagSolver) goGetSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration, resultChannel chan *model.FinalConfiguration) {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		solver.VisitedConfigurations++
		totalWeight := problemInstance.CalculateTotalWeight(configuration.Flags)
		totalPrice := problemInstance.CalculateTotalPrice(configuration.Flags)
		fitsInBag := totalWeight <= problemInstance.Bag.Capacity
		fulfilsMinimumPrice := totalPrice >= problemInstance.MinimumPrice
		if fitsInBag && fulfilsMinimumPrice {
			var configurationCopy = configuration.Flags.Clone()
			resultChannel <- &configurationCopy
		} else {
			resultChannel <- nil
		}
	} else {
		if configuration.MaskIndex > 6 {
			resultChannel <- solver.getSolutionRec(problemInstance, configuration)
			return
		}
		leftChannel := make(chan *model.FinalConfiguration)
		rightChannel := make(chan *model.FinalConfiguration)
		//continue with recursion, set currently changed boolean to false
		originalMask := configuration.MaskIndex
		configuration.SetNextFlag(false)

		go solver.goGetSolutionRec(problemInstance, configuration, leftChannel)
		rightConfiguration := configuration.Clone()
		rightConfiguration.MaskIndex = originalMask

		//continue with recursion, set currently changed boolean to true
		rightConfiguration.SetNextFlag(true)
		go solver.goGetSolutionRec(problemInstance, rightConfiguration, rightChannel)
		var leftSolution, rightSolution *model.FinalConfiguration
		//leftSolution = <-leftChannel
		//rightSolution = <-rightChannel
		for i := 0; i < 2; i++ {
			select {
			case leftSolution = <-leftChannel:
			case rightSolution = <-rightChannel:
			}
		}

		resultChannel <- pickBetterConfiguration(leftSolution, rightSolution, problemInstance)
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
