package rewards

import "math"

const (
	OneSmesh = 10000000000 // smidge per smesh

	// Total issuance figures

	TotalIssuance = OneSmesh * 2400000000 // 2.4bn smesh
	TotalVaulted  = OneSmesh * 120000000  // 120mn smesh
	TotalSubsidy  = TotalIssuance - TotalVaulted

	// Time-based constants

	LayerTime = 5                  // in minutes
	OneYear   = 365.2425 * 24 * 60 // in minutes
	HalfLife  = OneYear * 29.32233 // exponential decay half-life
)

var (
	// LayersPerHalving is the number of layers after which the per-layer subsidy halves
	LayersPerHalving = uint32(math.Floor(HalfLife / LayerTime))
	Lambda           = math.Log(2) / float64(LayersPerHalving)
)

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(effectiveGenesis uint32, layerID uint32) uint64 {
	if layerID < effectiveGenesis {
		return 0
	}
	effectiveLayer := layerID - effectiveGenesis

	// add one because layers are zero-indexed
	return uint64(TotalSubsidy * (1 - math.Exp(-Lambda*float64(effectiveLayer+1))))
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
