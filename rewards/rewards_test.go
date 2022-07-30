package rewards

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_IssuanceBeforeEffectiveGenesis(t *testing.T) {
	subsidy := TotalAccumulatedSubsidyAtLayer(1, 0)
	totalSubsidy := TotalSubsidyAtLayer(1, 0)
	assert.Zero(t, subsidy)
	assert.Zero(t, totalSubsidy)
}

func Test_IssuanceAtEffectiveGenesis(t *testing.T) {
	expectedIssuance := uint64(TotalSubsidy * (1 - math.Exp(-Lambda)))

	subsidy := TotalAccumulatedSubsidyAtLayer(1, 1)
	totalSubsidy := TotalSubsidyAtLayer(1, 1)
	assert.Equal(t, expectedIssuance, subsidy)
	assert.Equal(t, expectedIssuance, totalSubsidy)
}

// test layer-to-layer decay in total issuance
func Test_LayerToLayer(t *testing.T) {
	expectedTotalBefore := uint64(TotalSubsidy * (1 - math.Exp(-Lambda)))
	expectedTotalAfter := uint64(TotalSubsidy * (1 - math.Exp(-Lambda*2)))

	subsidyBefore := TotalAccumulatedSubsidyAtLayer(0, 0)
	subsidyAfter := TotalAccumulatedSubsidyAtLayer(0, 1)

	assert.Equal(t, expectedTotalBefore, subsidyBefore)
	assert.Equal(t, expectedTotalAfter, subsidyAfter)
}

// test issuance halving
func Test_Halving(t *testing.T) {
	expectedLayerSubsidyBefore := uint64(TotalSubsidy * (1 - math.Exp(-Lambda)))
	expectedLayerSubsidyAfter :=
		uint64(TotalSubsidy*(1-math.Exp(-Lambda*float64(HalfLife+1)))) -
			uint64(TotalSubsidy*(1-math.Exp(-Lambda*float64(HalfLife))))

	subsidyBefore := TotalSubsidyAtLayer(0, 0)
	subsidyAfter := TotalSubsidyAtLayer(0, HalfLife)

	assert.Equal(t, expectedLayerSubsidyBefore, subsidyBefore)
	assert.Equal(t, expectedLayerSubsidyAfter, subsidyAfter)
	assert.Equal(t, uint64(2), subsidyBefore/subsidyAfter)
}

func Test_IssuanceAtFinalLayer(t *testing.T) {
}
