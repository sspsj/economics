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

// test issuance decay layer to layer, expect it to decay by precisely lambda
//func Test_LayerToLayer(t *testing.T) {
//}

// test issuance halving
func Test_Halving(t *testing.T) {
	expectedLayerSubsidyBefore := uint64(TotalSubsidy * (1 - math.Exp(-Lambda)))
	expectedLayerSubsidyAfter :=
		uint64(TotalSubsidy*(1-math.Exp(-Lambda*float64(LayersPerHalving+1)))) -
			uint64(TotalSubsidy*(1-math.Exp(-Lambda*float64(LayersPerHalving))))

	subsidyBefore := TotalSubsidyAtLayer(0, 0)
	subsidyAfter := TotalSubsidyAtLayer(0, LayersPerHalving)

	assert.Equal(t, expectedLayerSubsidyBefore, subsidyBefore)
	assert.Equal(t, expectedLayerSubsidyAfter, subsidyAfter)
	assert.Equal(t, uint64(2), subsidyBefore/subsidyAfter)
}

func Test_IssuanceAtFinalLayer(t *testing.T) {
}
