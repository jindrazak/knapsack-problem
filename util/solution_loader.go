package util

import (
	"bufio"
	"github.com/jindrazak/knapsack-problem/model"
	"log"
	"os"
	"strconv"
	"strings"
)

type SolutionLoader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func CreateSolutionLoader(filename string) SolutionLoader {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return SolutionLoader{file: file, scanner: bufio.NewScanner(file)}
}

func (SolutionLoader *SolutionLoader) Close() {
	SolutionLoader.file.Close()
}

func (SolutionLoader *SolutionLoader) Next() bool {
	return SolutionLoader.scanner.Scan()
}

func (SolutionLoader *SolutionLoader) Current() model.Solution {
	return buildSolution(SolutionLoader.scanner.Text())
}

func buildSolution(textData string) model.Solution {
	numbers := strings.Fields(textData)
	id, _ := strconv.Atoi(numbers[0])
	itemsCount, _ := strconv.Atoi(numbers[1])
	price, _ := strconv.Atoi(numbers[2])
	itemsPresent := make([]bool, itemsCount)

	for i := 0; i < itemsCount; i++ {
		itemPresentInt, _ := strconv.Atoi(numbers[i+3])
		itemsPresent[i] = itemPresentInt == 1
	}
	return model.Solution{Id: id, Price: price, Configuration: itemsPresent}
}
