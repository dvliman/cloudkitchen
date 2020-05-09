package main

import "time"

type Order struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Temp      string  `json:"temp"`
	ShelfLife int     `json:"shelfLife"`
	DecayRate float64 `json:"decayRate"`
}

type OrderReceived struct {
	Order      Order
	QueuedTime time.Time
	PickupTime *time.Timer
}
