package main

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spacemeshos/economics/constants"
	"github.com/spacemeshos/economics/rewards"
	"github.com/spacemeshos/economics/vesting"
	"github.com/tcnksm/go-input"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	currentDate, tickInterval, endLayer := getParams()
	log.Printf("genesis date is %s\n", currentDate)
	log.Printf("tick interval is %d\n", tickInterval)
	log.Printf("last layer is %d\n", endLayer)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"layer",
		"date",
		"vaultNewVest",
		"vaultTotalVest",
		"vaultPctVested",
		"vaultTotal",
		"subsidyNew",
		"subsidyTotal",
		"circulatingTotal",
		"issuanceTotal",
	})

	p := message.NewPrinter(language.English)

	vaultTotal := uint64(constants.TotalVaulted)
	vaultVested := uint64(0)
	//subsidyNew := uint64(0)
	subsidyTotal := uint64(0)
	circulatingTotal := uint64(0)
	issuanceTotal := uint64(0)
	effectiveGenesis := uint32(0)

	oneLayer, _ := time.ParseDuration("5m")

	var vaultNewVest, subsidyNew uint64

	// note: we could optimize this and just step by tick interval, but we do the simplest possible thing here and get
	// as close as possible to reality by stepping through every single layer
	for layerID := uint32(0); layerID <= endLayer; layerID++ {
		// update vault
		vaultVested = vesting.AccumulatedVestAtLayer(effectiveGenesis, layerID)
		vaultNewVest += vesting.VestAtLayer(effectiveGenesis, layerID)
		circulatingTotal += vaultNewVest

		// add new issuance
		subsidyTotalNew := rewards.TotalAccumulatedSubsidyAtLayer(effectiveGenesis, layerID)
		subsidyThisLayer := subsidyTotalNew - subsidyTotal
		circulatingTotal += subsidyThisLayer
		issuanceTotal += subsidyThisLayer
		subsidyNew += subsidyThisLayer
		subsidyTotal = subsidyTotalNew

		if layerID%tickInterval == 0 {
			t.AppendRow([]interface{}{
				layerID,
				currentDate.Format("2006-01-02"),
				p.Sprintf("%d", vaultNewVest/constants.OneSmesh),
				p.Sprintf("%d", vaultVested/constants.OneSmesh),
				p.Sprintf("%0.2f", float64(vaultVested)/float64(vaultTotal)),
				p.Sprintf("%d", vaultTotal/constants.OneSmesh),
				p.Sprintf("%d", subsidyNew/constants.OneSmesh),
				p.Sprintf("%d", subsidyTotal/constants.OneSmesh),
				p.Sprintf("%d", circulatingTotal/constants.OneSmesh),
				p.Sprintf("%d", issuanceTotal/constants.OneSmesh),
			})

			// reset these
			vaultNewVest = 0
			subsidyNew = 0
		}
		currentDate = currentDate.Add(oneLayer)
		//fmt.Printf(".")
	}
	t.Render()
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
