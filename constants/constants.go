package constants

const (
	OneSmesh = 10000000000 // smidge per smesh

	// Vaults and vesting

	TotalVaulted = OneSmesh * 120000000 // 120mn smesh
	VestStart    = 525950               // one year, in layers
	VestEnd      = 4 * VestStart
	VestLayers   = VestEnd - VestStart
	VestPerLayer = TotalVaulted / VestLayers

	// Total issuance figures

	TotalIssuance = OneSmesh * 2400000000 // 2.4bn smesh
	TotalSubsidy  = TotalIssuance - TotalVaulted

	HalfLife = 3100000 // in layers, ~30 years
)
