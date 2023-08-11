package internal

import "time"

type StopLossTrigger struct {
	Type       StopLossTriggerType
	Value      float64
	ValidUntil time.Time
}

type StopLossOrderEvent struct {
	Type                OrderType
	Price               float64
	Volume              float64
	ValidDays           int
	PriceType           StopLossPriceType
	ShortSellingAllowed bool
}
