package main

import (
	"fmt"

	"github.com/spacemeshos/economics/rewards"
)

func main() {
	for i := uint32(0); i < 100; i++ {
		layerReward := rewards.SubsidyByLayer(i)
		fmt.Println("Layer: ", i, "; Reward: ", layerReward)
	}
}
