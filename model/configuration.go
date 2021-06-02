package model

type FinalConfiguration []bool

type PartialConfiguration struct {
	Flags     FinalConfiguration
	MaskIndex int
}

func MakePartialConfiguration(size int) PartialConfiguration {
	return PartialConfiguration{Flags: make([]bool, size), MaskIndex: 0}
}

func (configuration FinalConfiguration) Clone() FinalConfiguration {
	copiedConfiguration := make([]bool, len(configuration))
	copy(copiedConfiguration, configuration)
	return copiedConfiguration
}

func (configuration PartialConfiguration) Clone() PartialConfiguration {
	copiedFlags := configuration.Flags.Clone()
	return PartialConfiguration{Flags: copiedFlags, MaskIndex: configuration.MaskIndex}
}

func (configuration *PartialConfiguration) SetNextFlag(newFlag bool) {
	configuration.Flags[configuration.MaskIndex] = newFlag
	configuration.MaskIndex++
}
