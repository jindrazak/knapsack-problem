package algorithm

import "github.com/jindrazak/knapsack-problem/model"

type BruteforceBagSolver struct {
	VisitedConfigurations int
}

func (solver *BruteforceBagSolver) Reset() {
	solver.VisitedConfigurations = 0
}

//
//func (solver *BruteforceBagSolver) getSolution(problemInstance model.ProblemInstance) model.FinalConfiguration {
//	var bestConfiguration model.FinalConfiguration;
//	bestPrice := 0;
//
//
//	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
//
//	for  {
//		totalWeight := problemInstance.CalculateTotalWeight(configuration);
//		totalPrice := problemInstance.CalculateTotalPrice(configuration);
//		if totalWeight <= problemInstance.Bag.Capacity && totalPrice > bestPrice {
//			bestConfiguration = configuration.clone();
//			bestPrice = totalPrice;
//		}
//		if getNextConfiguration(configuration) {break;}
//	}
//
//
//        return bestConfiguration;
//    }
//
//    func getNextConfiguration(configuration model.PartialConfiguration) bool {
//        for i := 0; i < len(configuration.Flags); i++ {
//            if configuration.Flags[i] {
//                configuration.Flags[i] = false;
//            }else{
//                configuration.Flags[i] = true;
//                return true;
//            }
//        }
//        return false;
//    }
//

func (solver *BruteforceBagSolver) GetSolution(problemInstance model.ProblemInstance) *model.FinalConfiguration {
	configuration := model.MakePartialConfiguration(len(problemInstance.Items))
	return solver.getSolutionRec(problemInstance, configuration)
}

//todo pass all by reference
func (solver *BruteforceBagSolver) getSolutionRec(problemInstance model.ProblemInstance, configuration model.PartialConfiguration) *model.FinalConfiguration {
	if configuration.MaskIndex == len(configuration.Flags) { //All booleans are set.
		solver.VisitedConfigurations++
		totalWeight := problemInstance.CalculateTotalWeight(configuration.Flags)
		totalPrice := problemInstance.CalculateTotalPrice(configuration.Flags)
		if totalWeight <= problemInstance.Bag.Capacity &&
			totalPrice >= problemInstance.MinimumPrice {
			finalConfiguration := make([]bool, len(configuration.Flags))
			copy(finalConfiguration, configuration.Flags)
			return (*model.FinalConfiguration)(&finalConfiguration)
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
