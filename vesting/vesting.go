package vesting

import (
	"log"

	"github.com/spacemeshos/economics/constants"
)

func AccumulatedVestAtLayer(layersAfterGenesis uint32) uint64 {
	if layersAfterGenesis < constants.VestStart {
		return 0
	} else if layersAfterGenesis >= constants.VestEnd {
		return constants.TotalVaulted
	}

	// Note: this rounds down to the nearest int number of smidge below the intended vest as of the input layer.
	// No need to check for overflow on the subtraction but we can overflow on the multiplication.
	numLayers := uint64(layersAfterGenesis - constants.VestStart)
	vest := constants.VestPerLayer * numLayers
	if vest/constants.VestPerLayer != numLayers {
		log.Fatal("integer overflow")
	}
	return constants.VestedAtCliff + vest
}

func VestAtLayer(layersAfterGenesis uint32) uint64 {
	// base case: no vesting before vest start, no vesting after vest end
	if layersAfterGenesis < constants.VestStart {
		return 0
	} else if layersAfterGenesis > constants.VestEnd {
		return 0
	}

	// vest as of the previous layer
	var prevLayerAccumulatedVest, curLayerAccumulatedVest uint64
	if layersAfterGenesis > 0 {
		prevLayerAccumulatedVest = AccumulatedVestAtLayer(layersAfterGenesis - 1)
	}

	// intended vest as of this layer
	curLayerAccumulatedVest = AccumulatedVestAtLayer(layersAfterGenesis)

	// return the difference
	return curLayerAccumulatedVest - prevLayerAccumulatedVest
}
