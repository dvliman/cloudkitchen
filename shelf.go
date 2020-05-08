package main

import "math/rand"

type Shelf struct {
	Name                 string
	AllowableTemperature Temperature
	Capacity             int
	Orders               []OrderReceived
}

func (s *Shelf) IsFull() bool {
	return len(s.Orders) >= s.Capacity
}

func (s *Shelf) StoreOrder(order OrderReceived) {
	s.Orders = append(s.Orders, order)
}

func (s *Shelf) ThrowAwayExpiredOrder() {
	s.filter(func(od OrderReceived) bool {
		return s.computeShelfLife(od) > 0
	})
}

func (s *Shelf) DispenseFood(orderID string) {
	s.filter(func(od OrderReceived) bool {
		return od.Order.ID != orderID
	})
}

func (s *Shelf) RandomlyDiscardOneOrder() {
	s.DispenseFood(s.selectRandomOrderID())
}

func (s *Shelf) selectRandomOrderID() string {
	random := rand.Intn(len(s.Orders))
	return s.Orders[random].Order.ID
}

func (s *Shelf) filter(predicate func(OrderReceived) bool) {
	var xs []OrderReceived

	for _, x := range s.Orders {
		if predicate(x) {
			xs = append(xs, x)
		}
	}

	s.Orders = xs
}

func (s *Shelf) decayModifier() int {
	if s.AllowableTemperature == Any {
		return 2
	}

	return 1
}

func (s *Shelf) computeShelfLife(receivedOrder OrderReceived) float64 {
	order := receivedOrder.Order
	orderAge := receivedOrder.QueuedTime.Unix()

	return (float64(order.ShelfLife) - order.DecayRate) * float64(orderAge) * float64(s.decayModifier()) /
		float64(order.ShelfLife)
}
