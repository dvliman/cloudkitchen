package main

import "errors"

type Temperature string

const (
	Hot    Temperature = "hot"
	Cold   Temperature = "cold"
	Frozen Temperature = "frozen"
	Any    Temperature = "any"
)

var temperatureLookup = map[string]Temperature{
	"hot":    Hot,
	"cold":   Cold,
	"frozen": Frozen,
	"any":    Any,
}

var (
	ErrInvalidTemperatureLookup = errors.New("InvalidTemperatureLookup")
	ErrInvalidCapacity          = errors.New("InvalidCapacity")
	ErrEmptyShelfOrders         = errors.New("EmptyShelf")
)
