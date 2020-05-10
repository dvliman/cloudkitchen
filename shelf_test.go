package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewShelf(t *testing.T) {
	s1, err := NewShelf("shelf1", []Temperature{}, 0)
	assert.NoError(t, err)
	assert.True(t, s1.IsFull())

	s2, err := NewShelf("shelf2", []Temperature{}, 1)
	assert.NoError(t, err)
	assert.False(t, s2.IsFull())

	s3, err := NewShelf("shelf3", []Temperature{}, -1)
	assert.EqualError(t, err, ErrInvalidCapacity.Error())
	assert.Nil(t, s3)
}

func TestShelfPlaceOrder(t *testing.T) {
	s, err := NewShelf("shelf", []Temperature{}, 1)
	assert.NoError(t, err)

	before := len(s.Orders)
	s.PlaceOrder(OrderReceived{})
	after := len(s.Orders)

	assert.Equal(t, 1, after-before)
	assert.True(t, s.IsFull())
}

func TestShelfPlaceOrderWhenFull(t *testing.T) {
	s, err := NewShelf("shelf", []Temperature{}, 0)
	assert.NoError(t, err)

	before := len(s.Orders)
	s.PlaceOrder(OrderReceived{})
	after := len(s.Orders)

	assert.Equal(t, before, after)
}

func TestShelfDecayModifier(t *testing.T) {
	assert.Equal(t, 1, (&Shelf{AllowableTemperatures: []Temperature{Hot}}).decayModifier())
	assert.Equal(t, 1, (&Shelf{AllowableTemperatures: []Temperature{Cold}}).decayModifier())
	assert.Equal(t, 1, (&Shelf{AllowableTemperatures: []Temperature{Frozen}}).decayModifier())
	assert.Equal(t, 2, (&Shelf{AllowableTemperatures: []Temperature{}}).decayModifier())
}

func TestShelfComputeShellLife(t *testing.T) {
	tests := []struct {
		ShelfLife        int
		DecayRate        float64
		ShelfTemperature []Temperature
	}{
		{10, 1, []Temperature{Hot}},   // float-value: 8500259669165362p-73, float-to-int: 0
		{0, 10, []Temperature{Hot}},   // float-value: -Inf, float-to-int: -9223372036854775808
		{10, 0.5, []Temperature{Hot}}, // float-value: 0p-1074, float-to-int: 0
		{-1, 10, []Temperature{Hot}},  // float-value: 0p-1074, float-to-int: 0
		{1, 1, []Temperature{Hot}},    // float-value: 0p-1074, float-to-int: 0
		{10, 10, []Temperature{Hot}},  // float-value: 0p-1074, float-to-int: 0
	}

	for _, test := range tests {
		shelf := &Shelf{AllowableTemperatures: test.ShelfTemperature}

		order := OrderReceived{
			QueuedTime: time.Now().UTC(),
			Order: Order{
				ShelfLife: test.ShelfLife, // what unit?
				DecayRate: test.DecayRate, // what unit?
			},
		}

		result := shelf.computeShelfLife(order)
		PrintResult(result)
	}
}

func PrintResult(x int) {}

func TestShelfRemoveOrderAtIndex(t *testing.T) {
	s := &Shelf{Orders: []OrderReceived{
		{Order: Order{ID: "1"}},
		{Order: Order{ID: "2"}},
	}}
	assert.Len(t, s.Orders, 2)
	s.RemoveOrderAtIndex(0)
	assert.Len(t, s.Orders, 1)
	assert.Equal(t, s.Orders[0].Order.ID, "2")
}

func TestShelfRemoveOrderByID(t *testing.T) {
	s := &Shelf{Orders: []OrderReceived{
		{Order: Order{ID: "1"}},
	}}
	originalSize := len(s.Orders)

	// when: remove ID that doesn't exist
	// assert: nothing happened
	assert.False(t, s.RemoveOrderByID("2")) // remove that doesn't exist
	assert.Equal(t, len(s.Orders), originalSize)

	// when: remove existing ID
	// assert: removed
	assert.True(t, s.RemoveOrderByID("1"))
	assert.Equal(t, len(s.Orders), originalSize-1)
}

func TestShelfGetRandomOrderIndex(t *testing.T) {
	emptyShelf := &Shelf{}
	_, err := emptyShelf.GetRandomOrderIndex()
	assert.Error(t, err, ErrEmptyShelfOrders.Error())

	oneItem := &Shelf{Orders: []OrderReceived{{}}}
	oneIndex, err := oneItem.GetRandomOrderIndex()
	assert.NoError(t, err)

	// len([orders])=1, random=[0, 1]
	assert.GreaterOrEqual(t, oneIndex, 0)
	assert.LessOrEqual(t, oneIndex, 1)
}

func TestShelfGetOrderIDs(t *testing.T) {
	s, err := NewShelf("shelf", []Temperature{}, 1)
	assert.NoError(t, err)

	s.PlaceOrder(OrderReceived{Order: Order{ID: "1"}})
	assert.Contains(t, s.GetOrderIDs(), "1")
}
