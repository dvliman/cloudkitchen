package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Order struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Temp      string  `json:"temp"`
	ShelfLife int     `json:"shelfLife"`
	DecayRate float64 `json:"decayRate"`
}

type QueuedOrder struct {
	Order      Order
	QueuedTime time.Time
}

type Temperature string

const (
	Hot    Temperature = "hot"
	Cold   Temperature = "cold"
	Frozen Temperature = "frozen"
	Any    Temperature = "any"
)

type Shelf struct {
	Name                 string
	AllowableTemperature Temperature
	Capacity             int
	CookedOrders         []QueuedOrder
}

func (s *Shelf) IsFull() bool {
	return len(s.CookedOrders) >= s.Capacity
}

func (s *Shelf) StoreCookedOrder(order QueuedOrder) {
	s.CookedOrders = append(s.CookedOrders, order)
}

func (s *Shelf) ThrowAwayWastedFood() {
	var cookedOrdersToKeep []QueuedOrder

	for _, cookedOrder := range s.CookedOrders {
		if computeShelfLife(s, cookedOrder) > 0 {
			cookedOrdersToKeep = append(cookedOrdersToKeep, cookedOrder)
		}
	}

	s.CookedOrders = cookedOrdersToKeep
}

func (s *Shelf) RandomlyDiscardOneFood() {
	var cookedOrdersToKeep []QueuedOrder
	random := rand.Intn(len(s.CookedOrders))

	for i, cookedOrder := range s.CookedOrders {
		if i != random {
			cookedOrdersToKeep = append(cookedOrdersToKeep, cookedOrder)
		}
	}

	s.CookedOrders = cookedOrdersToKeep
}

func shelfDecayModifier(s *Shelf) int {
	if s.AllowableTemperature == Any {
		return 2
	}

	return 1
}

func computeShelfLife(s *Shelf, queuedOrder QueuedOrder) float64 {
	order := queuedOrder.Order
	orderAge := queuedOrder.QueuedTime.Unix()

	return (float64(order.ShelfLife) - order.DecayRate) * float64(orderAge) * float64(shelfDecayModifier(s)) /
		float64(order.ShelfLife)
}

func main() {
	kitchen := NewKitchen()
	kitchen.PeriodicallyRemoveExpiredOrder()

	orders := readOrdersFromFile()

	rateLimiter := time.Tick(rateLimit)
	queuedOrders := make(chan QueuedOrder, len(orders))

	acceptOrders(rateLimiter, queuedOrders, orders) // TODO: make rate limiter configurable
	shelfCookedOrder(kitchen, queuedOrders)
}

func readOrdersFromFile() []Order {
	jsonFile, err := os.Open("./orders.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	var orders []Order
	if err := json.Unmarshal(bytes, &orders); err != nil {
		panic(err)
	}

	return orders
}

const rateLimit = time.Second / 2 // 2 orders per seconds

func acceptOrders(rateLimiter <-chan time.Time, queuedOrders chan<- QueuedOrder, orders []Order) {
	fmt.Printf("accepting orders: %d\n", len(orders))

	for i, order := range orders {
		<-rateLimiter
		fmt.Printf("count: %d, time: %s\n", i, time.Now().UTC().String())
		queuedOrders <- QueuedOrder{Order: order, QueuedTime: time.Now().UTC()}
	}

	fmt.Printf("done accepting orders")
}

type Kitchen struct { // store the state of shelves with CookedOrderCount
	Shelves map[string]*Shelf
}

func NewKitchen() Kitchen {
	return Kitchen{Shelves: map[string]*Shelf{
		"hot":      {Name: "Hot Shelf", AllowableTemperature: Hot, Capacity: 10},
		"cold":     {Name: "Cold Shelf", AllowableTemperature: Cold, Capacity: 10},
		"frozen":   {Name: "Frozen Shelf", AllowableTemperature: Frozen, Capacity: 10},
		"overflow": {Name: "Overflow Shelf", AllowableTemperature: Any, Capacity: 15},
	}}
}

func (k Kitchen) ShelfOrder(queuedOrder QueuedOrder) {
	primaryShelf := k.Shelves[queuedOrder.Order.Temp]

	if !primaryShelf.IsFull() {
		primaryShelf.StoreCookedOrder(queuedOrder)

	} else {
		secondaryShelf := k.Shelves["overflow"]
		if secondaryShelf.IsFull() {
			secondaryShelf.RandomlyDiscardOneFood()
		}
		secondaryShelf.StoreCookedOrder(queuedOrder)
	}
}

func (k Kitchen) PeriodicallyRemoveExpiredOrder() {
	go func() {
		for _, shelf := range k.Shelves {
			shelf.ThrowAwayWastedFood()
		}
	}()
}

func shelfCookedOrder(kitchen Kitchen, queuedOrders <-chan QueuedOrder) {
	for {
		select {
		case queuedOrder := <-queuedOrders:
			kitchen.ShelfOrder(queuedOrder)
		}
	}
}
