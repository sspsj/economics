package rewards

import (
	"github.com/spacemeshos/economics/constants"

	"github.com/ericlagergren/decimal"
	"log"
)

var (
	Ctx    = decimal.Context128
	One    = decimal.WithContext(Ctx).SetUint64(1)
	LogTwo = Ctx.Log(decimal.WithContext(Ctx), decimal.WithContext(Ctx).SetUint64(2))

	// TenYears contains one extra layer to account for the effective genesis (zero) layer.
	TenYears          = decimal.WithContext(Ctx).SetUint64(10*constants.OneYear + 1)
	IssuanceNum       = decimal.WithContext(Ctx).SetUint64(constants.TenYearTarget - constants.TotalVaulted)
	IssuanceDenom     = decimal.WithContext(Ctx).SetUint64(constants.TotalSubsidy)
	IssuanceFrac      = Ctx.Sub(decimal.WithContext(Ctx), One, Ctx.Quo(decimal.WithContext(Ctx), IssuanceNum, IssuanceDenom))
	HalfLife          = Ctx.Mul(decimal.WithContext(Ctx), decimal.WithContext(Ctx).Neg(TenYears), Ctx.Quo(decimal.WithContext(Ctx), LogTwo, Ctx.Log(decimal.WithContext(Ctx), IssuanceFrac)))
	Lambda            = Ctx.Quo(decimal.WithContext(Ctx), LogTwo, HalfLife)
	NegLambda         = decimal.WithContext(Ctx).Neg(Lambda)
	TotalSubsidy      = decimal.WithContext(Ctx).SetUint64(constants.TotalSubsidy)
	FinalIssuanceFrac = Ctx.Quo(decimal.WithContext(Ctx), Ctx.Sub(decimal.WithContext(Ctx), TotalSubsidy, One), TotalSubsidy)
	FinalLayer        = Ctx.Quo(decimal.WithContext(Ctx), Ctx.Log(decimal.WithContext(Ctx), Ctx.Sub(decimal.WithContext(Ctx), One, FinalIssuanceFrac)), NegLambda)
)

func getUnroundedAccumulatedSubsidy(layersAfterEffectiveGenesis uint32) *decimal.Big {
	// add one because layers are zero-indexed and we want > 0 issuance in the first effective genesis layer
	layerCount := decimal.WithContext(Ctx).SetUint64(uint64(layersAfterEffectiveGenesis + 1))
	expInner := Ctx.Mul(decimal.WithContext(Ctx), NegLambda, layerCount)
	expOuter := Ctx.Exp(decimal.WithContext(Ctx), expInner)
	supplyMultiplier := Ctx.Sub(decimal.WithContext(Ctx), One, expOuter)
	return Ctx.Mul(decimal.WithContext(Ctx), TotalSubsidy, supplyMultiplier)
}

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	unroundedSubsidy := getUnroundedAccumulatedSubsidy(layersAfterEffectiveGenesis)
	if ret, ok := unroundedSubsidy.Uint64(); !ok {
		log.Fatal("unable to convert subsidy to uint")
		return 0
	} else {
		return ret
	}
}

// TotalSubsidyAtLayer returns the total subsidy issued in the layer
func TotalSubsidyAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	subsidyAtLayer := TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis)
	var subsidyPrevLayer uint64
	if layersAfterEffectiveGenesis > 0 {
		subsidyPrevLayer = TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis - 1)
	}

	// Calculate as the difference between the total issuance as of the previous layer and the total issuance as of the
	// current layer
	return subsidyAtLayer - subsidyPrevLayer
}
