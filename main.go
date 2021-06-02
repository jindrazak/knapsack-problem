package main

import (
	"fmt"
	"github.com/jindrazak/knapsack-problem/algorithm"
	"github.com/jindrazak/knapsack-problem/util"
	"reflect"
	"strconv"
	"time"
)

func main() {
	n := "20"
	setId := "N"
	bruteforceSolver := algorithm.BruteforceBagSolver{}
	solveSet(setId, n, &bruteforceSolver)
	// solveSet(setId, n, BranchBoundBagSolver());
	setId = "Z"
	//solveSet(setId, n, &bruteforceSolver)
	// solveSet(setId, n, new BranchBoundBagSolver());

}

func solveSet(directory string, n string, solver *algorithm.BruteforceBagSolver) {
	instanceLoader := initializeInstanceLoader(directory, n)
	solutionLoader := initializeSolutionLoader(directory, n)

	println("Starting to solve set " + directory + "/" + n + " using " + reflect.TypeOf(solver).Name())

	var times []int64
	for instanceLoader.Next() && solutionLoader.Next() {
		problemInstance := instanceLoader.Current()
		solution := solutionLoader.Current()

		solver.Reset()

		start := time.Now()
		calculatedSolution := solver.GetSolution(problemInstance)
		timeElapsed := time.Since(start)
		times = append(times, timeElapsed.Nanoseconds())
		hasSolution := solution.IsSolvable(problemInstance)
		foundSolution := calculatedSolution != nil

		if hasSolution && foundSolution {
			if !reflect.DeepEqual(solution.Configuration, *calculatedSolution) {
				println("Error. Solution does not match calculatedSolution.")
				fmt.Printf("Expected: %v\n", solution.Configuration)
				fmt.Printf("Got: %v\n", *calculatedSolution)
			}
		} else if !hasSolution && foundSolution {
			println("Error. Found invalid solution.")
		} else if hasSolution && !foundSolution {
			println("Error. Not found existing solution.")
		}

		var resultString string
		if calculatedSolution != nil {
			resultString = "solvable"
		} else {
			resultString = "not solvable"
		}
		generalInfo := "Instance '" + strconv.Itoa(problemInstance.Id) + "' solved. "
		resultInfo := "Result: " + resultString + ". "
		visitedStatesInfo := "Visited States: " + strconv.Itoa(solver.VisitedConfigurations) + ". "
		elapsedTimeInfo := "Elapsed time: " + strconv.FormatInt(timeElapsed.Nanoseconds(), 10) + " ns"
		println(generalInfo + resultInfo + visitedStatesInfo + elapsedTimeInfo)
	}
	var total int64
	for _, nanoseconds := range times {
		total += nanoseconds
	}
	average := total / int64(len(times))
	println("Average duration:" + strconv.FormatInt(average, 10))
	println("Set '" + directory + "/" + n + "' completed using '" + reflect.TypeOf(*solver).String())
}

func initializeInstanceLoader(setId string, n string) util.InstanceLoader {
	filename := "instances/" + setId + "R/" + setId + "R" + n + "_inst.dat"
	return util.CreateInstanceLoader(filename)
}

func initializeSolutionLoader(setId string, n string) util.SolutionLoader {
	filename := "instances/" + setId + "R/" + setId + "K" + n + "_sol.dat"
	return util.CreateSolutionLoader(filename)
}
