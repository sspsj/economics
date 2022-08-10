package rewards

import (
	"github.com/ericlagergren/decimal"
	"github.com/spacemeshos/economics/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Lambda(t *testing.T) {
	// https://www.wolframalpha.com/input?i=IntegerPart%5Blog%282%29%2F3265325.445944065837653721294562894*10e39%5D
	expectedLambda, ok := decimal.WithContext(Ctx).SetString("2122750678407627052571300326387314")
	assert.True(t, ok)
	scaleFactor, ok := decimal.WithContext(Ctx).SetString("10e39")
	assert.True(t, ok)
	actualLambda := Ctx.Mul(decimal.WithContext(Ctx), Lambda, scaleFactor)
	assert.Equal(t, 0, actualLambda.Cmp(expectedLambda),
		"expected lambda %d got %d with half life %f", expectedLambda, actualLambda, HalfLife)
}

func Test_Rounding(t *testing.T) {
	// check that lambda multiplier is nonzero, make sure it rounds down to nearest uint
	layerID := uint32(99)
	unrounded := getUnroundedAccumulatedSubsidy(layerID)
	rounded := TotalAccumulatedSubsidyAtLayer(layerID)
	assert.Equal(t, 1, unrounded.Cmp(new(decimal.Big).SetUint64(rounded)),
		"expected subsidy to be rounded")
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
	startLayerID := uint32(9)
	endLayerID := startLayerID + 10

	accumulatedRounddown := decimal.WithContext(Ctx)
	for layerID := startLayerID + 1; layerID <= endLayerID; layerID++ {
		subsidyUnrounded := getUnroundedAccumulatedSubsidy(layerID)
		assert.Equal(t, 1, subsidyUnrounded.Sign(), "expected positive subsidy value")

		subsidyBigIntPart := decimal.WithContext(Ctx).Copy(subsidyUnrounded)
		subsidyBigIntPart.Context.RoundingMode = decimal.ToZero
		subsidyBigIntPart.RoundToInt()

		rounddown := decimal.WithContext(Ctx).Sub(subsidyUnrounded, subsidyBigIntPart)
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

func Test_TenYearIssuance(t *testing.T) {
	// total issuance at ten years should be 600M
	// total subsidy at ten years should be 450M
	expectedTenYearSubsidy := uint64(constants.TenYearTarget - constants.TotalVaulted)

	tenYears := constants.OneYear * 10
	tenYearsU32 := uint32(tenYears)
	tenYearSubsidyRaw := getUnroundedAccumulatedSubsidy(tenYearsU32)
	tenYearSubsidy := TotalAccumulatedSubsidyAtLayer(tenYearsU32)
	assert.Equal(t, 1, tenYearSubsidyRaw.Cmp(new(decimal.Big).SetUint64(tenYearSubsidy)),
		"expected subsidy to be rounded")

	// we expect to be one smidge short due to rounding
	assert.Equal(t, expectedTenYearSubsidy-1, tenYearSubsidy,
		"expected total subsidy of %d at ten years, got %d raw %f",
		expectedTenYearSubsidy-1, tenYearSubsidy, tenYearSubsidyRaw)
}

// test hardcoded subsidy in sampled layers
func Test_Subsidy(t *testing.T) {
	testValues := []struct {
		layerID              uint32
		expectedSubsidyLayer uint64
		expectedSubsidyTotal uint64
	}{
		{0, 0, 0},
		{1, 477618851948, 477618851948},
		{10, 477617939470, 4776183957091},
		{100, 477608814783, 47761383334833},
		{1000, 477517577499, 477568212936019},
		{10000, 476606162709, 4771123282240424},
		{100000, 467587145594, 47258525358305362},
		{1000000, 386270606901, 430330032010906191},
		{10000000, 57171902772, 1980670693994627720},
		{100000000, 288, 2249999998641081355},
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
	// https://www.wolframalpha.com/input?i=-1051200*log%282%29%2Flog%281-%28450%2F2250%29%29
	expectedHalfLife, ok := new(decimal.Big).SetString("3265325.445944065837653721294562894")
	assert.True(t, ok)
	assert.Equal(t, 0, HalfLife.Cmp(expectedHalfLife))

	expectedSubsidyAtHalflife := uint64(constants.TotalSubsidy / 2)
	issuanceMargin := uint64(150 * constants.OneSmesh)

	// half life is not an integer so we cannot test this precisely, just test the nearest layers
	lastLayerBeforeHalflife, ok := decimal.WithContext(Ctx).Copy(HalfLife).RoundToInt().Uint64()
	assert.True(t, ok)
	lastLayerBeforeHalflifeU32 := uint32(lastLayerBeforeHalflife)
	assert.Equal(t, lastLayerBeforeHalflife, uint64(lastLayerBeforeHalflifeU32))

	totalBeforeHalfLife := TotalAccumulatedSubsidyAtLayer(lastLayerBeforeHalflifeU32)

	// expect it to be within the margin
	assert.Less(t, expectedSubsidyAtHalflife-totalBeforeHalfLife, issuanceMargin)
	assert.Positive(t, expectedSubsidyAtHalflife-totalBeforeHalfLife)

	firstLayerAfterHalfLifeU32 := uint32(lastLayerBeforeHalflife + 1)
	assert.Equal(t, lastLayerBeforeHalflife+1, uint64(firstLayerAfterHalfLifeU32))

	totalAfterHalfLife := TotalAccumulatedSubsidyAtLayer(firstLayerAfterHalfLifeU32)

	// expect it to be within the margin
	assert.Less(t, totalAfterHalfLife-expectedSubsidyAtHalflife, issuanceMargin)
	assert.Positive(t, totalAfterHalfLife-expectedSubsidyAtHalflife)
}

// test issuance of final smidge
func Test_FinalLayer(t *testing.T) {
	expectedFinalLayer := 199069359
	finalLayer, ok := FinalLayer.Uint64()
	assert.True(t, ok)
	finalLayerUint32 := uint32(finalLayer)
	assert.Equal(t, finalLayer, uint64(finalLayerUint32))

	// check against hardcoded number
	assert.Equal(t, uint32(expectedFinalLayer), finalLayerUint32,
		"expected final layer %d to be %d", finalLayerUint32, expectedFinalLayer)

	// increment by one to account for rounding down
	finalLayerUint32++

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
