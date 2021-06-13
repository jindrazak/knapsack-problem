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
	solveSet(setId, n, algorithm.BruteforceBagSolver{})
}

func solveSet(directory string, n string, algorithm algorithm.BagSolver) {
	instanceLoader := initializeInstanceLoader(directory, n)
	solutionLoader := initializeSolutionLoader(directory, n)

	println("Starting to solve set " + directory + "/" + n)

	var durations []int64
	for instanceLoader.Next() && solutionLoader.Next() {
		problemInstance := instanceLoader.Current()
		solution := solutionLoader.Current()

		start := time.Now()
		calculatedSolution := algorithm.CalculateSolution(problemInstance)
		timeElapsed := time.Since(start)
		durations = append(durations, timeElapsed.Nanoseconds())
		hasSolution := solution.IsSolvable(problemInstance)
		foundSolution := calculatedSolution.Configuration != nil

		if hasSolution && foundSolution {
			if !reflect.DeepEqual(solution.Configuration, *(calculatedSolution.Configuration)) {
				fmt.Printf("Expected: %v\n", solution.Configuration)
				fmt.Printf("Got: %v\n", *(calculatedSolution.Configuration))
				panic("Error. Solution does not match calculatedSolution.")
			}
		} else if !hasSolution && foundSolution {
			panic("Error. Found invalid solution.")
		} else if hasSolution && !foundSolution {
			panic("Error. Not found existing solution.")
		}

		var resultString string
		if calculatedSolution.Configuration != nil {
			resultString = "SOLVABLE"
		} else {
			resultString = "NOT SOLVABLE"
		}
		generalInfo := "Instance '" + strconv.Itoa(problemInstance.Id) + "' solved. "
		resultInfo := "Result: " + resultString + ". "
		elapsedTimeInfo := "Elapsed time: " + strconv.FormatInt(timeElapsed.Nanoseconds(), 10) + " ns. "
		visitedConfigurationsInfo := "Visited configurations: " + strconv.Itoa(calculatedSolution.VisitedConfigurations) + ". "

		println(generalInfo + visitedConfigurationsInfo + elapsedTimeInfo + resultInfo)
	}

	println("Average duration:" + strconv.FormatInt(average(durations), 10))
	println("Set '" + directory + "/" + n + "' completed")
}

func average(series []int64) int64 {
	var total int64
	for _, nanoseconds := range series {
		total += nanoseconds
	}
	return total / int64(len(series))
}

func initializeInstanceLoader(setId string, n string) util.InstanceLoader {
	filename := "instances/" + setId + "R/" + setId + "R" + n + "_inst.dat"
	return util.CreateInstanceLoader(filename)
}

func initializeSolutionLoader(setId string, n string) util.SolutionLoader {
	filename := "instances/" + setId + "R/" + setId + "K" + n + "_sol.dat"
	return util.CreateSolutionLoader(filename)
}
