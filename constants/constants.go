package constants

const (
	OneSmesh = 1000000000 // 1e9 (1bn) smidge per smesh

	// Vaults and vesting

	// one decade after genesis 600M, 25% of that in vault
	// vesting: one year cliff, three years after that linear (per layer)

	TotalVaulted = OneSmesh * 150000000 // 150mn smesh
	VestStart    = 105190               // one year, in layers
	VestEnd      = 4 * VestStart        // four years post-genesis, three years post-vesting start
	VestLayers   = VestEnd - VestStart  // three years, in layers

	// VestPerLayer is rounded down to the nearest int. We make up for this rounding in the code.
	VestPerLayer = TotalVaulted / VestLayers

	// Total issuance figures

	TotalIssuance = OneSmesh * 2400000000 // 2.4bn smesh
	TotalSubsidy  = TotalIssuance - TotalVaulted

	HalfLife = 3267565 // in layers, ~30 years
)
