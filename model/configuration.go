package model

type FinalConfiguration []bool

type PartialConfiguration struct {
	Flags     FinalConfiguration
	MaskIndex int
}

func MakePartialConfiguration(size int) PartialConfiguration {
	return PartialConfiguration{Flags: make([]bool, size), MaskIndex: 0}
}

func (configuration *PartialConfiguration) SetNextFlag(newFlag bool) {
	configuration.Flags[configuration.MaskIndex] = newFlag
	configuration.MaskIndex++
}
