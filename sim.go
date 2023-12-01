package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spacemeshos/economics/constants"
	"github.com/spacemeshos/economics/rewards"
	"github.com/spacemeshos/economics/vesting"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tcnksm/go-input"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var qFlag = flag.Bool("q", false, "quiet mode (noninteractive)")

func main() {
	// parse flags
	flag.Parse()

	oneLayer := 5 * time.Minute

	currentDate, tickInterval, endLayer := getParams()
	log.Printf("genesis is %s\n", currentDate)
	log.Printf("effective genesis is/issuance begins %s\n", currentDate.Add(effectiveGenesis*oneLayer))
	log.Printf("tick interval is %d layers\n", tickInterval)
	log.Printf("last layer is %d\n", endLayer)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"layer",
		"epoch",
		"date",
		"vaultNewVest",
		"vaultTotalVest",
		"vaultPctVest",
		"vaultTotal",
		"subsidyPerLayer",
		"subsidyNew",
		"subsidyTotal",
		"circulatingTotal",
		"issuanceTotal",
		"pctVault",
		"pctCirculating",
		"pctFinalIssuance",
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
		{Number: 6, Align: text.AlignRight},
		{Number: 8, Align: text.AlignRight},
		{Number: 9, Align: text.AlignRight},
		{Number: 10, Align: text.AlignRight},
		{Number: 11, Align: text.AlignRight},
		{Number: 12, Align: text.AlignRight},
		{Number: 14, Align: text.AlignRight},
		{Number: 15, Align: text.AlignRight},
	})
	t.SetCaption("Please note:\n" +
		"- All figures in SMESH (rounded down)\n" +
		"- No coins are issued in the first two epochs\n" +
		"- Figures represent maximum issuance (and do not account for empty layers)\n")

	p := message.NewPrinter(language.English)

	pw := progress.NewWriter()
	pw.SetUpdateFrequency(time.Millisecond * 100)

	// don't render progress bar in quiet mode
	if !*qFlag {
		go pw.Render()
		defer pw.Stop()
	}
	tracker := progress.Tracker{Total: int64(endLayer), Units: progress.Units{
		Formatter:        progress.FormatNumber,
		Notation:         " layers",
		NotationPosition: progress.UnitsNotationPositionAfter,
	}}
	pw.AppendTracker(&tracker)
	trackerTickInterval := 1000

	vaultTotal := uint64(constants.TotalVaulted)
	issuanceTotal := vaultTotal // vaulted amount is issued but not circulating yet

	var vaultVested, subsidyTotal, circulatingTotal, vaultNewVest, subsidyNew uint64

	// note: we could optimize this and just step by tick interval, but we do the simplest possible thing here and get
	// as close as possible to reality by stepping through every single layer
	for layerID := uint32(0); layerID <= endLayer; layerID++ {
		// update vault
		// vault vesting is calculated on the basis of layers post-genesis
		vaultVested = vesting.AccumulatedVestAtLayer(layerID)
		vestThisLayer := vesting.VestAtLayer(layerID)
		vaultNewVest += vestThisLayer
		circulatingTotal += vestThisLayer

		// add new issuance
		// issuance is calculated on the basis of layers post-effective genesis
		// and no issuance occurs before effective genesis
		var subsidyTotalNew, subsidyThisLayer uint64
		if layerID >= effectiveGenesis {
			// calculate effective layer, i.e., layers post-effective-genesis
			effectiveLayer := layerID - effectiveGenesis
			subsidyTotalNew = rewards.TotalAccumulatedSubsidyAtLayer(effectiveLayer)
			subsidyThisLayer = subsidyTotalNew - subsidyTotal
		}

		circulatingTotal += subsidyThisLayer
		issuanceTotal += subsidyThisLayer
		subsidyNew += subsidyThisLayer
		subsidyTotal = subsidyTotalNew

		// increment here in case tick interval is really big
		if layerID > 0 && layerID%uint32(trackerTickInterval) == 0 {
			tracker.Increment(int64(trackerTickInterval))
		}

		if layerID%tickInterval == 0 || layerID == endLayer {
			t.AppendRow(table.Row{
				layerID,
				layerID / constants.OneEpoch,
				currentDate.Format("2006-01-02"),
				p.Sprintf("%7d", vaultNewVest/constants.OneSmesh),
				p.Sprintf("%11d", vaultVested/constants.OneSmesh),
				p.Sprintf("%7.2f%%", 100*float64(vaultVested)/float64(vaultTotal)),
				p.Sprintf("%d", vaultTotal/constants.OneSmesh),
				p.Sprintf("%7d", subsidyThisLayer/constants.OneSmesh),
				p.Sprintf("%7d", subsidyNew/constants.OneSmesh),
				p.Sprintf("%11d", subsidyTotal/constants.OneSmesh),
				p.Sprintf("%11d", circulatingTotal/constants.OneSmesh),
				p.Sprintf("%11d", issuanceTotal/constants.OneSmesh),
				p.Sprintf("%7.2f%%", 100*float64(vaultTotal)/float64(issuanceTotal)),
				p.Sprintf("%7.2f%%", 100*float64(circulatingTotal)/float64(issuanceTotal)),
				p.Sprintf("%7.2f%%", 100*float64(issuanceTotal)/float64(constants.TotalIssuance)),
			})

			// reset these
			vaultNewVest = 0
			subsidyNew = 0
		}
		currentDate = currentDate.Add(oneLayer)
	}
	tracker.MarkAsDone()
	t.Render()
}

const (
	// Actual mainnet genesis
	defaultGenesisDateStr = "20230714"

	// One mainnet epoch
	defaultTickInterval = constants.OneEpoch

	// Effective genesis occurs two epochs post-genesis
	effectiveGenesis = 2 * constants.OneEpoch

	// Issuance begins at effective genesis; we reach the ten year target ten years post-effective genesis
	defaultEndLayer = 10*constants.OneYear + effectiveGenesis
)

var defaultGenesisDate, _ = time.Parse("20060102", defaultGenesisDateStr)

func getParams() (time.Time, uint32, uint32) {
	// short-circuit UI in quiet mode
	if *qFlag {
		return defaultGenesisDate, defaultTickInterval, defaultEndLayer
	}

	ui := &input.UI{}
	var genesisDate time.Time
	if genesisDateStr, err := ui.Ask("genesis date (YYYYMMDD)", &input.Options{
		Default:   defaultGenesisDateStr,
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

	defaultTickIntervalStr := fmt.Sprintf("%d (one epoch/two weeks)", defaultTickInterval)
	var tickInterval int
	if tickIntervalStr, err := ui.Ask("layer tick interval", &input.Options{
		Default:   defaultTickIntervalStr,
		HideOrder: true,
		Required:  true,
		Loop:      true,
		ValidateFunc: func(s string) (err error) {
			_, err = strconv.Atoi(s)
			return
		},
	}); err != nil {
		log.Fatal(err)
	} else if tickIntervalStr == defaultTickIntervalStr {
		tickInterval = defaultTickInterval
	} else {
		tickInterval, _ = strconv.Atoi(tickIntervalStr)
	}

	defaultEndLayerStr := fmt.Sprintf("%d (ten years)", defaultEndLayer)
	var endLayer int
	if endLayerStr, err := ui.Ask("end layer", &input.Options{
		Default:   defaultEndLayerStr,
		HideOrder: true,
		Required:  true,
		Loop:      true,
		ValidateFunc: func(s string) (err error) {
			_, err = strconv.Atoi(s)
			return
		},
	}); err != nil {
		log.Fatal(err)
	} else if endLayerStr == defaultEndLayerStr {
		endLayer = defaultEndLayer
	} else {
		endLayer, _ = strconv.Atoi(endLayerStr)
	}

	return genesisDate, uint32(tickInterval), uint32(endLayer)
}
