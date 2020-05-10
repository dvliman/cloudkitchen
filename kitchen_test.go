package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKitchenPlaceOrderToShelfAccordingly(t *testing.T) {
	k := NewKitchen()

	before1 := len(k.Shelves[Hot].Orders)
	k.AcceptOrder(OrderReceived{Order: Order{ID: "1", Temp: "hot"}})
	assert.Len(t, k.Shelves[Hot].Orders, before1+1)

	before2 := len(k.Shelves[Cold].Orders)
	k.AcceptOrder(OrderReceived{Order: Order{ID: "1", Temp: "cold"}})
	assert.Len(t, k.Shelves[Cold].Orders, before2+1)

	before3 := len(k.Shelves[Frozen].Orders)
	k.AcceptOrder(OrderReceived{Order: Order{ID: "1", Temp: "frozen"}})
	assert.Len(t, k.Shelves[Frozen].Orders, before3+1)

	before4 := len(k.Shelves[Any].Orders)
	k.AcceptOrder(OrderReceived{Order: Order{ID: "1", Temp: "unknown"}})
	assert.Len(t, k.Shelves[Any].Orders, before4+1)
}

func TestKitchenPlaceOrderToOverflowIfFull(t *testing.T) {
	k := NewKitchen()
	for i := range make([]int, 10) {
		k.AcceptOrder(OrderReceived{Order: Order{ID: strconv.Itoa(i), Temp: "hot"}})
	}
	k.AcceptOrder(OrderReceived{Order: Order{ID: "11", Temp: "hot"}})
	assert.Len(t, k.Shelves[Hot].Orders, 10)
	assert.Len(t, k.Shelves[Any].Orders, 1)
}

func TestKitchenRandomlyDiscardFromOverflowIfNoMoreRoom(t *testing.T) {
	k := NewKitchen()
	for i := range make([]int, 10) {
		k.AcceptOrder(OrderReceived{Order: Order{ID: strconv.Itoa(i), Temp: "hot"}})
		k.AcceptOrder(OrderReceived{Order: Order{ID: strconv.Itoa(i), Temp: "cold"}})
		k.AcceptOrder(OrderReceived{Order: Order{ID: strconv.Itoa(i), Temp: "frozen"}})
		k.AcceptOrder(OrderReceived{Order: Order{ID: strconv.Itoa(i), Temp: "any"}})
	}
	for i := range make([]int, 5) { // overflow default capacity is 15
		k.AcceptOrder(OrderReceived{Order: Order{ID: strconv.Itoa(i), Temp: "any"}})
	}

	k.AcceptOrder(OrderReceived{Order: Order{ID: "last-one", Temp: "hot"}})
	assert.Len(t, k.Shelves[Hot].Orders, 10)
	assert.Len(t, k.Shelves[Cold].Orders, 10)
	assert.Len(t, k.Shelves[Frozen].Orders, 10)
	assert.Len(t, k.Shelves[Any].Orders, 15) // capped
}

func TestKitchenPickupOrderByID(t *testing.T) {
	k := NewKitchen()
	k.AcceptOrder(OrderReceived{Order: Order{ID: "1", Temp: "hot"}})

	before := len(k.Shelves[Hot].Orders)
	assert.True(t, k.PickupOrderByID("1"))
	assert.Len(t, k.Shelves[Hot].Orders, before-1)

	assert.False(t, k.PickupOrderByID("1")) // after removed, not found
}

func TestSelectUnknownTemperature(t *testing.T) {
	k := NewKitchen()
	shelf, err := k.selectShelfByTemperature("unknown")
	assert.Nil(t, shelf)
	assert.Error(t, err, ErrInvalidTemperatureLookup.Error())
}

func TestKitchenRemoveExpiredOrders(t *testing.T) {
	k := NewKitchen()
	k.AcceptOrder(OrderReceived{Order: Order{
		ID:        "1",
		Name:      "some-food-name",
		Temp:      "hot",
		ShelfLife: 0,
		DecayRate: 100,
	}})
	assert.Len(t, k.Shelves[Hot].Orders, 1)
	k.RemoveExpiredOrders()
	assert.Len(t, k.Shelves[Hot].Orders, 0)
}

func TestKitchenShelvesContent(t *testing.T) {
	k := NewKitchen()
	assert.NotEmpty(t, k.ShelvesContent())
}
