package vesting

import (
	"github.com/spacemeshos/economics/constants"
	"log"
)

func AccumulatedVestAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	if layersAfterEffectiveGenesis < constants.VestStart {
		return 0
	} else if layersAfterEffectiveGenesis >= constants.VestEnd {
		return constants.TotalVaulted
	}

	// Note: this rounds down to the nearest int number of smidge below the intended vest as of the input layer.
	// No need to check for overflow on the subtraction but we can overflow on the multiplication.
	numLayers := uint64(layersAfterEffectiveGenesis - constants.VestStart)
	vest := constants.VestPerLayer * numLayers
	if vest/constants.VestPerLayer != numLayers {
		log.Fatal("integer overflow")
	}
	return constants.VestedAtCliff + vest
}

func VestAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	// base case: no vesting before vest start, no vesting after vest end
	if layersAfterEffectiveGenesis < constants.VestStart {
		return 0
	} else if layersAfterEffectiveGenesis > constants.VestEnd {
		return 0
	}

	// vest as of the previous layer
	var prevLayerAccumulatedVest, curLayerAccumulatedVest uint64
	if layersAfterEffectiveGenesis > 0 {
		prevLayerAccumulatedVest = AccumulatedVestAtLayer(layersAfterEffectiveGenesis - 1)
	}

	// intended vest as of this layer
	curLayerAccumulatedVest = AccumulatedVestAtLayer(layersAfterEffectiveGenesis)

	// return the difference
	return curLayerAccumulatedVest - prevLayerAccumulatedVest
}
