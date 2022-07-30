package vesting

import "github.com/spacemeshos/economics/constants"

func AccumulatedVestAtLayer(effectiveGenesis uint32, layerID uint32) uint64 {
	if layerID < effectiveGenesis {
		return 0
	}
	effectiveLayer := layerID - effectiveGenesis
	if effectiveLayer < constants.VestStart {
		return 0
	} else if effectiveLayer > constants.VestEnd {
		return constants.TotalVaulted
	}
	return constants.VestPerLayer * uint64(effectiveLayer-constants.VestStart)
}

func VestAtLayer(effectiveGenesis uint32, layerID uint32) uint64 {
	if layerID < effectiveGenesis {
		return 0
	}
	effectiveLayer := layerID - effectiveGenesis
	if effectiveLayer < constants.VestStart {
		return 0
	} else if effectiveLayer > constants.VestEnd {
		return 0
	}
	return constants.VestPerLayer
}
