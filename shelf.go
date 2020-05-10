package main

import (
	"math/rand"
	"time"
)

type Shelf struct {
	Name                  string
	AllowableTemperatures []Temperature
	Capacity              int
	Orders                []OrderReceived
}

func NewShelf(name string, allowableTemperature []Temperature, capacity int) (*Shelf, error) {
	if capacity < 0 {
		return nil, ErrInvalidCapacity
	}

	return &Shelf{
		Name:                  name,
		AllowableTemperatures: allowableTemperature,
		Capacity:              capacity,
		Orders:                []OrderReceived{},
	}, nil
}

func (s *Shelf) IsFull() bool {
	return len(s.Orders) >= s.Capacity
}

func (s *Shelf) PlaceOrder(order OrderReceived) {
	if !s.IsFull() {
		s.Orders = append(s.Orders, order)
	}
}

func (s *Shelf) RemoveExpiredOrders() {
	s.filter(func(_ int, od OrderReceived) bool {
		return s.computeShelfLife(od) > 0
	})
}

func (s *Shelf) RemoveOrderByID(orderID string) bool {
	return s.filter(func(_ int, od OrderReceived) bool {
		return od.Order.ID != orderID
	})
}

func (s *Shelf) RemoveOrderAtIndex(index int) bool {
	return s.filter(func(i int, od OrderReceived) bool {
		return i != index
	})
}

func (s *Shelf) GetRandomOrderIndex() (int, error) {
	if len(s.Orders) == 0 {
		return 0, ErrEmptyShelfOrders
	}
	return rand.Intn(len(s.Orders)), nil
}

func (s *Shelf) filter(predicate func(int, OrderReceived) bool) bool {
	var xs []OrderReceived

	for i, x := range s.Orders {
		if predicate(i, x) {
			xs = append(xs, x)
		}
	}

	before := len(s.Orders)
	s.Orders = xs
	after := len(s.Orders)

	return (before - after) > 0 // any filtered?
}

func (s *Shelf) decayModifier() int {
	if len(s.AllowableTemperatures) == 1 {
		return 1
	}

	return 2
}

func (s *Shelf) computeShelfLife(receivedOrder OrderReceived) int {
	order := receivedOrder.Order
	orderAge := time.Since(receivedOrder.QueuedTime).Seconds() // nanoseconds

	lifeWithDecay := float64(order.ShelfLife) - order.DecayRate
	lifeDecayedInTheShelf := lifeWithDecay * orderAge * float64(s.decayModifier())
	actualShellLife := lifeDecayedInTheShelf / float64(order.ShelfLife)
	// fmt.Printf("float-value: %b, float-to-int: %d\n", actualShellLife, int(actualShellLife))
	return int(actualShellLife)
}

func (s *Shelf) GetOrderIDs() []string {
	var xs []string
	for _, order := range s.Orders {
		xs = append(xs, order.Order.ID)
	}
	return xs
}
