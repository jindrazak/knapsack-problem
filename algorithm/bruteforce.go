package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BruteforceBagSolver struct {
}

func (bagSolver BruteforceBagSolver) CalculateSolution(problemInstance model.ProblemInstance) (*model.FinalConfiguration, int) {
	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
	resultChannel := make(chan *model.FinalConfiguration)
	visitedConfigurationsChannel := make(chan int)

	go bagSolver.goGetSolutionRec(problemInstance, configuration, resultChannel, visitedConfigurationsChannel)
	result := <-resultChannel
	visitedConfigurations := <-visitedConfigurationsChannel
	return result, visitedConfigurations
}

func (bagSolver BruteforceBagSolver) getSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) (*model.FinalConfiguration, int) {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		return bagSolver.copyConfigurationIfValid(problemInstance, configuration), 1
	}
	//continue with recursion, set currently changed boolean to false
	originalMask := configuration.MaskIndex
	configuration.SetNextFlag(false)
	leftSolution, leftVisitedConfigurations := bagSolver.getSolutionRec(problemInstance, configuration)

	//continue with recursion, set currently changed boolean to true
	configuration.MaskIndex = originalMask
	configuration.SetNextFlag(true)
	rightSolution, rightVisitedConfigurations := bagSolver.getSolutionRec(problemInstance, configuration)
	return bagSolver.pickBetterConfiguration(leftSolution, rightSolution, problemInstance), leftVisitedConfigurations + rightVisitedConfigurations
}

func (bagSolver BruteforceBagSolver) goGetSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration, resultChannel chan *model.FinalConfiguration, visitedConfigurationsChannel chan int) {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		resultChannel <- bagSolver.copyConfigurationIfValid(problemInstance, configuration)
		return
	}

	if configuration.MaskIndex > 6 {
		//There's no need to traverse the whole tree using goroutines. At some point, continue sequentially.
		//Empirically found that 6th level in the tree is approximately the right spot
		result, visitedConfigurations := bagSolver.getSolutionRec(problemInstance, configuration)
		resultChannel <- result
		visitedConfigurationsChannel <- visitedConfigurations
		return
	}
	leftChannel := make(chan *model.FinalConfiguration)
	rightChannel := make(chan *model.FinalConfiguration)
	leftVisitedConfigurationsChannel := make(chan int)
	rightVisitedConfigurationsChannel := make(chan int)
	//continue with recursion, set currently changed boolean to false
	leftConfiguration := configuration
	rightConfiguration := leftConfiguration.Clone()
	leftConfiguration.SetNextFlag(false)
	rightConfiguration.SetNextFlag(true)

	//let the slave gophers do their job
	go bagSolver.goGetSolutionRec(problemInstance, leftConfiguration, leftChannel, leftVisitedConfigurationsChannel)
	go bagSolver.goGetSolutionRec(problemInstance, rightConfiguration, rightChannel, rightVisitedConfigurationsChannel)
	var leftSolution, rightSolution *model.FinalConfiguration
	var leftVisitedConfigurations, rightVisitedConfigurations int
	for i := 0; i < 4; i++ {
		select {
		case leftSolution = <-leftChannel:
		case rightSolution = <-rightChannel:
		case leftVisitedConfigurations = <-leftVisitedConfigurationsChannel:
		case rightVisitedConfigurations = <-rightVisitedConfigurationsChannel:
		}
	}

	resultChannel <- bagSolver.pickBetterConfiguration(leftSolution, rightSolution, problemInstance)
	visitedConfigurationsChannel <- rightVisitedConfigurations + leftVisitedConfigurations
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
