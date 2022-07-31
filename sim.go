package main

import (
	"github.com/tcnksm/go-input"
	"log"
	"strconv"
	"time"
)

func main() {
	genesisDate, tickInterval, endLayer := getParams()
	log.Printf("genesis date is %s\n", genesisDate)
	log.Printf("tick interval is %d\n", tickInterval)
	log.Printf("last layer is %d\n", endLayer)
}

func getParams() (time.Time, uint32, uint32) {
	// get params
	ui := &input.UI{}
	var genesisDate time.Time
	if genesisDateStr, err := ui.Ask("genesis date (YYYYMMDD)", &input.Options{
		Default:   "20230101",
		HideOrder: true,
		Required:  true,
		Loop:      true,
		ValidateFunc: func(s string) (err error) {
			_, err = time.Parse("20060102", s)
			return
		},
	}); err != nil {
		log.Fatal(err)
	} else {
		genesisDate, _ = time.Parse("20060102", genesisDateStr)
	}

	defaultTickVal := "2016 (one week)"
	var tickInterval int
	if tickIntervalStr, err := ui.Ask("layer tick interval", &input.Options{
		Default:   defaultTickVal,
		HideOrder: true,
		Required:  true,
		Loop:      true,
		ValidateFunc: func(s string) (err error) {
			_, err = strconv.Atoi(s)
			return
		},
	}); err != nil {
		log.Fatal(err)
	} else if tickIntervalStr == defaultTickVal {
		tickInterval = 2016
	} else {
		tickInterval, _ = strconv.Atoi(tickIntervalStr)
	}

	defaultEndLayer := "1051920 (ten years)"
	var endLayer int
	if endLayerStr, err := ui.Ask("end layer", &input.Options{
		Default:   defaultEndLayer,
		HideOrder: true,
		Required:  true,
		Loop:      true,
		ValidateFunc: func(s string) (err error) {
			_, err = strconv.Atoi(s)
			return
		},
	}); err != nil {
		log.Fatal(err)
	} else if endLayerStr == defaultEndLayer {
		endLayer = 1051920
	} else {
		endLayer, _ = strconv.Atoi(endLayerStr)
	}

	return genesisDate, uint32(tickInterval), uint32(endLayer)
	//for i := uint32(0); i < 100; i++ {
	//	layerReward := rewards.SubsidyByLayer(i)
	//	fmt.Println("Layer: ", i, "; Reward: ", layerReward)
	//}
}
