package backend

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type memStore struct{ items []Record }

func (m *memStore) Save(r Record) (Record, error) {
	r.ID = int64(len(m.items) + 1)
	m.items = append([]Record{r}, m.items...)
	return r, nil
}
func (m *memStore) List(limit int) ([]Record, error) { return m.items, nil }

func mk(spend, cups, price float64) Record {
	d := DayLog{Spend: spend, CupsSold: cups, PricePerCup: price}
	in, _ := json.Marshal(d)
	return Record{Input: in, Headline: d.CostPerCup(), Label: ""}
}

func TestSummarize_CostPerCupAndMargin(t *testing.T) {
	// Day1: 600 spend, 100 cups (₹6/cup); Day2: 800 spend, 100 cups (₹8/cup); price ₹10.
	s := Summarize([]Record{mk(600, 100, 10), mk(800, 100, 10)})
	if math.Abs(s.AvgCostPerCup-7) > 1e-9 { // 1400/200
		t.Fatalf("avgCostPerCup=%v want 7", s.AvgCostPerCup)
	}
	if math.Abs(s.AvgMargin-3) > 1e-9 {
		t.Fatalf("avgMargin=%v want 3", s.AvgMargin)
	}
}

func TestValidate(t *testing.T) {
	if err := (DayLog{Spend: 600, CupsSold: 100}).Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []DayLog{{CupsSold: 0}, {CupsSold: 10, Spend: -1}} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad %d accepted", i)
		}
	}
}

func TestLogAndSummary(t *testing.T) {
	srv := NewServer(&memStore{})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/log",
		strings.NewReader(`{"date":"2026-08-18","spend":600,"cupsSold":100,"pricePerCup":10}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("log %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/summary", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"avgCostPerCup":6`) {
		t.Fatalf("summary=%s", rec.Body.String())
	}
}
