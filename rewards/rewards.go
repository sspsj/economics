package rewards

import "math"
import "github.com/spacemeshos/economics/constants"

var (
	Lambda = math.Log(2) / float64(constants.HalfLife)
)

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(effectiveGenesis uint32, layerID uint32) uint64 {
	if layerID < effectiveGenesis {
		return 0
	}
	effectiveLayer := layerID - effectiveGenesis

	// add one because layers are zero-indexed
	return uint64(constants.TotalSubsidy * (1 - math.Exp(-Lambda*float64(effectiveLayer+1))))
}

// TotalSubsidyAtLayer returns the total subsidy issued in the layer
func TotalSubsidyAtLayer(effectiveGenesis uint32, layerID uint32) uint64 {
	subsidyAtLayer := TotalAccumulatedSubsidyAtLayer(effectiveGenesis, layerID)
	subsidyPrevLayer := uint64(0)
	if layerID > 0 {
		subsidyPrevLayer = TotalAccumulatedSubsidyAtLayer(effectiveGenesis, layerID-1)
	}

	// Calculate as the difference between the total issuance as of the previous layer and the total issuance as of the
	// current layer
	return subsidyAtLayer - subsidyPrevLayer
}
