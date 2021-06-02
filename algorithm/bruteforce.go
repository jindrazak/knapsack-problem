package algorithm

import "github.com/jindrazak/knapsack-problem/model"

func CalculateBruteforceSolution(problemInstance model.ProblemInstance) *model.FinalConfiguration {
	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
	resultChannel := make(chan *model.FinalConfiguration)

	go goGetSolutionRec(problemInstance, configuration, resultChannel)
	result := <-resultChannel
	return result
}

func getSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) *model.FinalConfiguration {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		if isValidConfiguration(problemInstance, configuration.Flags) {
			var configurationCopy = configuration.Flags.Clone()
			return &configurationCopy
		} else {
			return nil
		}
	} else {
		//continue with recursion, set currently changed boolean to false
		originalMask := configuration.MaskIndex
		configuration.SetNextFlag(false)
		leftSolution := getSolutionRec(problemInstance, configuration)

		//continue with recursion, set currently changed boolean to true
		configuration.MaskIndex = originalMask
		configuration.SetNextFlag(true)
		rightSolution := getSolutionRec(problemInstance, configuration)
		return pickBetterConfiguration(leftSolution, rightSolution, problemInstance)
	}

}

func goGetSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration, resultChannel chan *model.FinalConfiguration) {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		if isValidConfiguration(problemInstance, configuration.Flags) {
			var configurationCopy = configuration.Flags.Clone()
			resultChannel <- &configurationCopy
		} else {
			resultChannel <- nil
		}
	} else {
		if configuration.MaskIndex > 6 {
			//There's no need to traverse the whole tree using goroutines. At some point, continue sequentially.
			//Empirically found that 6th level in the tree is approximately the right spot
			resultChannel <- getSolutionRec(problemInstance, configuration)
			return
		}
		leftChannel := make(chan *model.FinalConfiguration)
		rightChannel := make(chan *model.FinalConfiguration)
		//continue with recursion, set currently changed boolean to false
		leftConfiguration := configuration
		rightConfiguration := leftConfiguration.Clone()
		leftConfiguration.SetNextFlag(false)
		rightConfiguration.SetNextFlag(true)

		//let the slave gophers do their job
		go goGetSolutionRec(problemInstance, leftConfiguration, leftChannel)
		go goGetSolutionRec(problemInstance, rightConfiguration, rightChannel)
		var leftSolution, rightSolution *model.FinalConfiguration
		for i := 0; i < 2; i++ {
			select {
			case leftSolution = <-leftChannel:
			case rightSolution = <-rightChannel:
			}
		}

		resultChannel <- pickBetterConfiguration(leftSolution, rightSolution, problemInstance)
	}

}

func isValidConfiguration(instance model.ProblemInstance, configuration model.FinalConfiguration) bool {
	totalWeight := instance.CalculateTotalWeight(configuration)
	totalPrice := instance.CalculateTotalPrice(configuration)
	fitsInBag := totalWeight <= instance.Bag.Capacity
	fulfilsMinimumPrice := totalPrice >= instance.MinimumPrice
	return fitsInBag && fulfilsMinimumPrice
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
