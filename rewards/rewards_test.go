package rewards

import (
	"github.com/ericlagergren/decimal"
	"github.com/spacemeshos/economics/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test issuance halving
func Test_Halving(t *testing.T) {

	// subtract one because we zero index layers
	totalAtHalfLife := TotalAccumulatedSubsidyAtLayer(constants.HalfLife - 1)

	assert.Equal(t, 0, totalAtHalfLife.Mul(totalAtHalfLife, new(decimal.Big).SetUint64(2)).Cmp(new(decimal.Big).SetUint64(constants.TotalSubsidy)))
}

func Test_IssuanceAtFinalLayer(t *testing.T) {
}
