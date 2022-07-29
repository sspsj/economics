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
	Lambda = math.Log(2) / (HalfLife / LayerTime)
)

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(layerID uint32) uint64 {
	return uint64(TotalSubsidy * (1 - math.Exp(-Lambda*float64(layerID+2))))
}

// TotalSubsidyAtLayer returns the total subsidy issued in the layer
func TotalSubsidyAtLayer(layerID uint32) uint64 {
	// Calculate as the difference between the total issuance as of the previous layer and the total issuance as of the
	// current layer
	return TotalAccumulatedSubsidyAtLayer(layerID) - TotalAccumulatedSubsidyAtLayer(layerID-1)
}
