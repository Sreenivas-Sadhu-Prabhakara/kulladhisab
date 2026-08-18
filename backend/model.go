package backend

import (
	"encoding/json"
	"fmt"
)

// DayLog is one day's consumable spend against cups/plates served. The selling
// price is fixed by habit, so the true cost per cup is what moves.
type DayLog struct {
	Date        string  `json:"date"`
	Spend       float64 `json:"spend"`       // milk + tea + sugar + gas for the day
	CupsSold    float64 `json:"cupsSold"`    // approximate tally, by design
	PricePerCup float64 `json:"pricePerCup"` // your fixed selling price
}

// Validate reports whether the DayLog is well formed.
func (d DayLog) Validate() error {
	if d.Spend < 0 || d.PricePerCup < 0 {
		return fmt.Errorf("spend and price cannot be negative")
	}
	if d.CupsSold <= 0 {
		return fmt.Errorf("cups sold must be positive")
	}
	return nil
}

// CostPerCup is spend divided by cups served.
func (d DayLog) CostPerCup() float64 { return d.Spend / d.CupsSold }

// Summary aggregates the cost-per-cup trend.
type Summary struct {
	Days           int     `json:"days"`
	TotalSpend     float64 `json:"totalSpend"`
	TotalCups      float64 `json:"totalCups"`
	AvgCostPerCup  float64 `json:"avgCostPerCup"`
	AvgPricePerCup float64 `json:"avgPricePerCup"`
	AvgMargin      float64 `json:"avgMargin"`
}

// Summarize blends the logged days into an overall cost-per-cup and margin.
func Summarize(records []Record) Summary {
	var s Summary
	var priceSum float64
	for _, r := range records {
		var d DayLog
		if json.Unmarshal(r.Input, &d) != nil {
			continue
		}
		s.Days++
		s.TotalSpend += d.Spend
		s.TotalCups += d.CupsSold
		priceSum += d.PricePerCup
	}
	if s.TotalCups > 0 {
		s.AvgCostPerCup = s.TotalSpend / s.TotalCups
	}
	if s.Days > 0 {
		s.AvgPricePerCup = priceSum / float64(s.Days)
	}
	s.AvgMargin = s.AvgPricePerCup - s.AvgCostPerCup
	return s
}

// parseEntry decodes+validates a day log; headline is its cost per cup.
func parseEntry(raw []byte) (float64, string, error) {
	var d DayLog
	if err := json.Unmarshal(raw, &d); err != nil {
		return 0, "", fmt.Errorf("invalid json")
	}
	if err := d.Validate(); err != nil {
		return 0, "", err
	}
	return d.CostPerCup(), d.Date, nil
}
