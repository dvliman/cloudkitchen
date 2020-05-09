package main

import (
	"errors"
	"math/rand"
	"time"
)

type Kitchen struct {
	RateLimiter    <-chan time.Time
	OrdersReceived chan OrderReceived
	Shelves        map[Temperature]*Shelf
}

func NewKitchen(ordersPerSecond int, removeExpiredOrder bool) Kitchen {
	k := Kitchen{
		RateLimiter: time.Tick(time.Second / time.Duration(ordersPerSecond)),

		Shelves: map[Temperature]*Shelf{
			Hot:    {Name: "Hot Shelf", AllowableTemperature: Hot, Capacity: 10},
			Cold:   {Name: "Cold Shelf", AllowableTemperature: Cold, Capacity: 10},
			Frozen: {Name: "Frozen Shelf", AllowableTemperature: Frozen, Capacity: 10},
			Any:    {Name: "Overflow Shelf", AllowableTemperature: Any, Capacity: 15},
		},
	}

	if removeExpiredOrder {
		go func() {
			k.periodicallyRemoveExpiredOrder()
		}()
	}
	return k
}

func (k Kitchen) AcceptOrders(orders []Order) {
	for _, order := range orders {
		<-k.RateLimiter

		k.OrdersReceived <- OrderReceived{
			Order:      order,
			QueuedTime: time.Now().UTC(),
			PickupTime: time.AfterFunc(randomPickupTime(), func() {
				k.DispenseOrder(order.ID, order.Temp)
			}),
		}
	}

	k.storeToShelves()
}

func (k Kitchen) periodicallyRemoveExpiredOrder() {
	ticker := time.NewTicker(time.Minute) // TODO: make this configurable
	for {
		select {
		case <-ticker.C:
			for _, shelf := range k.Shelves {
				shelf.ThrowAwayExpiredOrder()
			}
		}
	}
}

func (k Kitchen) storeToShelves() {
	for {
		select {
		case od := <-k.OrdersReceived:
			primaryShelf, err := k.selectShelfByTemperature(od.Order.Temp)
			if err != nil {
				return
			}

			if !primaryShelf.IsFull() {
				primaryShelf.StoreOrder(od)
				return
			}

			overflowShelf := k.Shelves[Any]
			if overflowShelf.IsFull() {
				overflowShelf.RandomlyDiscardOneOrder()
			}

			overflowShelf.StoreOrder(od)
		}
	}
}

func (k Kitchen) DispenseOrder(orderID string, orderTemperature string) {
	shelf, err := k.selectShelfByTemperature(orderTemperature)
	if err != nil {
		return
	}
	shelf.DispenseFood(orderID)
}

func (k Kitchen) selectShelfByTemperature(orderTemperature string) (*Shelf, error) {
	temperature, found := temperatureLookup[orderTemperature]
	if !found {
		return nil, errors.New("selectShelfByTemperature/InvalidTemperatureLookup")
	}

	return k.Shelves[temperature], nil
}

func randomPickupTime() time.Duration {
	min := 2
	max := 3
	random := rand.Intn(max-min) + min
	return time.Second * time.Duration(random)
}
