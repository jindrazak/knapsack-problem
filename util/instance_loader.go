package util

import (
	"bufio"
	"github.com/jindrazak/knapsack-problem/model"
	"log"
	"os"
	"strconv"
	"strings"
)

type InstanceLoader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func CreateInstanceLoader(filename string) InstanceLoader {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return InstanceLoader{file: file, scanner: bufio.NewScanner(file)}
}

func (instanceLoader *InstanceLoader) Close() {
	instanceLoader.file.Close()
}

func (instanceLoader *InstanceLoader) Next() bool {
	return instanceLoader.scanner.Scan()
}

func (instanceLoader *InstanceLoader) Current() model.ProblemInstance {
	return buildProblemInstance(instanceLoader.scanner.Text())
}

func buildProblemInstance(textData string) model.ProblemInstance {
	numbers := strings.Fields(textData)
	id, _ := strconv.Atoi(numbers[0])
	itemsCount, _ := strconv.Atoi(numbers[1])
	capacity, _ := strconv.Atoi(numbers[2])
	minimumPrice, _ := strconv.Atoi(numbers[3])
	items := make([]model.Item, itemsCount)

	for i := 0; i < itemsCount; i++ {
		weight, _ := strconv.Atoi(numbers[(i*2)+4])
		price, _ := strconv.Atoi(numbers[(i*2)+5])
		items[i] = model.Item{Price: price, Weight: weight}
	}
	return model.ProblemInstance{Id: id, Items: items, Bag: model.Bag{Capacity: capacity}, MinimumPrice: minimumPrice}
}
