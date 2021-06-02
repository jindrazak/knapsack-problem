package model

type Solution struct {
	Id            int
	Price         int
	Configuration FinalConfiguration
}

func (solution Solution) IsSolvable(problemInstance ProblemInstance) bool {
	solutionWeight := 0
	for i := 0; i < len(solution.Configuration); i++ {
		if solution.Configuration[i] {
			solutionWeight += problemInstance.Items[i].Weight
		}
	}
	return problemInstance.MinimumPrice <= solution.Price &&
		solutionWeight <= problemInstance.Bag.Capacity
}
