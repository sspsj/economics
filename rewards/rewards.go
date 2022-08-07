package rewards

import (
	"github.com/ericlagergren/decimal"
	"github.com/spacemeshos/economics/constants"
	"log"
)

var (
	Ctx          = decimal.Context128
	Two          = decimal.New(2, 0)
	HalfLife     = decimal.New(constants.HalfLife, 0)
	M            = new(decimal.Big)
	LogTwo       = Ctx.Log(M, Two)
	Lambda       = LogTwo.Quo(LogTwo, HalfLife)
	NegLambda    = new(decimal.Big).Copy(Lambda).Neg(Lambda)
	TotalSubsidy = new(decimal.Big).SetUint64(constants.TotalSubsidy)
)

func getUnroundedAccumulatedSubsidy(layersAfterEffectiveGenesis uint32) *decimal.Big {
	// add one because layers are zero-indexed and we want > 0 issuance in the first effective genesis layer
	layerCount := new(decimal.Big).SetUint64(uint64(layersAfterEffectiveGenesis + 1))
	expInner := new(decimal.Big).Mul(NegLambda, layerCount)
	expOuter := new(decimal.Big)
	Ctx.Exp(expOuter, expInner)
	one := decimal.New(1, 0)
	supplyMultiplier := one.Sub(one, expOuter)
	return new(decimal.Big).Mul(TotalSubsidy, supplyMultiplier)
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
