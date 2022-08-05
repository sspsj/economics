package rewards

import (
	"github.com/ericlagergren/decimal"
	"github.com/spacemeshos/economics/constants"
)

var (
	Ctx      = decimal.Context128
	Two      = decimal.New(2, 0)
	HalfLife = decimal.New(constants.HalfLife, 0)
	M        = new(decimal.Big)
	LogTwo   = Ctx.Log(M, Two)
	Lambda   = LogTwo.Quo(LogTwo, HalfLife)
)

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis uint32) *decimal.Big {
	// add one because layers are zero-indexed
	totalSubsidy := new(decimal.Big).SetUint64(constants.TotalSubsidy)
	layerCount := new(decimal.Big).SetUint64(uint64(layersAfterEffectiveGenesis + 1))
	negLambda := new(decimal.Big).Copy(Lambda).Neg(Lambda)
	expInner := new(decimal.Big).Mul(negLambda, layerCount)
	expOuter := new(decimal.Big)
	Ctx.Exp(expOuter, expInner)
	one := decimal.New(1, 0)
	supplyMultiplier := one.Sub(one, expOuter)
	return totalSubsidy.Mul(totalSubsidy, supplyMultiplier)
}

// TotalSubsidyAtLayer returns the total subsidy issued in the layer
func TotalSubsidyAtLayer(layersAfterEffectiveGenesis uint32) *decimal.Big {
	subsidyAtLayer := TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis)
	subsidyPrevLayer := new(decimal.Big)
	if layersAfterEffectiveGenesis > 0 {
		subsidyPrevLayer = TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis - 1)
	}

	// Calculate as the difference between the total issuance as of the previous layer and the total issuance as of the
	// current layer
	return subsidyAtLayer.Sub(subsidyAtLayer, subsidyPrevLayer)
}
