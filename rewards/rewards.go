package rewards

import (
	"github.com/ericlagergren/decimal"
	"github.com/spacemeshos/economics/constants"
	"log"
)

var (
	Ctx               = decimal.Context128
	One               = decimal.New(1, 0)
	Two               = decimal.New(2, 0)
	HalfLife          = decimal.New(constants.HalfLife, 0)
	LogTwo            = Ctx.Log(new(decimal.Big), Two)
	Lambda            = Ctx.Quo(new(decimal.Big), LogTwo, HalfLife)
	NegLambda         = decimal.WithContext(Ctx).Copy(Lambda).Neg(Lambda)
	TotalSubsidy      = decimal.WithContext(Ctx).SetUint64(constants.TotalSubsidy)
	FinalIssuanceFrac = Ctx.Quo(new(decimal.Big), Ctx.Sub(new(decimal.Big), TotalSubsidy, One), TotalSubsidy)
	FinalLayer        = Ctx.Quo(new(decimal.Big), Ctx.Log(new(decimal.Big), new(decimal.Big).Sub(One, FinalIssuanceFrac)), NegLambda)
)

func getUnroundedAccumulatedSubsidy(layersAfterEffectiveGenesis uint32) *decimal.Big {
	// add one because layers are zero-indexed and we want > 0 issuance in the first effective genesis layer
	layerCount := new(decimal.Big).SetUint64(uint64(layersAfterEffectiveGenesis + 1))
	expInner := Ctx.Mul(new(decimal.Big), NegLambda, layerCount)
	expOuter := new(decimal.Big)
	Ctx.Exp(expOuter, expInner)
	one := decimal.New(1, 0)
	supplyMultiplier := Ctx.Sub(one, one, expOuter)
	return Ctx.Mul(new(decimal.Big), TotalSubsidy, supplyMultiplier)
}

// TotalAccumulatedSubsidyAtLayer returns the total accumulated block subsidy paid by the protocol as of the given
// layer, denominated in smidge.
func TotalAccumulatedSubsidyAtLayer(layersAfterEffectiveGenesis uint32) uint64 {
	unroundedSubsidy := getUnroundedAccumulatedSubsidy(layersAfterEffectiveGenesis)
	if ret, ok := unroundedSubsidy.Uint64(); !ok {
		log.Panicln("unable to convert subsidy to uint")
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
