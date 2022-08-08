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
	// note that there's no way to test the accuracy of the rounddown since the issuance at any given layer is
	// simply defined to be the difference between the total subsidy at the previous layer and the total subsidy
	// at the layer
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
}

// test hardcoded subsidy in sampled layers
func Test_Subsidy(t *testing.T) {
	testValues := []struct {
		layerID              uint32
		expectedSubsidyLayer uint64
		expectedSubsidyTotal uint64
	}{
		{0, 477291497137, 477291497137},
		{10, 477290484662, 5250200899893},
		{100, 477281372481, 48205929913896},
		{1000, 477190260232, 477718117773694},
		{10000, 476280093827, 4768332952709150},
		{100000, 467273365577, 47226944479755904},
		{1000000, 386061943727, 430065653103871634},
		{10000000, 57215889114, 1980278615961832312},
		{100000000, 293, 2249999998621166071},
		{1000000000, 0, 2250000000000000000},
	}
	for _, testTuple := range testValues {
		subsidyLayer := TotalSubsidyAtLayer(testTuple.layerID)
		subsidyTotal := TotalAccumulatedSubsidyAtLayer(testTuple.layerID)
		assert.Equal(t, testTuple.expectedSubsidyLayer, subsidyLayer,
			"expected layer %d subsidy %d to equal %d", testTuple.layerID, subsidyLayer, testTuple.expectedSubsidyLayer)
		assert.Equal(t, testTuple.expectedSubsidyTotal, subsidyTotal,
			"expected layer %d total subsidy %d to equal %d", testTuple.layerID, subsidyTotal, testTuple.expectedSubsidyTotal)
	}
}

// test issuance halving
func Test_Halving(t *testing.T) {
	// subtract one because we zero index layers
	totalAtHalfLife := TotalAccumulatedSubsidyAtLayer(constants.HalfLife - 1)
	assert.Equal(t, uint64(constants.TotalSubsidy), totalAtHalfLife*2,
		"expected total accum. subsidy %d at halfway point %d to be half of total subsidy %d",
		totalAtHalfLife, constants.HalfLife, constants.TotalSubsidy)
}

// test issuance of final smidge
func Test_FinalLayer(t *testing.T) {
	finalLayer, ok := FinalLayer.Uint64()
	assert.True(t, ok)
	finalLayerUint32 := uint32(finalLayer)
	assert.Equal(t, finalLayer, uint64(finalLayerUint32))

	// check against hardcoded number
	assert.Equal(t, uint32(199205893), finalLayerUint32)

	// that final smidge will never be issued since, beyond this point, all issuance will be rounded down to zero
	expectedFinalTotalIssuance := uint64(constants.TotalSubsidy) - 1
	subsidyLayer := TotalSubsidyAtLayer(finalLayerUint32)
	subsidyTotal := TotalAccumulatedSubsidyAtLayer(finalLayerUint32)
	assert.Equal(t, uint64(1), subsidyLayer,
		"expected final layer %d subsidy %d to equal %d", finalLayerUint32, subsidyLayer, 1)
	assert.Equal(t, expectedFinalTotalIssuance, subsidyTotal,
		"expected final layer %d total subsidy %d to equal %d", finalLayerUint32, subsidyTotal, expectedFinalTotalIssuance)

	// one layer later we expect issuance to go to zero
	subsidyLayerBeyond := TotalSubsidyAtLayer(finalLayerUint32 + 1)
	subsidyTotalBeyond := TotalAccumulatedSubsidyAtLayer(finalLayerUint32 + 1)
	assert.Equal(t, uint64(0), subsidyLayerBeyond,
		"expected final layer +1 %d subsidy %d to equal %d", finalLayerUint32+1, subsidyLayerBeyond, 0)
	assert.Equal(t, expectedFinalTotalIssuance, subsidyTotalBeyond,
		"expected final layer +1 %d total subsidy %d to equal %d", finalLayerUint32+1, subsidyTotalBeyond, expectedFinalTotalIssuance)
}
