package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BruteforceBagSolver struct {
	visitedConfigurations int
}

func (bagSolver BruteforceBagSolver) CalculateSolution(problemInstance model.ProblemInstance) *model.FinalConfiguration {
	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
	resultChannel := make(chan *model.FinalConfiguration)

	go bagSolver.goGetSolutionRec(problemInstance, configuration, resultChannel)
	result := <-resultChannel
	return result
}

func (bagSolver BruteforceBagSolver) getSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) *model.FinalConfiguration {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		return bagSolver.copyConfigurationIfValid(problemInstance, configuration)
	}
	//continue with recursion, set currently changed boolean to false
	originalMask := configuration.MaskIndex
	configuration.SetNextFlag(false)
	leftSolution := bagSolver.getSolutionRec(problemInstance, configuration)

	//continue with recursion, set currently changed boolean to true
	configuration.MaskIndex = originalMask
	configuration.SetNextFlag(true)
	rightSolution := bagSolver.getSolutionRec(problemInstance, configuration)
	return bagSolver.pickBetterConfiguration(leftSolution, rightSolution, problemInstance)
}

func (bagSolver BruteforceBagSolver) goGetSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration, resultChannel chan *model.FinalConfiguration) {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		resultChannel <- bagSolver.copyConfigurationIfValid(problemInstance, configuration)
		return
	}

	if configuration.MaskIndex > 6 {
		//There's no need to traverse the whole tree using goroutines. At some point, continue sequentially.
		//Empirically found that 6th level in the tree is approximately the right spot
		resultChannel <- bagSolver.getSolutionRec(problemInstance, configuration)
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
	go bagSolver.goGetSolutionRec(problemInstance, leftConfiguration, leftChannel)
	go bagSolver.goGetSolutionRec(problemInstance, rightConfiguration, rightChannel)
	var leftSolution, rightSolution *model.FinalConfiguration
	for i := 0; i < 2; i++ {
		select {
		case leftSolution = <-leftChannel:
		case rightSolution = <-rightChannel:
		}
	}

	resultChannel <- bagSolver.pickBetterConfiguration(leftSolution, rightSolution, problemInstance)

}

func (_ BruteforceBagSolver) pickBetterConfiguration(a, b *model.FinalConfiguration, instance model.ProblemInstance) *model.FinalConfiguration {
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

func (bagSolver BruteforceBagSolver) copyConfigurationIfValid(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) *model.FinalConfiguration {
	if !problemInstance.IsValidConfiguration(configuration.Flags) {
		return nil
	}
	configurationClone := configuration.Flags.Clone()
	return &configurationClone
}
