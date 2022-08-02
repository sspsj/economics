package rewards

import "math"
import "github.com/spacemeshos/economics/constants"

var (
	Lambda = math.Log(2) / float64(constants.HalfLife)
)

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	// add one because layers are zero-indexed
	return uint64(constants.TotalSubsidy * (1 - math.Exp(-Lambda*float64(layersAfterEffectiveGenesis+1))))
}

// TotalSubsidyAtLayer returns the total subsidy issued in the layer
func TotalSubsidyAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	subsidyAtLayer := TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis)
	subsidyPrevLayer := uint64(0)
	if layersAfterEffectiveGenesis > 0 {
		subsidyPrevLayer = TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis - 1)
	}

	// Calculate as the difference between the total issuance as of the previous layer and the total issuance as of the
	// current layer
	return subsidyAtLayer - subsidyPrevLayer
}
