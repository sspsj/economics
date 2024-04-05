package constants

const (
	OneSmesh = 1000000000 // 1e9 (1bn) smidge per smesh

	// Vaults and vesting

	TotalVaulted  = OneSmesh * 150000000 // 150mn smesh
	CliffRatio    = 0                    // pct vest at cliff
	OneEpoch      = 4032                 // mainnet epochs are two weeks long
	OneYear       = 105120               // 365 days, in 5-minute intervals
	VestStart     = OneYear              // one year, in layers
	VestEnd       = 4 * OneYear          // four years post-genesis, three years post-vesting start
	VestLayers    = VestEnd - VestStart  // three years, in layers (exclusive of start layer, inclusive of end layer)
	VestedAtCliff = uint64(CliffRatio * TotalVaulted)

	// VestPerLayer is rounded down to the nearest int. We make up for this rounding in the code.
	VestPerLayer = (TotalVaulted - VestedAtCliff) / VestLayers

	// Total issuance figures

	TenYearTarget = OneSmesh * 600000000
	TotalIssuance = OneSmesh * 2400000000 // 2.4bn smesh
	TotalSubsidy  = TotalIssuance - TotalVaulted
)
