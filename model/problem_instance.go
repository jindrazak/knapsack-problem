package model

type ProblemInstance struct {
	Id           int
	Items        []Item
	Bag          Bag
	MinimumPrice int
}

func (problemInstance ProblemInstance) CalculateTotalPrice(configuration FinalConfiguration) int {
	totalPrice := 0
	for i := 0; i < len(configuration); i++ {
		if configuration[i] {
			totalPrice += problemInstance.Items[i].Price
		}
	}
	return totalPrice
}

func (problemInstance ProblemInstance) CalculateTotalWeight(configuration FinalConfiguration) int {
	totalWeight := 0
	for i := 0; i < len(configuration); i++ {
		if configuration[i] {
			totalWeight += problemInstance.Items[i].Weight
		}
	}
	return totalWeight
}

func (problemInstance ProblemInstance) CalculateMaxPossiblePrice(configuration PartialConfiguration) int {
	totalPrice := 0
	for i := 0; i < configuration.MaskIndex; i++ {
		if configuration.Flags[i] == true {
			totalPrice += problemInstance.Items[i].Price
		}
	}
	for i := configuration.MaskIndex; i < len(configuration.Flags); i++ {
		totalPrice += problemInstance.Items[i].Price
	}
	return totalPrice
}