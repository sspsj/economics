package rewards

import (
	"github.com/ericlagergren/decimal"
	"github.com/spacemeshos/economics/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Rounding(t *testing.T) {
	// check that lambda multiplier is nonzero, make sure it rounds down to nearest uint
	layerID := uint32(99)
	unrounded := getUnroundedAccumulatedSubsidy(layerID)
	rounded := TotalAccumulatedSubsidyAtLayer(layerID)
	assert.NotEqual(t, unrounded, rounded, "expected subsidy to be rounded")
	unroundedFloat, _ := unrounded.Float64()
	assert.Greater(t, unroundedFloat, float64(rounded))
	unroundedRounded, _ := unrounded.Uint64()
	assert.Equal(t, unroundedRounded, rounded, "expected rounding of unrounded to equal rounded value")
}

func Test_Accumulation(t *testing.T) {
	// check that layer to layer values get added to accumulated total
	layerID := uint32(99)
	totalAtLayer := TotalAccumulatedSubsidyAtLayer(layerID)
	nextLayerSubsidy := TotalSubsidyAtLayer(layerID + 1)
	totalNextLayer := TotalAccumulatedSubsidyAtLayer(layerID + 1)
	assert.Equal(t, totalAtLayer+nextLayerSubsidy, totalNextLayer)
}

func Test_Rounddown(t *testing.T) {
	// in low layer ranges we expect round down at each layer
	startLayerID := uint32(99)
	endLayerID := startLayerID + 10

	accumulatedRounddown := new(decimal.Big)
	for layerID := startLayerID + 1; layerID <= endLayerID; layerID++ {
		subsidyUnrounded := getUnroundedAccumulatedSubsidy(layerID)
		assert.Equal(t, 1, subsidyUnrounded.Sign(), "expected positive subsidy value")

		subsidyBigIntPart := new(decimal.Big).Copy(subsidyUnrounded)
		subsidyBigIntPart.Context.RoundingMode = decimal.ToZero
		subsidyBigIntPart.RoundToInt()

		rounddown := new(decimal.Big).Sub(subsidyUnrounded, subsidyBigIntPart)
		assert.Equal(t, 1, rounddown.Sign(), "expected positive rounddown value")
		rounddownFloat, ok := rounddown.Float64()
		assert.True(t, ok)

		// note: we expect there to _always_ be a rounddown
		assert.Greater(t, rounddownFloat, float64(0), "expected rounddown between zero and one")
		assert.Less(t, rounddownFloat, float64(1), "expected rounddown between zero and one")
		accumulatedRounddown.Add(accumulatedRounddown, rounddown)
	}

	// expect substantial total rounddown
	// float not uint since we want to check that amount before rounding > 1
	accumulatedRounddownFloat, ok := accumulatedRounddown.Float64()
	assert.True(t, ok)
	assert.Greater(t, accumulatedRounddownFloat, float64(1), "expected rounddown > 1 smidge")
	accumulatedRounddownUint, ok := accumulatedRounddown.Uint64()
	assert.True(t, ok)
	assert.Greater(t, accumulatedRounddownUint, uint64(1), "expected rounddown > 1 smidge")

	// expect total issuance for final layer to equal per-layer issuance plus accumulated rounddown (rounded down!)
	//totalDifferential, ok := accumulatedDifferential.Uint64()
	//assert.True(t, ok)
	//assert.Equalf(t, issuanceAtStart+totalDifferential+accumulatedRounddownUint, issuanceAtEnd,
	//	"expected start issuance %d plus accumulated differential %d plus rounddown uint part %d (total %d) "+
	//		"to equal final layer issuance %d",
	//	issuanceAtStart, totalDifferential, accumulatedRounddownUint, issuanceAtStart+totalDifferential+accumulatedRounddownUint,
	//	issuanceAtEnd)
}

// test hardcoded issuance in first layer
func Test_IssuanceAtFirstLayer(t *testing.T) {

}

// test issuance halving
func Test_Halving(t *testing.T) {
	// subtract one because we zero index layers
	totalAtHalfLife := TotalAccumulatedSubsidyAtLayer(constants.HalfLife - 1)

	// add hardcoded value, check single layer issuance

	//assert.Equal(t, 0, totalAtHalfLife.Mul(totalAtHalfLife, new(decimal.Big).SetUint64(2)).Cmp(new(decimal.Big).SetUint64(constants.TotalSubsidy)))
	assert.Equal(t, totalAtHalfLife*2, constants.TotalSubsidy)
}

// test issuance of final smidge
func Test_IssuanceAtFinalLayer(t *testing.T) {
}
