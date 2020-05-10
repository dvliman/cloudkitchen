package main

import (
	"fmt"
	"log"
	"strings"
)

type Kitchen struct {
	Shelves map[Temperature]*Shelf
}

func NewKitchen() Kitchen {
	return Kitchen{
		Shelves: map[Temperature]*Shelf{
			Hot: {
				Name:                  "Hot Shelf",
				AllowableTemperatures: []Temperature{Hot},
				Capacity:              10,
			},
			Cold: {
				Name:                  "Cold Shelf",
				AllowableTemperatures: []Temperature{Cold},
				Capacity:              10,
			},
			Frozen: {
				Name:                  "Frozen Shelf",
				AllowableTemperatures: []Temperature{Frozen},
				Capacity:              10,
			},
			Any: { // not exactly map{temperature->shelf}
				Name:                  "Overflow Shelf",
				AllowableTemperatures: []Temperature{Hot, Cold, Frozen},
				Capacity:              15,
			},
		},
	}
}

func (k Kitchen) AcceptOrder(order OrderReceived) {
	firstOptionShelf, err := k.selectShelfByTemperature(order.Order.Temp)
	if err != nil {
		log.Printf("AcceptOrder: can not selectShelfByTemperature: "+
			"order.ID=%s order.Temperature=%s\n", order.Order.ID, order.Order.Temp)
		k.overflowShelf().PlaceOrder(order)
		return
	}

	if !firstOptionShelf.IsFull() {
		firstOptionShelf.PlaceOrder(order)
		return
	}

	if k.overflowShelf().IsFull() {
		hasRoom := k.canMoveOneOrderToAnotherShelf()
		if !hasRoom {
			toDiscard, _ := k.overflowShelf().GetRandomOrderIndex()
			k.overflowShelf().RemoveOrderAtIndex(toDiscard)
		}
	}

	k.overflowShelf().PlaceOrder(order)
}

func (k Kitchen) canMoveOneOrderToAnotherShelf() bool {
	for i, order := range k.overflowShelf().Orders {
		if targetShelf, err := k.selectShelfByTemperature(order.Order.Temp); err != nil && targetShelf != nil {
			if !targetShelf.IsFull() {
				targetShelf.PlaceOrder(order)
				k.overflowShelf().RemoveOrderAtIndex(i)
				return true
			}
		}
	}
	return false
}

func (k Kitchen) RemoveExpiredOrders() {
	for _, shelf := range k.Shelves {
		shelf.RemoveExpiredOrders()
	}
}

func (k Kitchen) PickupOrderByID(orderID string) bool {
	for _, shelf := range k.Shelves { // order might be placed in different shelf if overflow
		removed := shelf.RemoveOrderByID(orderID)
		if removed {
			return true
		}
	}

	return false
}

func (k Kitchen) selectShelfByTemperature(orderTemperature string) (*Shelf, error) {
	temperature, found := temperatureLookup[orderTemperature]
	if !found {
		return nil, ErrInvalidTemperatureLookup
	}

	return k.Shelves[temperature], nil
}

func (k Kitchen) overflowShelf() *Shelf {
	return k.Shelves[Any]
}

func (k Kitchen) ShelvesContent() string {
	var sb strings.Builder

	for _, shelf := range k.Shelves {
		_, _ = fmt.Fprintf(&sb, "\tShelf: name=%s capacity=%d allowableTemperatures=%s ordersCount=%d orders=%s\n",
			shelf.Name, shelf.Capacity, shelf.AllowableTemperatures, len(shelf.Orders), shelf.GetOrderIDs())
	}

	return sb.String()
}
