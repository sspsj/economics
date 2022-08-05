package constants

const (
	OneSmesh = 1000000000 // 1e9 (1bn) smidge per smesh

	// Vaults and vesting

	TotalVaulted  = OneSmesh * 150000000 // 150mn smesh
	CliffRatio    = 0.25                 // pct vest at cliff
	VestStart     = 105190               // one year, in layers
	VestEnd       = 4 * VestStart        // four years post-genesis, three years post-vesting start
	VestLayers    = VestEnd - VestStart  // three years, in layers
	VestedAtCliff = uint64(CliffRatio * TotalVaulted)

	// VestPerLayer is rounded down to the nearest int. We make up for this rounding in the code.
	VestPerLayer = (TotalVaulted - VestedAtCliff) / VestLayers

	// Total issuance figures

	TotalIssuance = OneSmesh * 2400000000 // 2.4bn smesh
	TotalSubsidy  = TotalIssuance - TotalVaulted

	HalfLife = 3267565 // in layers, ~31 years, s.t. total issuance at 10 years is ~600m
)
