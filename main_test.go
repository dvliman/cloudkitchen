package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeFirstTwo(t *testing.T) {
	head1, tail1 := takeFirstTwo([]Order{})
	assert.EqualValues(t, emptyOrder(), head1)
	assert.EqualValues(t, emptyOrder(), tail1)

	head2, tail2 := takeFirstTwo([]Order{{ID: "1"}})
	assert.EqualValues(t, []Order{{ID: "1"}}, head2)
	assert.EqualValues(t, emptyOrder(), tail2)

	head3, tail3 := takeFirstTwo([]Order{{ID: "1"}, {ID: "2"}})
	assert.EqualValues(t, []Order{{ID: "1"}, {ID: "2"}}, head3)
	assert.EqualValues(t, emptyOrder(), tail3)

	head4, tail4 := takeFirstTwo([]Order{{ID: "1"}, {ID: "2"}, {ID: "3"}})
	assert.EqualValues(t, []Order{{ID: "1"}, {ID: "2"}}, head4)
	assert.EqualValues(t, []Order{{ID: "3"}}, tail4)

	head5, tail5 := takeFirstTwo([]Order{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}})
	assert.EqualValues(t, []Order{{ID: "1"}, {ID: "2"}}, head5)
	assert.EqualValues(t, []Order{{ID: "3"}, {ID: "4"}}, tail5)

	head6, tail6 := takeFirstTwo([]Order{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}})
	assert.EqualValues(t, []Order{{ID: "1"}, {ID: "2"}}, head6)
	assert.EqualValues(t, []Order{{ID: "3"}, {ID: "4"}, {ID: "5"}}, tail6)
}

func emptyOrder() []Order {
	return []Order{}
}
