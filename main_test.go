package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeFirst(t *testing.T) {
	head1, tail1 := first([]Order{})
	assert.Nil(t, head1)
	assert.EqualValues(t, emptyOrder(), tail1)

	head2, tail2 := first([]Order{{ID: "1"}})
	assert.EqualValues(t, &Order{ID: "1"}, head2)
	assert.EqualValues(t, emptyOrder(), tail2)

	head3, tail3 := first([]Order{{ID: "1"}, {ID: "2"}})
	assert.EqualValues(t, &Order{ID: "1"}, head3)
	assert.EqualValues(t, []Order{{ID: "2"}}, tail3)
}

func emptyOrder() []Order {
	return []Order{}
}
