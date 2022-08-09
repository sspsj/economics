package vesting

import (
	"github.com/spacemeshos/economics/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Genesis(t *testing.T) {
	// we expect zero vest in genesis layer
	assert.Zero(t, AccumulatedVestAtLayer(0),
		"expected zero vest at genesis")
	assert.Zero(t, VestAtLayer(0),
		"expected zero vest at genesis")
}

func Test_Cliff(t *testing.T) {
	// expect zero vest before cliff
	assert.Zero(t, AccumulatedVestAtLayer(constants.VestStart-1),
		"expected zero vest before cliff")
	assert.Zero(t, VestAtLayer(constants.VestStart-1),
		"expected zero vest before cliff")

	// expect only cliff amount at vest start
	assert.Equal(t, constants.VestedAtCliff, AccumulatedVestAtLayer(constants.VestStart),
		"expected cliff at vest start")
	assert.Equal(t, constants.VestedAtCliff, VestAtLayer(constants.VestStart),
		"expected cliff at vest start")

	// expect regular issuance to begin after cliff
	assert.Equal(t, constants.VestedAtCliff+constants.VestPerLayer, AccumulatedVestAtLayer(constants.VestStart+1),
		"expected cliff at vest start")
	assert.Equal(t, constants.VestPerLayer, VestAtLayer(constants.VestStart+1),
		"expected cliff at vest start")

	// test rounding in VestPerLayer
	actualTotalVested := constants.VestPerLayer * constants.VestLayers
	expectedTotalVested := constants.TotalVaulted - constants.VestedAtCliff
	vestingDifferential := expectedTotalVested - actualTotalVested

	// it should be less, but the difference should be less than one smesh
	assert.Greater(t, vestingDifferential, uint64(0), "expected per layer vesting to be rounded down")
	assert.Less(t, vestingDifferential, uint64(constants.OneSmesh), "expected total vest differential to be < 1 SMESH")

	// test end of vesting
	assert.Equal(t, constants.VestedAtCliff+constants.VestPerLayer*(constants.VestLayers-1), AccumulatedVestAtLayer(constants.VestEnd-1),
		"unexpected accumulated vest in penultimate vesting layer")
	assert.Equal(t, constants.VestPerLayer, VestAtLayer(constants.VestEnd-1),
		"unexpected layer vest in final penultimate layer")

	// this is the final layer of vesting, and the above shortfall should also be accounted for here
	assert.Equal(t, uint64(constants.TotalVaulted), AccumulatedVestAtLayer(constants.VestEnd),
		"unexpected accumulated vest in final vesting layer")
	assert.Equal(t, constants.VestPerLayer+vestingDifferential, VestAtLayer(constants.VestEnd),
		"unexpected layer vest in final vesting layer")

	// expect no additional vesting beyond this layer
	assert.Equal(t, uint64(constants.TotalVaulted), AccumulatedVestAtLayer(constants.VestEnd+1),
		"unexpected additional accumulated vesting after final vesting layer")
	assert.Zero(t, VestAtLayer(constants.VestEnd+1),
		"unexpected additional layer vesting after final vesting layer")
}
