package rewards

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IssuanceBeforeGenesis(t *testing.T) {
	genesisSubsidy := TotalAccumulatedSubsidyAtLayer(1, 0)
	genesisTotalSubsidy := TotalSubsidyAtLayer(1, 0)
	assert.Zero(t, genesisSubsidy)
	assert.Zero(t, genesisTotalSubsidy)
}
