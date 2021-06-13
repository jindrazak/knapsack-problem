package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BruteforceBagSolver struct {
}

func (bagSolver BruteforceBagSolver) CalculateSolution(problemInstance model.ProblemInstance) model.CalculatedSolution {
	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
	resultChannel := make(chan model.CalculatedSolution)

	go bagSolver.goGetSolutionRec(problemInstance, configuration, resultChannel)
	return <-resultChannel
}

func (bagSolver BruteforceBagSolver) getSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) model.CalculatedSolution {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		return model.CalculatedSolution{
			Configuration:         bagSolver.copyConfigurationIfValid(problemInstance, configuration),
			VisitedConfigurations: 1,
		}
	}
	//continue with recursion, set currently changed boolean to false
	originalMask := configuration.MaskIndex
	configuration.SetNextFlag(false)
	leftSolution := bagSolver.getSolutionRec(problemInstance, configuration)

	//continue with recursion, set currently changed boolean to true
	configuration.MaskIndex = originalMask
	configuration.SetNextFlag(true)
	rightSolution := bagSolver.getSolutionRec(problemInstance, configuration)
	return bagSolver.mergeCalculatedSolutions(leftSolution, rightSolution, problemInstance)
}

func (bagSolver BruteforceBagSolver) goGetSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration, resultChannel chan model.CalculatedSolution) {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		resultChannel <- model.CalculatedSolution{
			Configuration:         bagSolver.copyConfigurationIfValid(problemInstance, configuration),
			VisitedConfigurations: 1,
		}
		return
	}

	if configuration.MaskIndex > 6 {
		//There's no need to traverse the whole tree using goroutines. At some point, continue sequentially.
		//Empirically found that 6th level in the tree is approximately the right spot
		resultChannel <- bagSolver.getSolutionRec(problemInstance, configuration)
		return
	}
	leftChannel := make(chan model.CalculatedSolution)
	rightChannel := make(chan model.CalculatedSolution)
	//continue with recursion, set currently changed boolean to false
	leftConfiguration := configuration
	rightConfiguration := leftConfiguration.Clone()
	leftConfiguration.SetNextFlag(false)
	rightConfiguration.SetNextFlag(true)

	//let the slave gophers do their job
	go bagSolver.goGetSolutionRec(problemInstance, leftConfiguration, leftChannel)
	go bagSolver.goGetSolutionRec(problemInstance, rightConfiguration, rightChannel)
	var leftSolution, rightSolution model.CalculatedSolution
	for i := 0; i < 2; i++ {
		select {
		case leftSolution = <-leftChannel:
		case rightSolution = <-rightChannel:
		}
	}

	resultChannel <- bagSolver.mergeCalculatedSolutions(leftSolution, rightSolution, problemInstance)
}

func (bagSolver BruteforceBagSolver) mergeCalculatedSolutions(a, b model.CalculatedSolution, instance model.ProblemInstance) model.CalculatedSolution {
	return model.CalculatedSolution{
		Configuration:         instance.PickBetterConfiguration(a.Configuration, b.Configuration),
		VisitedConfigurations: a.VisitedConfigurations + b.VisitedConfigurations,
	}
}

func (bagSolver BruteforceBagSolver) copyConfigurationIfValid(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) *model.FinalConfiguration {
	if !problemInstance.IsValidConfiguration(configuration.Flags) {
		return nil
	}
	configurationClone := configuration.Flags.Clone()
	return &configurationClone
}
