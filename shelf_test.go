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
	d1, _ := time.ParseDuration("-50s")
	d2, _ := time.ParseDuration("-800s")

	now := time.Now().UTC()
	t1 := now.Add(d1)
	t2 := now.Add(d2)

	tests := []struct {
		ShelfLife        int     // in seconds
		DecayRate        float64 // deterioration modifier
		ShelfTemperature []Temperature
		QueuedTime       time.Time // when the order was received
		CurrentTime      time.Time // time when kitchen is expiring the shelves
		Expired          bool
	}{
		// from the example: shelfLife=300, decayRate=0.45
		{300, 0.45, []Temperature{Hot}, now, now, false}, // as it happens, expired = false
		{300, 1, []Temperature{Hot}, now, now, false},    // as it happens but high deterioration rate, expired = false
		{300, 0.45, []Temperature{Hot}, t1, now, false},  // within the shelflife * modifier, expired = false
		{300, 0.45, []Temperature{Hot}, t2, now, true},   // queued much longer than shelfLife, expired = true
		{300, 0, []Temperature{Hot}, t2, now, false},     // queued much longer than shelfLife but never deteriorate, expired = false
		{0, 1, []Temperature{Hot}, now, now, false},      // no shelfLife at all, expired = true (immediately)
	}

	for _, test := range tests {
		shelf := &Shelf{AllowableTemperatures: test.ShelfTemperature}

		order := OrderReceived{
			QueuedTime: test.QueuedTime,
			Order: Order{
				ShelfLife: test.ShelfLife,
				DecayRate: test.DecayRate,
			},
		}

		expired := shelf.computeShelfLife(order, test.CurrentTime) < 0.0
		assert.Equal(t, expired, test.Expired)
	}
}

func TestRemoveExpiredOrders(t *testing.T) {
	s, _ := NewShelf("some-name", []Temperature{Hot}, 1)
	s.PlaceOrder(OrderReceived{Order: Order{ShelfLife: 0, DecayRate: 0.45, Temp: "hot"}})

	before := len(s.Orders)
	s.RemoveExpiredOrders()
	after := len(s.Orders)

	assert.True(t, before-after == 1)
}

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
